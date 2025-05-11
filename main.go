package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"agent/bot"
	"agent/tools"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	apiClient := anthropic.NewClient()

	inputScanner := bufio.NewScanner(os.Stdin)
	readInput := func() (string, bool) {
		if !inputScanner.Scan() {
			return "", false
		}
		return inputScanner.Text(), true
	}

	tools := []tools.ToolDefinition{
		tools.FileReaderTool,
		tools.DirectoryListerTool,
		tools.FileEditorTool,
		tools.CommandRunnerTool,
		tools.TextSearchTool,
	}

	systemPrompt := bot.LoadSystemPrompt("content/prompts/system.md")

	chatBot := bot.InitializeBot(&apiClient, readInput, tools, systemPrompt)

	if err := chatBot.Execute(context.TODO()); err != nil {
		fmt.Printf("Execution Error: %s\n", err.Error())
	}
}