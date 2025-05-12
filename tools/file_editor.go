package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ttacon/chalk"
)

type FileEditorInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path to the file to be edited."`
	OldText  string `json:"old_text" jsonschema_description:"Text to be replaced."`
	NewText  string `json:"new_text" jsonschema_description:"Text to replace with."`
}

var FileEditorInputFormat = GenerateSchema[FileEditorInput]()

var FileEditorTool = ToolDefinition{
	ToolName:        "file_editor",
	ToolDescription: `Edits a file by replacing specified text. If the file does not exist, it will be created and shown as a full addition.`,
	InputFormat:     FileEditorInputFormat,
	Handler:         EditFileContent,
}

func EditFileContent(input json.RawMessage) (string, error) {
	var editInput FileEditorInput
	if err := json.Unmarshal(input, &editInput); err != nil {
		return "", err
	}
	if editInput.FilePath == "" || editInput.OldText == editInput.NewText {
		return "", fmt.Errorf("invalid input parameters")
	}

	origBytes, err := os.ReadFile(editInput.FilePath)
	if err != nil {
		if os.IsNotExist(err) && editInput.OldText == "" {
			return createNewFileWithDiff(editInput.FilePath, editInput.NewText)
		}
		return "", err
	}
	original := string(origBytes)

	updated := strings.Replace(original, editInput.OldText, editInput.NewText, -1)
	if original == updated && editInput.OldText != "" {
		return "", fmt.Errorf("text to replace not found")
	}

	fmt.Println(chalk.Bold.TextStyle(editInput.FilePath))

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(original, updated, false)
	dmp.DiffCleanupSemantic(diffs)

	for _, d := range diffs {
		lines := strings.SplitAfter(d.Text, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			switch d.Type {
			case diffmatchpatch.DiffInsert:
				fmt.Print(chalk.Green.Color("+"+line), chalk.Reset)
			case diffmatchpatch.DiffDelete:
				fmt.Print(chalk.Red.Color("-"+line), chalk.Reset)
			case diffmatchpatch.DiffEqual:
				fmt.Print(chalk.White.Color(" "+line), chalk.Reset)
			}
		}
	}

	fmt.Println()

	if err := os.WriteFile(editInput.FilePath, []byte(updated), 0644); err != nil {
		return "", err
	}
	return "Edit successful", nil
}

func createNewFileWithDiff(path, content string) (string, error) {
	fmt.Println(chalk.Bold.TextStyle(path))

	for _, line := range strings.SplitAfter(content, "\n") {
		if line == "" {
			continue
		}
		fmt.Print(chalk.Green.Color("+"+line), chalk.Reset)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return "File created and diff shown", nil
}
