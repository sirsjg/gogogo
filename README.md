# ✨ GOGOGO Code Agent

GOGOGO is a lightweight command-line agentic code assistant written in Go built by enthusiasts, for enthusiasts. Currently Snazz only supports Claude, but we are investigating other LLM support depending on interest.

> [!NOTE]
> This project is currently under active development. Contributions are highly encouraged — join in on the fun!

## ✨ Features

- 💬 Chat with Claude AI using the Anthropics SDK.
- 🛠️ Dynamically invoke tools for file reading, directory listing, and file editing.
- 🧩 Extendable architecture for adding custom tools.

## 📋 Prerequisites

- 🔧 Go 1.24.3 or later installed.
- 🔑 An Anthropics API key. Set it as an environment variable:

```bash
export ANTHROPIC_API_KEY="sk-****************************"
```

## 📥 Installation

1. Clone the repository:

```bash
git clone https://github.com/sirsjg/snazz-go.git
cd snazz-go
```

2. Install dependencies:

```bash
go mod tidy
```

## 🚀 Run

To start the chatbot, run the following command:

```bash
go run main.go
```

## 🧰 Tools

Current, the agent supports the following basic tools:

| Tool Name           | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| **FileReaderTool**  | Reads the content of a file.                                               |
| **DirectoryListerTool** | Lists the contents of a directory.                                         |
| **FileEditorTool**  | Edits a file based on input.                                                |
| **CommandRunnerTool** | Executes a safe shell command and returns stdout/stderr.                    |
| **TextSearchTool**  | Searches for a string or pattern in files under a given directory.           |
| **DiffViewerTool** | Displays the diffs of changes for a list of files                             |

## 📝 To do

### Tools

- [x] File Reader Tool
- [x] Directory Reader Tool
- [x] File Editor Tool
- [x] Command Runner Tool
- [x] Text Search Tool

### Features

- [x] Shortcuts menu
- [ ] Syntax Highlighting
- [ ] Code diffs
- [ ] Tests
- [ ] Undo
- [ ] CI/CD
- [ ] Add MCP support
- [ ] CLI
- [ ] Brew
- [x] Token stats

## 🔧 Adding Tools

To add a new tool:

1. Create a new file for the tool in the `tools/` directory. For example, `tools/new_tool.go`.
2. Define the tool's structure and methods in the new file.
3. Implement the tool's logic in the same file, ensuring it adheres to the existing patterns and interfaces.
4. If the tool requires helper functions, add them to `tools/helpers.go`.
5. Register the tool in `main.go`.
6. Update the documentation in the `README.md` file to include the new tool.

## 📜 License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## 👥 Contributing

Contributions are welcome! Please open an issue or submit a pull request.
