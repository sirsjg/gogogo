package bot

import (
	"fmt"

	"github.com/ttacon/chalk"
)

func DisplayShortcuts() {
	fmt.Println(chalk.Yellow.Color(chalk.Bold.TextStyle("Shortcuts:")))
	fmt.Println()
	fmt.Println(chalk.Dim.TextStyle("/clear   - clear history"))
	fmt.Println(chalk.Dim.TextStyle("/tokens  - show/hide token usage"))
	fmt.Println(chalk.Dim.TextStyle("/exit    - leave"))
	fmt.Println()
}

func Write(message string, color chalk.Color, linebreak ...bool) {
	if len(linebreak) > 0 && !linebreak[0] {
		fmt.Print(color.Color(message))
	} else {
		fmt.Println(color.Color(message))
	}
}

func WriteError(message string) {
	fmt.Println(chalk.Red.Color(fmt.Sprintf("Error: %s", message)))
}