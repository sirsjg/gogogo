package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type ToolDefinition struct {
	ToolName        string                         `json:"tool_name"`
	ToolDescription string                         `json:"tool_description"`
	InputFormat     anthropic.ToolInputSchemaParam `json:"input_format"`
	Handler         func(input json.RawMessage) (string, error)
}

var FileReaderTool = ToolDefinition{
	ToolName:        "file_reader",
	ToolDescription: "Reads the content of a specified file.",
	InputFormat:     FileReaderInputFormat,
	Handler:         ReadFileContent,
}

type FileReaderInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path to the file to be read."`
}

var FileReaderInputFormat = GenerateSchema[FileReaderInput]()

func ReadFileContent(input json.RawMessage) (string, error) {
	var fileInput FileReaderInput
	if err := json.Unmarshal(input, &fileInput); err != nil {
		return "", err
	}

	content, err := os.ReadFile(fileInput.FilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

var DirectoryListerTool = ToolDefinition{
	ToolName:        "directory_lister",
	ToolDescription: "Lists files and directories at a specified path.",
	InputFormat:     DirectoryListerInputFormat,
	Handler:         ListDirectoryContents,
}

type DirectoryListerInput struct {
	DirectoryPath string `json:"directory_path,omitempty" jsonschema_description:"Path to the directory to list contents from."`
}

var DirectoryListerInputFormat = GenerateSchema[DirectoryListerInput]()

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

var FileEditorTool = ToolDefinition{
	ToolName: "file_editor",
	ToolDescription: `Edits a file by replacing specified text.

If the file does not exist, it will be created.`,
	InputFormat: FileEditorInputFormat,
	Handler:     EditFileContent,
}

type FileEditorInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path to the file to be edited."`
	OldText  string `json:"old_text" jsonschema_description:"Text to be replaced."`
	NewText  string `json:"new_text" jsonschema_description:"Text to replace with."`
}

var FileEditorInputFormat = GenerateSchema[FileEditorInput]()

func EditFileContent(input json.RawMessage) (string, error) {
	var editInput FileEditorInput
	if err := json.Unmarshal(input, &editInput); err != nil {
		return "", err
	}

	if editInput.FilePath == "" || editInput.OldText == editInput.NewText {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editInput.FilePath)
	if err != nil {
		if os.IsNotExist(err) && editInput.OldText == "" {
			return createNewFile(editInput.FilePath, editInput.NewText)
		}
		return "", err
	}

	updatedContent := strings.Replace(string(content), editInput.OldText, editInput.NewText, -1)
	if string(content) == updatedContent && editInput.OldText != "" {
		return "", fmt.Errorf("text to replace not found")
	}

	if err := os.WriteFile(editInput.FilePath, []byte(updatedContent), 0644); err != nil {
		return "", err
	}

	return "Edit successful", nil
}

var CommandRunnerTool = ToolDefinition{
	ToolName:        "command_runner",
	ToolDescription: "Executes a safe shell command and returns stdout/stderr.",
	InputFormat: GenerateSchema[struct {
		Command          string   `json:"command" jsonschema_description:"The shell command to execute."`
		Args             []string `json:"args,omitempty" jsonschema_description:"Optional arguments for the command."`
		WorkingDirectory string   `json:"working_directory,omitempty" jsonschema_description:"Optional working directory for the command."`
	}](),
	Handler: CommandRunnerHandler,
}