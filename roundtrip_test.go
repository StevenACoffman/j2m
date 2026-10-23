package j2m

import (
	"os"
	"testing"
)

func TestJiraRoundTrip(t *testing.T) {
	sample, err := os.ReadFile("j2m.jira")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		jira string
	}{
		{"headers", "h1. One\nh2. Two\nh3. Three\nh4. Four\nh5. Five\nh6. Six"},
		{"bold", "some *bold* text"},
		{"italic", "some _italic_ text"},
		{"bold and italic", "_*both*_"},
		{"monospaced", "{{mono}}"},
		{"monospaced with markup inside", "{{*x* _y_}}"},
		{"insert", "+inserted+"},
		{"superscript", "x^2^"},
		{"subscript", "H~2~O"},
		{"strikethrough", "a -deleted- b"},
		{"code block", "{code}\nx := 1\n{code}"},
		{"code block with language", "{code:go}\nx := 1\n{code}"},
		{"code block with markup inside", "{code}\nh1. *x* _y_\n# z\n{code}"},
		{"un-named link", "[http://example.com]"},
		{"named link", "[Example|http://example.com]"},
		{"image", "!http://example.com/i.png!"},
		{"image link", "[!http://example.com/i.png!|http://example.com]"},
		{"blockquote", "bq. quoted"},
		{"nested un-ordered list", "* a\n** b\n*** c\n* d"},
		{"nested ordered list", "# a\n## b\n### c\n# d"},
		{"formatting inside list", "* _i_ and *b*"},
		{"panel", "{panel:title=Note}\nbody\n{panel}\n"},
		{"table", "||H1||H2||\n|a|b|\n|c|d|"},
		{"table followed by text", "||H1||\n|a|\n|b|\n\nafter"},
		{"table after text", "before\n\n||H1||\n|a|\n|b|"},
		{"table after leading blank lines", "\n\n||H1||\n|a|\n|b|"},
		{"mixed document", "h1. Title\n\nSome *bold* and _italic_.\n\n* item\n** nested\n\n{code:go}\nx := 1\n{code}\n"},
		{"sample file", string(sample)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := JiraToMD(tt.jira)
			if got := MDToJira(md); got != tt.jira {
				t.Errorf("MDToJira(JiraToMD(%q)) = %q (via %q)", tt.jira, got, md)
			}
		})
	}
}

func TestMDRoundTrip(t *testing.T) {
	sample, err := os.ReadFile("j2m.md")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		md   string
	}{
		{"headers", "# One\n## Two\n### Three\n#### Four\n##### Five\n###### Six"},
		{"bold", "some **bold** text"},
		{"italic", "some *italic* text"},
		{"bold and italic", "***both***"},
		{"inline code", "`mono`"},
		{"inline code with markup inside", "`**x** _y_ <z>`"},
		{"strikethrough", "a ~~deleted~~ b"},
		{"ins tag", "<ins>inserted</ins>"},
		{"sup tag", "x<sup>2</sup>"},
		{"sub tag", "H<sub>2</sub>O"},
		{"code block", "```\nx := 1\n```"},
		{"code block with language", "```go\nx := 1\n```"},
		{"code block with markup inside", "```\n# *x* _y_\n> z\n1. q\n```"},
		{"un-named link", "<http://example.com>"},
		{"named link", "[Example](http://example.com)"},
		{"image", "![](http://example.com/i.png)"},
		{"image link", "[![](http://example.com/i.png)](http://example.com)"},
		{"blockquote", "> quoted"},
		{"nested un-ordered list", "* a\n  * b\n    * c\n* d"},
		{"nested ordered list", "1. a\n  1. b\n    1. c\n1. d"},
		{"formatting inside list", "* *i* and **b**"},
		{"panel", "| Note |\n| --- |\n| body |\n"},
		{"panel after blank line", "\n| Note |\n| --- |\n| body |\n"},
		{"table", "|H1|H2|\n| --- | --- |\n|a|b|\n|c|d|"},
		{"table after blank line", "\n|H1|H2|\n| --- | --- |\n|a|b|\n|c|d|"},
		{"table followed by text", "|H1|\n| --- |\n|a|\n|b|\n\nafter"},
		{"table after text", "before\n\n|H1|\n| --- |\n|a|\n|b|"},
		{"table after leading blank lines", "\n\n|H1|\n| --- |\n|a|\n|b|"},
		{"mixed document", "# Title\n\nSome **bold** and *italic*.\n\n* item\n  * nested\n\n```go\nx := 1\n```\n"},
		{"sample file", string(sample)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jira := MDToJira(tt.md)
			if got := JiraToMD(jira); got != tt.md {
				t.Errorf("JiraToMD(MDToJira(%q)) = %q (via %q)", tt.md, got, jira)
			}
		})
	}
}
