package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type TextSearchInput struct {
	DirectoryPath  string   `json:"directory_path" jsonschema_description:"Path to the directory to search in."`
	Pattern        string   `json:"pattern" jsonschema_description:"String or regex pattern to search for."`
	FileExtensions []string `json:"file_extensions,omitempty" jsonschema_description:"Optional list of file extensions to filter by."`
}

var TextSearchInputFormat = GenerateSchema[TextSearchInput]()

var TextSearchTool = ToolDefinition{
	ToolName:        "text_search",
	ToolDescription: "Searches for a string or pattern in files under a given directory.",
	InputFormat:     TextSearchInputFormat,
	Handler:         TextSearchHandler,
}

func TextSearchHandler(input json.RawMessage) (string, error) {
	var searchInput TextSearchInput
	if err := json.Unmarshal(input, &searchInput); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if searchInput.DirectoryPath == "" || searchInput.Pattern == "" {
		return "", fmt.Errorf("directory_path and pattern are required")
	}

	var matchedFiles []string
	err := filepath.Walk(searchInput.DirectoryPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if len(searchInput.FileExtensions) > 0 {
			ext := filepath.Ext(filePath)
			if !contains(searchInput.FileExtensions, ext) {
				return nil
			}
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		matched, err := regexp.Match(searchInput.Pattern, content)
		if err != nil {
			return err
		}

		if matched {
			matchedFiles = append(matchedFiles, filePath)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(matchedFiles)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}