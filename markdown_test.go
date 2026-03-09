package util

import (
	"testing"

	"github.com/slack-go/slack"
)

var (
	boolTrue = true
)

func TestConvertMarkdownTextToBlocks(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []slack.Block
	}{
		{
			name:     "heading",
			markdown: "# Title",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "Title",
						Emoji: &boolTrue,
					},
				},
			},
		},
		{
			name:     "paragraph",
			markdown: "This is a paragraph.",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is a paragraph.",
					},
				},
			},
		},
		{
			name:     "paragraph with bold",
			markdown: "This is **bold** text.",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is *bold* text.",
					},
				},
			},
		},
		{
			name:     "paragraph with italic",
			markdown: "This is *italic* text.",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is _italic_ text.",
					},
				},
			},
		},
		{
			name:     "paragraph with bold and italic 1",
			markdown: "This is **bold** and *italic* text.",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is *bold* and _italic_ text.",
					},
				},
			},
		},
		{
			name:     "paragraph with bold and italic 2",
			markdown: "This is **bold and *italic* mixed** text.",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is *bold and _italic_ mixed* text.",
					},
				},
			},
		},
		{
			name:     "unordered list",
			markdown: "- Item 1\n- Item 2",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Item 1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Item 2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "nested list",
			markdown: "- Parent 1\n  - Child 1\n  - Child 2\n- Parent 2",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Parent 1",
										},
									},
								},
							},
						},
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 1,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Child 1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Child 2",
										},
									},
								},
							},
						},
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Parent 2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "ordered list",
			markdown: "1. Item 1\n2. Item 2",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListOrdered,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Item 1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Item 2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "list with bold and italic",
			markdown: "- **bold** list item\n- *italic* list item",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "bold",
											Style: &slack.RichTextSectionTextStyle{
												Bold: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: " list item",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "italic",
											Style: &slack.RichTextSectionTextStyle{
												Italic: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: " list item",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "paragraph with hyperlink",
			markdown: "[hyperlink](https://example.com)",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "<https://example.com|hyperlink>",
					},
				},
			},
		},
		{
			name:     "multiple headings",
			markdown: "# Good morning\nGood morning\n# Good afternoon\nGood afternoon\n# Good evening\nGood evening",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "Good morning",
						Emoji: &boolTrue,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Good morning",
					},
				},
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "Good afternoon",
						Emoji: &boolTrue,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Good afternoon",
					},
				},
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "Good evening",
						Emoji: &boolTrue,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Good evening",
					},
				},
			},
		},
		{
			name:     "code block",
			markdown: "```\ncode block\n```",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextPreformatted{
							RichTextSection: slack.RichTextSection{
								Type: slack.RTEPreformatted,
								Elements: []slack.RichTextSectionElement{
									&slack.RichTextSectionTextElement{
										Type: slack.RTSEText,
										Text: "code block",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "multiline code block",
			markdown: "```\nfoo\nbar\n```",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextPreformatted{
							RichTextSection: slack.RichTextSection{
								Type: slack.RTEPreformatted,
								Elements: []slack.RichTextSectionElement{
									&slack.RichTextSectionTextElement{
										Type: slack.RTSEText,
										Text: "foo\nbar",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "code block with language specifier",
			markdown: "```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextPreformatted{
							RichTextSection: slack.RichTextSection{
								Type: slack.RTEPreformatted,
								Elements: []slack.RichTextSectionElement{
									&slack.RichTextSectionTextElement{
										Type: slack.RTSEText,
										Text: "func main() {\n    fmt.Println(\"hello\")\n}",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "inline code",
			markdown: "This is `inline code` text",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is `inline code` text",
					},
				},
			},
		},
		{
			name:     "blockquote",
			markdown: "> This is a quote",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextQuote{
							Type: slack.RTEQuote,
							Elements: []slack.RichTextSectionElement{
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "This is a quote",
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "multiple paragraphs",
			markdown: "Line 1\n\nLine 2\n\nLine 3",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Line 1",
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Line 2",
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Line 3",
					},
				},
			},
		},
		{
			name:     "complex markdown",
			markdown: "# Main Title\n\nThis is a paragraph with **bold** and *italic* text.\n\n- List item 1\n- Item with `code`\n- **Bold** item\n\n> *Quote* with **bold** can also be included\n\n[Link](https://example.com) and `inline code` in a paragraph.",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "Main Title",
						Emoji: &boolTrue,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "This is a paragraph with *bold* and _italic_ text.",
					},
				},
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:   slack.RTEList,
							Style:  slack.RTEListBullet,
							Indent: 0,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "List item 1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Item with ",
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "code",
											Style: &slack.RichTextSectionTextStyle{
												Code: true,
											},
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "Bold",
											Style: &slack.RichTextSectionTextStyle{
												Bold: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: " item",
										},
									},
								},
							},
						},
					},
				},
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextQuote{
							Type: slack.RTEQuote,
							Elements: []slack.RichTextSectionElement{
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "Quote",
									Style: &slack.RichTextSectionTextStyle{
										Italic: true,
									},
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: " with ",
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "bold",
									Style: &slack.RichTextSectionTextStyle{
										Bold: true,
									},
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: " can also be included",
								},
							},
						},
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "<https://example.com|Link> and `inline code` in a paragraph.",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMarkdownTextToBlocks(tt.markdown)
			if err != nil {
				t.Errorf("ConvertMarkdownTextToBlocks() returned error: %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ConvertMarkdownTextToBlocks() block count mismatch: got = %v, want %v", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].BlockType() != tt.want[i].BlockType() {
					t.Errorf("block type mismatch at index=%d, got=%v, want=%v", i, got[i].BlockType(), tt.want[i].BlockType())
				}

				if richTextBlock, ok := got[i].(*slack.RichTextBlock); ok {
					wantRichTextBlock := tt.want[i].(*slack.RichTextBlock)

					if len(richTextBlock.Elements) != len(wantRichTextBlock.Elements) {
						t.Errorf("RichTextBlock element count mismatch at index=%d, got=%v, want=%v",
							i, len(richTextBlock.Elements), len(wantRichTextBlock.Elements))
						continue
					}

					if list, ok := richTextBlock.Elements[0].(*slack.RichTextList); ok {
						wantList := wantRichTextBlock.Elements[0].(*slack.RichTextList)
						if list.Style != wantList.Style {
							t.Errorf("list style mismatch at index=%d, got=%v, want=%v",
								i, list.Style, wantList.Style)
						}

						if len(list.Elements) != len(wantList.Elements) {
							t.Errorf("list item count mismatch at index=%d, got=%v, want=%v",
								i, len(list.Elements), len(wantList.Elements))
							continue
						}

						for j, elem := range list.Elements {
							section := elem.(*slack.RichTextSection)
							wantSection := wantList.Elements[j].(*slack.RichTextSection)

							if len(section.Elements) != len(wantSection.Elements) {
								t.Errorf("section element count mismatch at index=%d,%d, got=%v, want=%v",
									i, j, len(section.Elements), len(wantSection.Elements))
								continue
							}

							textElem := section.Elements[0].(*slack.RichTextSectionTextElement)
							wantTextElem := wantSection.Elements[0].(*slack.RichTextSectionTextElement)
							if textElem.Text != wantTextElem.Text {
								t.Errorf("text content mismatch at index=%d,%d, got=%v, want=%v",
									i, j, textElem.Text, wantTextElem.Text)
							}
						}
					}
				}

				if headerBlock, ok := got[i].(*slack.HeaderBlock); ok {
					wantHeaderBlock := tt.want[i].(*slack.HeaderBlock)
					if headerBlock.Text.Text != wantHeaderBlock.Text.Text {
						t.Errorf("header text mismatch at index=%d, got=%v, want=%v",
							i, headerBlock.Text.Text, wantHeaderBlock.Text.Text)
					}
				}

				if sectionBlock, ok := got[i].(*slack.SectionBlock); ok {
					wantSectionBlock := tt.want[i].(*slack.SectionBlock)
					if sectionBlock.Text.Text != wantSectionBlock.Text.Text {
						t.Errorf("section text mismatch at index=%d, got=%v, want=%v",
							i, sectionBlock.Text.Text, wantSectionBlock.Text.Text)
					}
				}
			}
		})
	}
}
