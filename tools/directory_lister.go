package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DirectoryListerInput struct {
	DirectoryPath string `json:"directory_path,omitempty" jsonschema_description:"Path to the directory to list contents from."`
}

var DirectoryListerInputFormat = GenerateSchema[DirectoryListerInput]()

var DirectoryListerTool = ToolDefinition{
	ToolName:        "directory_lister",
	ToolDescription: "Lists files and directories at a specified path.",
	InputFormat:     DirectoryListerInputFormat,
	Handler:         ListDirectoryContents,
}

func ListDirectoryContents(input json.RawMessage) (string, error) {
	var dirInput DirectoryListerInput
	if err := json.Unmarshal(input, &dirInput); err != nil {
		return "", err
	}

	dir := "."
	if dirInput.DirectoryPath != "" {
		dir = dirInput.DirectoryPath
	}

	var items []string
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(dir, filePath)
		if err != nil {
			return err
		}

		if relativePath != "." {
			if info.IsDir() {
				items = append(items, relativePath+"/")
			} else {
				items = append(items, relativePath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(items)
	if err != nil {
		return "", err
	}

	return string(result), nil
}