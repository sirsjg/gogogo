package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

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