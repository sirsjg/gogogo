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
		fmt.Println()
	}

	acceptInput := true
	showTokenCount := false
	
	for {
		if acceptInput {
			fmt.Print("\033[1;31m>\033[0m ")
			userMessage, ok := b.readInput()
			if !ok {
				break
			}
			fmt.Println()

			// Check for shortcuts
			if len(userMessage) > 0 && (userMessage[0] == '?' || userMessage[0] == '/') {
				switch userMessage {
				case "/clear":
					dialogue = []anthropic.MessageParam{} // Clear history
					fmt.Println("\033[1;32mHistory cleared.\033[0m")
				case "/tokens":
					showTokenCount = !showTokenCount // Toggle token count display
					status := "hidden"
					if showTokenCount {
						status = "shown"
					}
					fmt.Printf("\033[1;32mToken count is now %s.\033[0m\n", status)
				case "/exit":
					fmt.Println("\033[1;32mExiting...\033[0m")
					return nil
				default:
					fmt.Println("\033[1;33mShortcuts:\033[0m")
					fmt.Println("\033[1;90m")
					fmt.Println("/clear   - clear history")
					fmt.Println("/tokens  - show/hide token count")
					fmt.Println("/exit    - leave")
					fmt.Printf("\033[0m")
				}
				fmt.Println() 
				continue
			}
			dialogue = append(dialogue, anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)))
		}

		response, err := b.processRequest(ctx, dialogue)
		if err != nil {
			return err
		}

		// Display raw JSON of response.Usage
		usageJSON, err := json.MarshalIndent(response.Usage, "", "  ")
		if err == nil {
			fmt.Printf("\033[1;34mRaw Usage JSON:\033[0m\n%s\n", string(usageJSON))
		}

		inputTokens := response.Usage.InputTokens   // Track input tokens
		outputTokens := response.Usage.OutputTokens // Track output tokens

		if showTokenCount {
			fmt.Printf("\033[1;33mToken Usage: Input: %d, Output: %d\033[0m\n", inputTokens, outputTokens) // Display token usage in yellow
		}

		dialogue = append(dialogue, response.ToParam())

		toolResponses := []anthropic.ContentBlockParamUnion{}
		for _, content := range response.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\033[1;90m%s\033[0m\n", content.Text) 
				fmt.Println()
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

	fmt.Printf("\033[1;35mTool:\033[0m %s(%s)\n", name, input)
	fmt.Println() 
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