package util

import (
	"testing"

	"github.com/slack-go/slack"
)

func TestConvertMarkdownTextToBlocks(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []slack.Block
	}{
		{
			name:     "見出しのテスト",
			markdown: "# タイトル",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "タイトル",
						Emoji: true,
					},
				},
			},
		},
		{
			name:     "段落のテスト",
			markdown: "これは段落です。",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは段落です。",
					},
				},
			},
		},
		{
			name:     "太文字混合の段落のテスト",
			markdown: "これは**太文字**です。",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは*太文字*です。",
					},
				},
			},
		},
		{
			name:     "斜体混合の段落のテスト",
			markdown: "これは*斜体*です。",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは_斜体_です。",
					},
				},
			},
		},
		{
			name:     "太文字と斜体混合の段落のテスト1",
			markdown: "これは**太文字**と*斜体*です。",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは*太文字*と_斜体_です。",
					},
				},
			},
		},
		{
			name:     "太文字と斜体混合の段落のテスト2",
			markdown: "これは**太文字の中の *斜体* **です。",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは*太文字の中の _斜体_ *です。",
					},
				},
			},
		},
		{
			name:     "リストのテスト",
			markdown: "- アイテム1\n- アイテム2",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:  slack.RTEList,
							Style: slack.RTEListBullet,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "アイテム1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "アイテム2",
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
			name:     "番号付きリストのテスト",
			markdown: "1. アイテム1\n2. アイテム2",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:  slack.RTEList,
							Style: slack.RTEListOrdered,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "アイテム1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "アイテム2",
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
			name:     "太文字と斜体が加わったリストのテスト",
			markdown: "- **太文字**のリスト\n- *斜体*のリスト",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:  slack.RTEList,
							Style: slack.RTEListBullet,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "太文字",
											Style: &slack.RichTextSectionTextStyle{
												Bold: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "のリスト",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "斜体",
											Style: &slack.RichTextSectionTextStyle{
												Italic: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "のリスト",
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
			name:     "ハイパーリンクが含まれた段落のテスト",
			markdown: "[ハイパーリンク](https://example.com)",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "<https://example.com|ハイパーリンク>",
					},
				},
			},
		},
		{
			name:     "複数の見出しのテスト",
			markdown: "# おはようございます\nおはようございます\n# こんにちは\nこんにちは\n# こんばんは\nこんばんは",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "おはようございます",
						Emoji: true,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "おはようございます",
					},
				},
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "こんにちは",
						Emoji: true,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "こんにちは",
					},
				},
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "こんばんは",
						Emoji: true,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "こんばんは",
					},
				},
			},
		},
		{
			name:     "コードブロックのテスト",
			markdown: "```\ncode block\n```",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextSection{
							Type: slack.RTESection,
							Elements: []slack.RichTextSectionElement{
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "code block",
									Style: &slack.RichTextSectionTextStyle{
										Code: true,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "インラインコードのテスト",
			markdown: "これは`インラインコード`です",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは`インラインコード`です",
					},
				},
			},
		},
		{
			name:     "引用のテスト",
			markdown: "> 引用文です",
			want: []slack.Block{
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextQuote{
							Type: slack.RTEQuote,
							Elements: []slack.RichTextSectionElement{
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "引用文です",
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "複数行の段落のテスト",
			markdown: "1行目\n2行目\n3行目",
			want: []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "1段目",
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "2段目",
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "3段目",
					},
				},
			},
		},
		{
			name:     "複合的なマークダウンのテスト",
			markdown: "# メインタイトル\n\nこれは**太字**と*斜体*が含まれた段落です。\n\n- リストアイテム1\n- `コード`を含むアイテム\n- **太字**なアイテム\n\n> *引用文*の中に**太字**を入れることもできます\n\n[リンク](https://example.com)を含む`インラインコード`な段落です。",
			want: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type:  slack.PlainTextType,
						Text:  "メインタイトル",
						Emoji: true,
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "これは**太字**と*斜体*が含まれた段落です。",
					},
				},
				&slack.RichTextBlock{
					Type: slack.MBTRichText,
					Elements: []slack.RichTextElement{
						&slack.RichTextList{
							Type:  slack.RTEList,
							Style: slack.RTEListBullet,
							Elements: []slack.RichTextElement{
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "リストアイテム1",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "コード",
											Style: &slack.RichTextSectionTextStyle{
												Code: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "を含むアイテム",
										},
									},
								},
								&slack.RichTextSection{
									Type: slack.RTESection,
									Elements: []slack.RichTextSectionElement{
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "太字",
											Style: &slack.RichTextSectionTextStyle{
												Bold: true,
											},
										},
										&slack.RichTextSectionTextElement{
											Type: slack.RTSEText,
											Text: "なアイテム",
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
									Text: "引用文",
									Style: &slack.RichTextSectionTextStyle{
										Italic: true,
									},
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "の中に",
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "太字",
									Style: &slack.RichTextSectionTextStyle{
										Bold: true,
									},
								},
								&slack.RichTextSectionTextElement{
									Type: slack.RTSEText,
									Text: "を入れることもできます",
								},
							},
						},
					},
				},
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "[リンク](https://example.com)を含む`インラインコード`な段落です。",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMarkdownTextToBlocks(tt.markdown)
			if err != nil {
				t.Errorf("ConvertMarkdownTextToBlocks() にエラーが発生しました: %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ConvertMarkdownTextToBlocks() のブロック数が異なります got = %v, want %v", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].BlockType() != tt.want[i].BlockType() {
					t.Errorf("ブロックタイプが異なります index=%d, got=%v, want=%v", i, got[i].BlockType(), tt.want[i].BlockType())
				}

				// RichTextBlockの場合は、内部の要素も検証
				if richTextBlock, ok := got[i].(*slack.RichTextBlock); ok {
					wantRichTextBlock := tt.want[i].(*slack.RichTextBlock)

					if len(richTextBlock.Elements) != len(wantRichTextBlock.Elements) {
						t.Errorf("RichTextBlock の要素数が異なります index=%d, got=%v, want=%v",
							i, len(richTextBlock.Elements), len(wantRichTextBlock.Elements))
						continue
					}

					// リストの検証
					if list, ok := richTextBlock.Elements[0].(*slack.RichTextList); ok {
						wantList := wantRichTextBlock.Elements[0].(*slack.RichTextList)
						if list.Style != wantList.Style {
							t.Errorf("リストのスタイルが異なります index=%d, got=%v, want=%v",
								i, list.Style, wantList.Style)
						}

						if len(list.Elements) != len(wantList.Elements) {
							t.Errorf("リストアイテムの数が異なります index=%d, got=%v, want=%v",
								i, len(list.Elements), len(wantList.Elements))
							continue
						}

						// 各リストアイテムの検証
						for j, elem := range list.Elements {
							section := elem.(*slack.RichTextSection)
							wantSection := wantList.Elements[j].(*slack.RichTextSection)

							if len(section.Elements) != len(wantSection.Elements) {
								t.Errorf("セクション要素の数が異なります index=%d,%d, got=%v, want=%v",
									i, j, len(section.Elements), len(wantSection.Elements))
								continue
							}

							// テキスト要素の検証
							textElem := section.Elements[0].(*slack.RichTextSectionTextElement)
							wantTextElem := wantSection.Elements[0].(*slack.RichTextSectionTextElement)
							if textElem.Text != wantTextElem.Text {
								t.Errorf("テキスト内容が異なります index=%d,%d, got=%v, want=%v",
									i, j, textElem.Text, wantTextElem.Text)
							}
						}
					}
				}

				// HeaderBlockの場合は、テキスト内容を検証
				if headerBlock, ok := got[i].(*slack.HeaderBlock); ok {
					wantHeaderBlock := tt.want[i].(*slack.HeaderBlock)
					if headerBlock.Text.Text != wantHeaderBlock.Text.Text {
						t.Errorf("ヘッダーのテキストが異なります index=%d, got=%v, want=%v",
							i, headerBlock.Text.Text, wantHeaderBlock.Text.Text)
					}
				}

				// SectionBlockの場合は、テキスト内容を検証
				if sectionBlock, ok := got[i].(*slack.SectionBlock); ok {
					wantSectionBlock := tt.want[i].(*slack.SectionBlock)
					if sectionBlock.Text.Text != wantSectionBlock.Text.Text {
						t.Errorf("セクションのテキストが異なります index=%d, got=%v, want=%v",
							i, sectionBlock.Text.Text, wantSectionBlock.Text.Text)
					}
				}
			}
		})
	}
}
