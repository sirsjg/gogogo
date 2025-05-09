# ✨ Sizzle Code Agent

Snazzy is a lightweight command-line coding agent built by enthusiasts, for enthusiasts. Whether you're navigating through file reading, directory listing, or quick file edits, this trusty bot is here to make your coding experience smooth and efficient.

**Note:** This project is currently under active development. Contributions are highly encouraged—join us in making Snazzy even better!

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
git clone https://github.com/sirsjg/snazzy.git
cd snazzy
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

The application supports the following tools:

1. **FileReaderTool**: Reads the content of a file.
2. **DirectoryListerTool**: Lists the contents of a directory.
3. **FileEditorTool**: Edits a file based on input.
4. **CommandRunnerTool**: Executes a safe shell command and returns stdout/stderr.

## 📝 To do

- [x] File Reader Tool
- [x] Directory Reader Tool
- [x] File Editor Tool
- [x] Command Runner Tool
- [ ] Text Search Tool
- [ ] Code Formatter
- [ ] Tests
- [ ] CI/CD
- [ ] Add MCP support

## 🔧 Adding Tools

To add a new tool:

1. Define the tool in `tools/tools.go`.
2. Implement the tool handler in `tools/helpers.go`.
3. Register the tool in `main.go`.

## 📜 License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## 👥 Contributing

Contributions are welcome! Please open an issue or submit a pull request.
