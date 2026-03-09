# 🚀 Slack API in Go Utility

[![Go Reference](https://pkg.go.dev/badge/github.com/takara2314/slack-go-util.svg)](https://pkg.go.dev/github.com/takara2314/slack-go-util)
[![CI](https://github.com/takara2314/slack-go-util/actions/workflows/test.yaml/badge.svg)](https://github.com/takara2314/slack-go-util/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/release/takara2314/slack-go-util.svg?style=flat-square)](https://github.com/takara2314/slack-go-util/releases)

A powerful utility package for slack-go that converts Markdown text into Slack's Block Kit format with ease! ✨

## ✨ Features
- 🔄 Convert Markdown to Slack Blocks
- 📚 Support for multiple Markdown elements:
    - Headers
    - Paragraphs
    - Lists
    - Nested lists
    - Code blocks
    - Blockquotes

## 📦 Installation
Install using Go Modules:

```sh
go get github.com/takara2314/slack-go-util
```

## 🚀 Usage
Here's a basic example of converting Markdown text to Slack Blocks:

```go
package main

import (
	"github.com/slack-go/slack"
	slackUtil "github.com/takara2314/slack-go-util"
)

func main() {
	// Initialize Slack API client
	api := slack.New("your-slack-token")

	markdown := `# Today's Tasks
**Project Updates**
- 🎯 Completed user authentication feature
  - Login flow implemented
  - Session management added
- 🐛 Fixed database connection issues
- 📱 Updated mobile responsive design

**Todo**
1. Review pull requests
2. Update documentation
3. Deploy to staging

> Don't forget team meeting at 2pm!`

	// Convert Markdown to Slack Blocks
	blocks, err := slackUtil.ConvertMarkdownTextToBlocks(markdown)
	if err != nil {
		panic(err)
	}

	// Send message to Slack
	_, _, err = api.PostMessage(
		"CHANNEL_ID",
		slack.MsgOptionBlocks(blocks...),
	)
	if err != nil {
		panic(err)
	}
}

```

The above code will send a beautifully formatted message to your Slack channel, including both bulleted and numbered lists! 📝

## 👥 Contributing
Contributions are welcome! 🎉 Feel free to:

- Report bugs
- Request features
- Submit pull requests

Please check out our CONTRIBUTING.md for guidelines.

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Support
If you have any questions or need support, please open an issue in the GitHub repository.
