package tools

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
)

type ToolDefinition struct {
	ToolName        string                         `json:"tool_name"`
	ToolDescription string                         `json:"tool_description"`
	InputFormat     anthropic.ToolInputSchemaParam `json:"input_format"`
	Handler         func(input json.RawMessage) (string, error)
}