package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("File created: %s", filePath), nil
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var instance T
	schema := reflector.Reflect(instance)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// CommandRunnerHandler executes a shell command and returns the output.
func CommandRunnerHandler(input json.RawMessage) (string, error) {
	var params struct {
		Command         string   `json:"command"`
		Args            []string `json:"args"`
		WorkingDirectory string   `json:"working_directory"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	cmd := exec.Command(params.Command, params.Args...)
	if params.WorkingDirectory != "" {
		cmd.Dir = params.WorkingDirectory
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}