package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type DiffInput struct {
	Files []string `json:"files"`
}

func colorizeDiff(diff string) string {
	var coloredDiff strings.Builder
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+"):
			coloredDiff.WriteString("\033[32m" + line + "\033[0m\n") // Green for additions
		case strings.HasPrefix(line, "-"):
			coloredDiff.WriteString("\033[31m" + line + "\033[0m\n") // Red for deletions
		default:
			coloredDiff.WriteString(line + "\n")
		}
	}
	return coloredDiff.String()
}

var DiffViewerTool = ToolDefinition{
	ToolName:        "diff_viewer",
	ToolDescription: "Displays the diffs of changes for a list of files.",
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

		return colorizeDiff(output.String()), nil
	},
}
