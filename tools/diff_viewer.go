package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/anthropics/anthropic-sdk-go"
)

type DiffInput struct {
	Files []string `json:"files"`
}

var DiffViewerTool = ToolDefinition{
	ToolName:        "diff_viewer",
	ToolDescription: "Displays the diffs of changes to a file.",
	InputFormat:     anthropic.ToolInputSchemaParam{Type: "object", Properties: map[string]anthropic.ToolInputSchemaParam{"files": {Type: "array"}}},
	Handler: func(input json.RawMessage) (string, error) {
		var diffInput DiffInput
		if err := json.Unmarshal(input, &diffInput); err != nil {
			return "", fmt.Errorf("invalid input format: %w", err)
		}

		var output bytes.Buffer
		for _, file := range diffInput.Files {
			cmd := exec.Command("git", "diff", file)
			cmd.Stdout = &output
			cmd.Stderr = &output
			if err := cmd.Run(); err != nil {
				return "", fmt.Errorf("error running git diff on %s: %w", file, err)
			}
		}

		return output.String(), nil
	},
}
