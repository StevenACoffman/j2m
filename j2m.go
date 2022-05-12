package j2m

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type jiration struct {
	re   *regexp.Regexp
	repl interface{}
}

// JiraToMD takes a string in Jira Markdown, and outputs Github Markdown
func JiraToMD(str string) string {
	jirations := []jiration{
		{ // UnOrdered Lists
			re: regexp.MustCompile(`(?m)^[ \t]*(\*+)\s+`),
			repl: func(groups []string) string {
				_, stars := groups[0], groups[1]
				return strings.Repeat("  ", len(stars)-1) + "* "
			},
		},
		{ //Ordered Lists
			re: regexp.MustCompile(`(?m)^[ \t]*(#+)\s+`),
			repl: func(groups []string) string {
				_, nums := groups[0], groups[1]
				return strings.Repeat("  ", len(nums)-1) + "1. "
			},
		},
		{ //Headers 1-6
			re: regexp.MustCompile(`(?m)^h([0-6])\.(.*)$`),
			repl: func(groups []string) string {
				_, level, content := groups[0], groups[1], groups[2]
				i, _ := strconv.Atoi(level)
				return strings.Repeat("#", i) + content
			},
		},
		{ // Bold
			re:   regexp.MustCompile(`\*(\S[^*]*)\*`),
			repl: "**$1**",
		},
		{ // Italic
			re:   regexp.MustCompile(`\b\_(\S[^_]*)\_`),
			repl: "*$1*",
		},
		{ // Monospaced text
			re:   regexp.MustCompile(`\{\{([^}]+)\}\}`),
			repl: "`$1`",
		},
		{ // Citations (buggy)
			re:   regexp.MustCompile(`\?\?((?:.[^?]|[^?].)+)\?\?`),
			repl: "<cite>$1</cite>",
		},
		{ // Inserts
			re:   regexp.MustCompile(`\+([^+]*)\+`),
			repl: "<ins>$1</ins>",
		},
		{ // Superscript
			re:   regexp.MustCompile(`\^([^^]*)\^`),
			repl: "<sup>$1</sup>",
		},
		{ // Subscript
			re:   regexp.MustCompile(`~([^~]*)~`),
			repl: "<sub>$1</sub>",
		},
		{ // Strikethrough
			re:   regexp.MustCompile(`(\s+)-(\S+.*?\S)-(\s+)`),
			repl: "$1~~$2~~$3",
		},
		{ // Code Block
			re:   regexp.MustCompile(`\{code(:([a-z]+))?([:|]?(title|borderStyle|borderColor|borderWidth|bgColor|titleBGColor)=.+?)*\}`),
			repl: "```$2",
		},
		{ // Code Block End
			re:   regexp.MustCompile(`{code}`),
			repl: "```",
		},
		{ // Pre-formatted text
			re:   regexp.MustCompile(`{noformat}`),
			repl: "```",
		},
		{ // Un-named Links
			re:   regexp.MustCompile(`(?U)\[([^|]+)\]`),
			repl: "<$1>",
		},
		{ // Images
			re:   regexp.MustCompile(`!(.+)!`),
			repl: "![]($1)",
		},
		{ // Named Links
			re:   regexp.MustCompile(`\[(.+?)\|(.+)\]`),
			repl: "[$1]($2)",
		},
		{ // Single Paragraph Blockquote
			re:   regexp.MustCompile(`(?m)^bq\.\s+`),
			repl: "> ",
		},
		{ // Remove color: unsupported in md
			re:   regexp.MustCompile(`(?m)\{color:[^}]+\}(.*)\{color\}`),
			repl: "$1",
		},
		{ // panel into table
			re:   regexp.MustCompile(`(?m)\{panel:title=([^}]*)\}\n?(.*?)\n?\{panel\}`),
			repl: "\n| $1 |\n| --- |\n| $2 |",
		},
		{ //table header
			re: regexp.MustCompile(`(?m)^[ \t]*((?:\|\|.*?)+\|\|)[ \t]*$`),
			repl: func(groups []string) string {
				_, headers := groups[0], groups[1]
				reBarred := regexp.MustCompile(`\|\|`)

				singleBarred := reBarred.ReplaceAllString(headers, "|")
				fillerRe := regexp.MustCompile(`\|[^|]+`)
				return "\n" + singleBarred + "\n" + fillerRe.ReplaceAllString(singleBarred, "| --- ")
			},
		},
		{ // remove leading-space of table headers and rows
			re:   regexp.MustCompile(`(?m)^[ \t]*\|`),
			repl: "|",
		},
	}
	for _, jiration := range jirations {
		switch v := jiration.repl.(type) {
		case string:
			str = jiration.re.ReplaceAllString(str, v)
		case func([]string) string:
			str = replaceAllStringSubmatchFunc(jiration.re, str, v)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	return str
}

// https://gist.github.com/elliotchance/d419395aa776d632d897
func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
