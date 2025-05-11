package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"agent/tools"

	"github.com/anthropics/anthropic-sdk-go"
)

type Bot struct {
	apiClient *anthropic.Client
	readInput func() (string, bool)
	tools     []tools.ToolDefinition
	systemPrompt string
}

func InitializeBot(
	apiClient *anthropic.Client,
	readInput func() (string, bool),
	tools []tools.ToolDefinition,
	systemPrompt string,
) *Bot {
	return &Bot{
		apiClient: apiClient,
		readInput: readInput,
		tools:     tools,
		systemPrompt: systemPrompt,
	}
}

func (b *Bot) Execute(ctx context.Context) error {
	dialogue := []anthropic.MessageParam{}

	banner, err := os.ReadFile("./content/banner.txt")
	if err == nil {
		fmt.Println(string(banner))
	}

	acceptInput := true
	for {
		if acceptInput {
			fmt.Print("\033[1;31mYou:\033[0m ")
			userMessage, ok := b.readInput()
			if !ok {
				break
			}

			dialogue = append(dialogue, anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)))
		}

		response, err := b.processRequest(ctx, dialogue)
		if err != nil {
			return err
		}
		dialogue = append(dialogue, response.ToParam())

		toolResponses := []anthropic.ContentBlockParamUnion{}
		for _, content := range response.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\033[1;92mAgent:\033[0m %s\n", content.Text)
			case "tool_use":
				result := b.invokeTool(content.ID, content.Name, content.Input)
				toolResponses = append(toolResponses, result)
			}
		}
		if len(toolResponses) == 0 {
			acceptInput = true
			continue
		}
		acceptInput = false
		dialogue = append(dialogue, anthropic.NewUserMessage(toolResponses...))
	}

	return nil
}

func (b *Bot) invokeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var tool tools.ToolDefinition
	var found bool
	for _, t := range b.tools {
		if t.ToolName == name {
			tool = t
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	fmt.Printf("\033[1;35mTool: %s(%s)\033[0m\n", name, input)
	result, err := tool.Handler(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, result, false)
}

func LoadSystemPrompt(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return string("")
	}
	return string(content)
}


func (b *Bot) processRequest(ctx context.Context, dialogue []anthropic.MessageParam) (*anthropic.Message, error) {
	toolParams := []anthropic.ToolUnionParam{}
	for _, t := range b.tools {
		toolParams = append(toolParams, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        t.ToolName,
				Description: anthropic.String(t.ToolDescription),
				InputSchema: t.InputFormat,
			},
		})
	}

	return b.apiClient.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  dialogue,
		Tools:     toolParams,
		System: []anthropic.TextBlockParam{
			{Text: b.systemPrompt},
		},
	})
}