package util

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var md = goldmark.New()

// ConvertMarkdownTextToBlocks converts a markdown text to a slice of slack blocks.
func ConvertMarkdownTextToBlocks(markdown string) ([]slack.Block, error) {
	source := []byte(markdown)
	doc := md.Parser().Parse(text.NewReader(source))
	blocks := []slack.Block{}

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.Kind() {
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			var text string
			lines := heading.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				text += string(line.Value(source))
			}
			emojiEnabled := true
			blocks = append(blocks, &slack.HeaderBlock{
				Type: slack.MBTHeader,
				Text: &slack.TextBlockObject{
					Type:  slack.PlainTextType,
					Text:  text,
					Emoji: &emojiEnabled,
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindParagraph:
			para := n.(*ast.Paragraph)
			var paraText string
			lines := para.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				paraText += string(line.Value(source))
			}
			mrkdwn := convertInlineMarkdownToMrkdwn(paraText)
			blocks = append(blocks, &slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: mrkdwn,
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindList:
			list := n.(*ast.List)
			// Collect all list items with their indent levels (handles nested lists)
			richTextElements := collectListItems(list, source, 0)
			blocks = append(blocks, &slack.RichTextBlock{
				Type:     slack.MBTRichText,
				Elements: richTextElements,
			})
			return ast.WalkSkipChildren, nil

		case ast.KindFencedCodeBlock:
			code := n.(*ast.FencedCodeBlock)
			var codeText string
			lines := code.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				codeText += string(line.Value(source))
			}

			blocks = append(blocks, &slack.RichTextBlock{
				Type: slack.MBTRichText,
				Elements: []slack.RichTextElement{
					&slack.RichTextPreformatted{
						RichTextSection: slack.RichTextSection{
							Type: slack.RTEPreformatted,
							Elements: []slack.RichTextSectionElement{
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: codeText,
								},
							},
						},
					},
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindCodeSpan:
			code := n.(*ast.CodeSpan)
			var text string
			lines := code.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				text += string(line.Value(source))
			}
			elements := []slack.RichTextSectionElement{
				&slack.RichTextSectionTextElement{
					Type: slack.RTSEText,
					Text: text,
					Style: &slack.RichTextSectionTextStyle{
						Code: true,
					},
				},
			}
			blocks = append(blocks, &slack.RichTextBlock{
				Type: slack.MBTRichText,
				Elements: []slack.RichTextElement{
					&slack.RichTextSection{
						Type:     slack.RTESection,
						Elements: elements,
					},
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindBlockquote:
			quote := n.(*ast.Blockquote)
			var quoteText string
			for child := quote.FirstChild(); child != nil; child = child.NextSibling() {
				if child.Kind() == ast.KindParagraph {
					lines := child.Lines()
					for i := 0; i < lines.Len(); i++ {
						line := lines.At(i)
						quoteText += string(line.Value(source))
					}
				}
			}
			blocks = append(blocks, &slack.RichTextBlock{
				Type: slack.MBTRichText,
				Elements: []slack.RichTextElement{
					&slack.RichTextQuote{
						Type: slack.RTEQuote,
						Elements: []slack.RichTextSectionElement{
							&slack.RichTextSectionTextElement{
								Type: slack.RTSEText,
								Text: quoteText,
							},
						},
					},
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindLink:
			link := n.(*ast.Link)
			var text string
			for c := link.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindText {
					textNode := c.(*ast.Text)
					text += string(textNode.Segment.Value(source))
				}
			}
			elements := []slack.RichTextSectionElement{
				&slack.RichTextSectionTextElement{
					Type: slack.RTSEText,
					Text: text,
				},
			}
			blocks = append(blocks, &slack.RichTextBlock{
				Type: slack.MBTRichText,
				Elements: []slack.RichTextElement{
					&slack.RichTextSection{
						Type:     slack.RTESection,
						Elements: elements,
					},
				},
			})
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// listItemWithIndent represents a list item with its indentation level and style
type listItemWithIndent struct {
	section *slack.RichTextSection
	indent  int
	style   slack.RichTextListElementType
}

// collectListItems recursively collects all list items from a list and its nested sublists,
// returning them as a flat slice of RichTextList elements with proper indent levels.
// This enables Slack's Block Kit to render nested lists correctly.
func collectListItems(list *ast.List, source []byte, indent int) []slack.RichTextElement {
	var items []listItemWithIndent
	style := getListStyle(list)

	for listItem := list.FirstChild(); listItem != nil; listItem = listItem.NextSibling() {
		if listItem.Kind() != ast.KindListItem {
			continue
		}

		// Process each child of the list item
		for child := listItem.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == ast.KindList {
				// Recursively collect nested list items
				nestedList := child.(*ast.List)
				nestedElements := collectListItemsFlat(nestedList, source, indent+1)
				items = append(items, nestedElements...)
			} else {
				// This is the content of the list item (paragraph, text, etc.)
				elements := parseInlineElements(child, source)
				if len(elements) > 0 {
					section := &slack.RichTextSection{
						Type:     slack.RTESection,
						Elements: elements,
					}
					items = append(items, listItemWithIndent{
						section: section,
						indent:  indent,
						style:   style,
					})
				}
			}
		}
	}

	// Convert flat items to grouped RichTextList elements
	return groupItemsByIndent(items)
}

// collectListItemsFlat is like collectListItems but returns listItemWithIndent for internal use
func collectListItemsFlat(list *ast.List, source []byte, indent int) []listItemWithIndent {
	var items []listItemWithIndent
	style := getListStyle(list)

	for listItem := list.FirstChild(); listItem != nil; listItem = listItem.NextSibling() {
		if listItem.Kind() != ast.KindListItem {
			continue
		}

		for child := listItem.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == ast.KindList {
				nestedList := child.(*ast.List)
				nestedElements := collectListItemsFlat(nestedList, source, indent+1)
				items = append(items, nestedElements...)
			} else {
				elements := parseInlineElements(child, source)
				if len(elements) > 0 {
					section := &slack.RichTextSection{
						Type:     slack.RTESection,
						Elements: elements,
					}
					items = append(items, listItemWithIndent{
						section: section,
						indent:  indent,
						style:   style,
					})
				}
			}
		}
	}

	return items
}

// groupItemsByIndent groups consecutive list items by their indent level and style,
// creating separate RichTextList elements for each group.
// This is required because Slack's Block Kit represents nested lists as separate
// RichTextList elements with incrementing indent values.
func groupItemsByIndent(items []listItemWithIndent) []slack.RichTextElement {
	if len(items) == 0 {
		return nil
	}

	var result []slack.RichTextElement
	var currentGroup []slack.RichTextElement
	currentIndent := items[0].indent
	currentStyle := items[0].style

	for _, item := range items {
		if item.indent != currentIndent || item.style != currentStyle {
			// Flush the current group
			if len(currentGroup) > 0 {
				result = append(result, &slack.RichTextList{
					Type:     slack.RTEList,
					Style:    currentStyle,
					Indent:   currentIndent,
					Elements: currentGroup,
				})
			}
			currentGroup = nil
			currentIndent = item.indent
			currentStyle = item.style
		}
		currentGroup = append(currentGroup, item.section)
	}

	// Flush the last group
	if len(currentGroup) > 0 {
		result = append(result, &slack.RichTextList{
			Type:     slack.RTEList,
			Style:    currentStyle,
			Indent:   currentIndent,
			Elements: currentGroup,
		})
	}

	return result
}

func getListStyle(list *ast.List) slack.RichTextListElementType {
	if list.IsOrdered() {
		return slack.RTEListOrdered
	}
	return slack.RTEListBullet
}

func parseInlineElements(n ast.Node, source []byte) []slack.RichTextSectionElement {
	var elements []slack.RichTextSectionElement
	var currentText string

	var process func(ast.Node, bool, bool)
	process = func(node ast.Node, isBold, isItalic bool) {
		if node == nil {
			return
		}

		switch node.Kind() {
		case ast.KindText:
			textNode := node.(*ast.Text)
			text := string(textNode.Segment.Value(source))
			if currentText != "" {
				elements = append(elements, &slack.RichTextSectionTextElement{
					Type: slack.RTSEText,
					Text: currentText,
				})
				currentText = ""
			}

			style := getTextStyle(isBold, isItalic)
			elements = append(elements, &slack.RichTextSectionTextElement{
				Type:  slack.RTSEText,
				Text:  text,
				Style: style,
			})

		case ast.KindEmphasis:
			emp := node.(*ast.Emphasis)
			newBold := isBold || emp.Level == 2
			newItalic := isItalic || emp.Level == 1
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				process(c, newBold, newItalic)
			}

		case ast.KindLink:
			link := node.(*ast.Link)
			var text string
			for c := link.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindText {
					textNode := c.(*ast.Text)
					text += string(textNode.Segment.Value(source))
				}
			}
			elements = append(elements, &slack.RichTextSectionLinkElement{
				Type: slack.RTSELink,
				Text: text,
				URL:  string(link.Destination),
			})

		case ast.KindCodeSpan:
			var text string
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindText {
					textNode := c.(*ast.Text)
					text += string(textNode.Segment.Value(source))
				}
			}
			if currentText != "" {
				elements = append(elements, &slack.RichTextSectionTextElement{
					Type: slack.RTSEText,
					Text: currentText,
				})
				currentText = ""
			}
			elements = append(elements, &slack.RichTextSectionTextElement{
				Type: slack.RTSEText,
				Text: text,
				Style: &slack.RichTextSectionTextStyle{
					Code: true,
				},
			})

		default:
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				process(c, isBold, isItalic)
			}
		}
	}

	process(n, false, false)

	if currentText != "" {
		elements = append(elements, &slack.RichTextSectionTextElement{
			Type: slack.RTSEText,
			Text: currentText,
		})
	}

	return elements
}

func getTextStyle(isBold, isItalic bool) *slack.RichTextSectionTextStyle {
	if !isBold && !isItalic {
		return nil
	}
	return &slack.RichTextSectionTextStyle{
		Bold:   isBold,
		Italic: isItalic,
	}
}

func convertInlineMarkdownToMrkdwn(markdown string) string {
	source := []byte(markdown)
	doc := md.Parser().Parse(text.NewReader(source))
	var result string

	var processNode func(ast.Node)
	processNode = func(n ast.Node) {
		if n == nil {
			return
		}

		switch n.Kind() {
		case ast.KindText:
			textNode := n.(*ast.Text)
			result += string(textNode.Segment.Value(source))

		case ast.KindEmphasis:
			emp := n.(*ast.Emphasis)
			switch emp.Level {
			case 2:
				result += "*"
				for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					processNode(c)
				}
				result += "*"
			case 1:
				result += "_"
				for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					processNode(c)
				}
				result += "_"
			}
			return

		case ast.KindLink:
			link := n.(*ast.Link)
			var text string
			for c := link.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindText {
					textNode := c.(*ast.Text)
					text += string(textNode.Segment.Value(source))
				}
			}
			result += fmt.Sprintf("<%s|%s>", string(link.Destination), text)
			return

		case ast.KindCodeSpan:
			var text string
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if c.Kind() == ast.KindText {
					textNode := c.(*ast.Text)
					text += string(textNode.Segment.Value(source))
				}
			}
			result += "`" + text + "`"
			return

		default:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				processNode(c)
			}
		}
	}

	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindParagraph {
			for child := c.FirstChild(); child != nil; child = child.NextSibling() {
				processNode(child)
			}
		}
	}

	return result
}
