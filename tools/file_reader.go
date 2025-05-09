package tools

import (
	"encoding/json"
	"os"
)

type FileReaderInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path to the file to be read."`
}

var FileReaderInputFormat = GenerateSchema[FileReaderInput]()

var FileReaderTool = ToolDefinition{
	ToolName:        "file_reader",
	ToolDescription: "Reads the content of a specified file.",
	InputFormat:     FileReaderInputFormat,
	Handler:         ReadFileContent,
}

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