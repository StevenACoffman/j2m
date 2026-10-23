package j2m

import "testing"

func TestMDToJira(t *testing.T) {
	tests := []struct {
		name string
		md   string
		want string
	}{
		{"h1", "# Title", "h1. Title"},
		{"h2", "## Title", "h2. Title"},
		{"h3", "### Title", "h3. Title"},
		{"h4", "#### Title", "h4. Title"},
		{"h5", "##### Title", "h5. Title"},
		{"h6", "###### Title", "h6. Title"},
		{"underlined h1", "Title\n=====", "h1. Title"},
		{"underlined h2", "Title\n-----", "h2. Title"},
		{"bold", "**bold**", "*bold*"},
		{"bold with underscores", "__bold__", "*bold*"},
		{"italic", "*italic*", "_italic_"},
		{"italic with underscores", "_italic_", "_italic_"},
		{"bold and italic", "***both***", "_*both*_"},
		{"bold and italic with underscores", "___both___", "_*both*_"},
		{"mixed emphasis on a line", "*i* **b** ***bi***", "_i_ *b* _*bi*_"},
		{"inline code", "`mono`", "{{mono}}"},
		{"inline code content is untouched", "`**not bold** <x>`", "{{**not bold** <x>}}"},
		{"strikethrough", "a ~~deleted~~ b", "a -deleted- b"},
		{"del tag", "<del>deleted</del>", "-deleted-"},
		{"ins tag", "<ins>inserted</ins>", "+inserted+"},
		{"sup tag", "x<sup>2</sup>", "x^2^"},
		{"sub tag", "H<sub>2</sub>O", "H~2~O"},
		{"code block", "```\nx := 1\n```", "{code}\nx := 1\n{code}"},
		{"code block with language", "```go\nx := 1\n```", "{code:go}\nx := 1\n{code}"},
		{"code block content is untouched", "```\n# *x* > [y](z)\n1. q\n```", "{code}\n# *x* > [y](z)\n1. q\n{code}"},
		{"multiple code blocks", "```\na\n```\n**b**\n```\nc\n```", "{code}\na\n{code}\n*b*\n{code}\nc\n{code}"},
		{"un-named link", "<http://example.com>", "[http://example.com]"},
		{"named link", "[Example](http://example.com)", "[Example|http://example.com]"},
		{"image", "![alt](http://example.com/i.png)", "!http://example.com/i.png!"},
		{"image link", "[![](http://example.com/i.png)](http://example.com)", "[!http://example.com/i.png!|http://example.com]"},
		{"blockquote", "> quoted", "bq. quoted"},
		{"un-ordered list", "* a\n* b", "* a\n* b"},
		{"nested un-ordered list", "* a\n  * b\n    * c", "* a\n** b\n*** c"},
		{"ordered list", "1. a\n2. b", "# a\n# b"},
		{"nested ordered list", "1. a\n  1. b\n    1. c", "# a\n## b\n### c"},
		{"formatting inside list", "* *i* and **b**", "* _i_ and *b*"},
		{"panel", "| Note |\n| --- |\n| body |\n", "{panel:title=Note}\nbody\n{panel}\n"},
		{"table", "| H1 | H2 |\n| --- | --- |\n| a | b |\n", "|| H1 || H2 ||\n| a | b |\n"},
		{"table at start of input", "| H1 |\n| --- |\n| a |\n| b |\n", "|| H1 ||\n| a |\n| b |\n"},
		{"leading blank line before table is kept", "\n| H1 |\n| --- |\n| a |\n| b |\n", "\n|| H1 ||\n| a |\n| b |\n"},
		{"blank line before table after text is dropped", "text\n\n| H1 |\n| --- |\n| a |\n| b |\n", "text\n|| H1 ||\n| a |\n| b |\n"},
		{"panel at start of input", "| Note |\n| --- |\n| body |\n", "{panel:title=Note}\nbody\n{panel}\n"},
		{"table after text without blank line is left alone", "text\n| H1 |\n| --- |\n| a |\n", "text\n| H1 |\n| --- |\n| a |\n"},
		{"table at end of input", "| H1 |\n| --- |\n| a |\n| b |", "|| H1 ||\n| a |\n| b |"},
		{"table with mismatched separator is left alone", "\n|H1|H2|\n| --- |\n|a|b|\n", "\n|H1|H2|\n| --- |\n|a|b|\n"},
		{"plain text", "nothing to see here", "nothing to see here"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MDToJira(tt.md); got != tt.want {
				t.Errorf("MDToJira(%q) = %q, want %q", tt.md, got, tt.want)
			}
		})
	}
}
