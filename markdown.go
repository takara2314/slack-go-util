// markdown.go
package util

import (
	"github.com/slack-go/slack"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var md = goldmark.New()

func ConvertMarkdownTextToBlocks(markdown string) []slack.Block {
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
			blocks = append(blocks, &slack.HeaderBlock{
				Type: slack.MBTHeader,
				Text: &slack.TextBlockObject{
					Type:  slack.PlainTextType,
					Text:  text,
					Emoji: true,
				},
			})
			return ast.WalkSkipChildren, nil

		case ast.KindParagraph:
			elements := parseInlineElements(n, source)
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

		case ast.KindList:
			list := n.(*ast.List)
			listItems := []slack.RichTextElement{}

			for listItem := n.FirstChild(); listItem != nil; listItem = listItem.NextSibling() {
				if listItem.Kind() != ast.KindListItem {
					continue
				}

				elements := parseInlineElements(listItem.FirstChild(), source)
				section := &slack.RichTextSection{
					Type:     slack.RTESection,
					Elements: elements,
				}
				listItems = append(listItems, section)
			}

			blocks = append(blocks, &slack.RichTextBlock{
				Type: slack.MBTRichText,
				Elements: []slack.RichTextElement{
					&slack.RichTextList{
						Type:     slack.RTEList,
						Style:    getListStyle(list),
						Elements: listItems,
					},
				},
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
					&slack.RichTextSection{
						Type: slack.RTESection,
						Elements: []slack.RichTextSectionElement{
							&slack.RichTextSectionTextElement{
								Type: slack.RTSEText,
								Text: codeText,
								Style: &slack.RichTextSectionTextStyle{
									Code: true,
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
					text += string(c.Text(source))
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
		panic(err)
	}

	return blocks
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
					text += string(c.Text(source))
				}
			}
			elements = append(elements, &slack.RichTextSectionLinkElement{
				Type: slack.RTSELink,
				Text: text,
				URL:  string(link.Destination),
			})

		case ast.KindCodeSpan:
			text := string(node.Text(source))
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
