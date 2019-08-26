[![GoDoc](https://godoc.org/github.com/StevenACoffman/j2m?status.svg)](https://godoc.org/github.com/StevenACoffman/j2m)
[![GoReportcard](https://goreportcard.com/badge/github.com/StevenACoffman/j2m?status.svg)](https://goreportcard.com/report/github.com/StevenACoffman/j2m)

# jira-to-md

## JIRA to MarkDown text format converter
Golang tool to convert from JIRA Markdown text formatting to GitHub Flavored MarkDown.

## Credits
This fun toy was heavily inspired by the J2M project by Fokke Zandbergen (http://j2m.fokkezb.nl/). Major credit to Fokke, kylefarris (and other contributors) for establishing the RexExp patterns for this to work. The maintained JavaScript fork I based this on is [here](https://github.com/kylefarris/J2M)

## Supported Conversions

* Headers (H1-H6)
* Bold
* Italic
* Bold + Italic
* Un-ordered lists
* Ordered lists
* Programming Language-specific code blocks (with help from herbert-venancio)
* Inline preformatted text spans
* Un-named links
* Named links
* Monospaced Text
* ~~Citations~~ (currently kinda buggy)
* Strikethroughs
* Inserts
* Superscripts
* Subscripts
* Single-paragraph blockquotes
* Tables
* Panels


## How to Use

### Markdown String

```
**Some bold things**
*Some italic stuff*
## H2
<http://google.com>
```

### Atlassian Wiki MarkUp Syntax (JIRA)

We'll refer to this as the `jira` variable in the examples below.

```
*Some bold things**
_Some italic stuff_
h2. H2
[http://google.com]
```

### Examples

```
cat j2m.jira | j2m
```
