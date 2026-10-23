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
	leadingNewlines := len(str) - len(strings.TrimLeft(str, "\n"))
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
		{ // Un-named Links
			re:   regexp.MustCompile(`(?U)\[([^|]+?)\]`),
			repl: "<$1>",
		},
		{ // Images
			re:   regexp.MustCompile(`!(.+)!`),
			repl: "![]($1)",
		},
		{ // Named Links
			re:   regexp.MustCompile(`\[(.+?)\|(.+?)\]`),
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

	// Code is converted up front and swapped for placeholders, so that the
	// jirations can't alter its contents
	codes := []string{}
	codeJirations := []jiration{
		{ // Code Block
			re: regexp.MustCompile(`\{code(:([a-z]+))?([:|]?(title|borderStyle|borderColor|borderWidth|bgColor|titleBGColor)=.+?)*\}((?:.|\n)*?)\{code\}`),
			repl: func(groups []string) string {
				_, lang, content := groups[0], groups[2], groups[5]
				codes = append(codes, "```"+lang+content+"```")
				return codePlaceholder(len(codes) - 1)
			},
		},
		{ // Pre-formatted text
			re: regexp.MustCompile(`\{noformat\}((?:.|\n)*?)\{noformat\}`),
			repl: func(groups []string) string {
				_, content := groups[0], groups[1]
				codes = append(codes, "```"+content+"```")
				return codePlaceholder(len(codes) - 1)
			},
		},
		{ // Monospaced text
			re: regexp.MustCompile(`\{\{([^}]+)\}\}`),
			repl: func(groups []string) string {
				_, content := groups[0], groups[1]
				codes = append(codes, "`"+content+"`")
				return codePlaceholder(len(codes) - 1)
			},
		},
	}
	for _, jiration := range append(codeJirations, jirations...) {
		switch v := jiration.repl.(type) {
		case string:
			str = jiration.re.ReplaceAllString(str, v)
		case func([]string) string:
			str = replaceAllStringSubmatchFunc(jiration.re, str, v)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	for i, code := range codes {
		str = strings.Replace(str, codePlaceholder(i), code, 1)
	}
	// Tables and panels get a blank line before them, which isn't wanted at
	// the start of the document
	if len(str)-len(strings.TrimLeft(str, "\n")) > leadingNewlines {
		str = strings.TrimPrefix(str, "\n")
	}
	return str
}

// MDToJira takes a string in Github Markdown, and outputs Jira Markdown
func MDToJira(str string) string {
	leadingNewlines := len(str) - len(strings.TrimLeft(str, "\n"))
	jirations := []jiration{
		{ // Tables
			re: regexp.MustCompile(`(?m)(?:\A|^\n)((?:\|.*?)+\|)[ \t]*\n((?:\|\s*?-{3,}\s*?)+\|)[ \t]*\n((?:(?:\|.*?)+\|[ \t]*(?:\n|\z))*)$`),
			repl: func(groups []string) string {
				match, headerLine, separatorLine, rowstr := groups[0], groups[1], groups[2], groups[3]
				headers := []string{}
				for _, header := range strings.Split(headerLine, "|") {
					if header != "" {
						headers = append(headers, header)
					}
				}
				separators := []string{}
				for _, separator := range strings.Split(separatorLine, "|") {
					if separator != "" {
						separators = append(separators, separator)
					}
				}
				if len(headers) != len(separators) {
					return match
				}

				rows := strings.Split(strings.TrimSuffix(rowstr, "\n"), "\n")
				if len(rows) == 1 && len(headers) == 1 {
					// Panel
					panelRe := regexp.MustCompile(`^\|(.*)[ \t]*\|`)
					return "{panel:title=" + strings.TrimSpace(headers[0]) + "}\n" +
						strings.TrimSpace(panelRe.ReplaceAllString(rowstr, "$1")) + "\n{panel}\n"
				}

				return "||" + strings.Join(headers, "||") + "||\n" + rowstr
			},
		},
		{ // Bold, Italic, and Combined (bold+italic)
			re: regexp.MustCompile(`\*\*\*(\S.*?)\*\*\*|___(\S.*?)___|\*\*(\S.*?)\*\*|__(\S.*?)__|\*(\S.*?)\*|_(\S.*?)_`),
			repl: func(groups []string) string {
				switch {
				case groups[1] != "":
					return "_*" + groups[1] + "*_"
				case groups[2] != "":
					return "_*" + groups[2] + "*_"
				case groups[3] != "":
					return "*" + groups[3] + "*"
				case groups[4] != "":
					return "*" + groups[4] + "*"
				case groups[5] != "":
					return "_" + groups[5] + "_"
				default:
					return "_" + groups[6] + "_"
				}
			},
		},
		{ // All Headers (# format)
			re: regexp.MustCompile(`(?m)^([#]+)(.*?)$`),
			repl: func(groups []string) string {
				_, level, content := groups[0], groups[1], groups[2]
				return "h" + strconv.Itoa(len(level)) + "." + content
			},
		},
		{ // Headers (H1 and H2 underlines)
			re: regexp.MustCompile(`(?m)^(.*?)\n([=-]+)$`),
			repl: func(groups []string) string {
				_, content, level := groups[0], groups[1], groups[2]
				if level[0] == '=' {
					return "h1. " + content
				}
				return "h2. " + content
			},
		},
		{ // Ordered lists
			re: regexp.MustCompile(`(?m)^([ \t]*)\d+\.\s+`),
			repl: func(groups []string) string {
				_, spaces := groups[0], groups[1]
				return strings.Repeat("#", len(spaces)/2+1) + " "
			},
		},
		{ // Un-Ordered Lists
			re: regexp.MustCompile(`(?m)^([ \t]*)\*\s+`),
			repl: func(groups []string) string {
				_, spaces := groups[0], groups[1]
				return strings.Repeat("*", len(spaces)/2+1) + " "
			},
		},
		{ // Strikethrough
			re:   regexp.MustCompile(`<del>(.*?)</del>`),
			repl: "-$1-",
		},
		{ // Inserts
			re:   regexp.MustCompile(`<ins>(.*?)</ins>`),
			repl: "+$1+",
		},
		{ // Superscript
			re:   regexp.MustCompile(`<sup>(.*?)</sup>`),
			repl: "^$1^",
		},
		{ // Subscript
			re:   regexp.MustCompile(`<sub>(.*?)</sub>`),
			repl: "~$1~",
		},
		{ // Other kind of strikethrough
			re:   regexp.MustCompile(`(\s+)~~(.*?)~~(\s+)`),
			repl: "$1-$2-$3",
		},
		{ // Images
			re:   regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`),
			repl: "!$1!",
		},
		{ // Named Link
			re:   regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`),
			repl: "[$1|$2]",
		},
		{ // Un-Named Link
			re:   regexp.MustCompile(`<([^>]+)>`),
			repl: "[$1]",
		},
		{ // Single Paragraph Blockquote
			re:   regexp.MustCompile(`(?m)^>`),
			repl: "bq.",
		},
	}

	// Code is converted up front and swapped for placeholders, so that the
	// jirations can't alter its contents
	codes := []string{}
	codeJirations := []jiration{
		{ // Named/Un-Named Code Block
			re: regexp.MustCompile("```(.+\n)?((?:.|\n)*?)```"),
			repl: func(groups []string) string {
				_, synt, content := groups[0], groups[1], groups[2]
				code := "{code}"
				if synt != "" {
					code = "{code:" + strings.ReplaceAll(synt, "\n", "") + "}\n"
				}
				codes = append(codes, code+content+"{code}")
				return codePlaceholder(len(codes) - 1)
			},
		},
		{ // Inline-Preformatted Text
			re: regexp.MustCompile("`([^`]+)`"),
			repl: func(groups []string) string {
				_, content := groups[0], groups[1]
				codes = append(codes, "{{"+content+"}}")
				return codePlaceholder(len(codes) - 1)
			},
		},
	}
	for _, jiration := range append(codeJirations, jirations...) {
		switch v := jiration.repl.(type) {
		case string:
			str = jiration.re.ReplaceAllString(str, v)
		case func([]string) string:
			str = replaceAllStringSubmatchFunc(jiration.re, str, v)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	for i, code := range codes {
		str = strings.Replace(str, codePlaceholder(i), code, 1)
	}
	// Tables and panels lose the blank line before them, which needs keeping
	// at the start of the document
	if len(str)-len(strings.TrimLeft(str, "\n")) < leadingNewlines {
		str = "\n" + str
	}
	return str
}

// codePlaceholder uses NUL bytes, which no jiration matches, to mark where code was
func codePlaceholder(i int) string {
	return "\x00" + strconv.Itoa(i) + "\x00"
}

// https://gist.github.com/elliotchance/d419395aa776d632d897
func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] < 0 { // optional group that didn't participate in the match
				groups = append(groups, "")
				continue
			}
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
