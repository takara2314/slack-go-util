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
