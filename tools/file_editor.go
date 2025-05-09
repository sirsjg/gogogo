package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type FileEditorInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path to the file to be edited."`
	OldText  string `json:"old_text" jsonschema_description:"Text to be replaced."`
	NewText  string `json:"new_text" jsonschema_description:"Text to replace with."`
}

var FileEditorInputFormat = GenerateSchema[FileEditorInput]()

var FileEditorTool = ToolDefinition{
	ToolName: "file_editor",
	ToolDescription: `Edits a file by replacing specified text.

If the file does not exist, it will be created.`,
	InputFormat: FileEditorInputFormat,
	Handler:     EditFileContent,
}

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