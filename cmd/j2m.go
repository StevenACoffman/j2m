///usr/bin/env go run "$0" "$@" ; exit "$?"
package main

import (
	"fmt"
	"github.com/StevenACoffman/j2m"
	"io/ioutil"
	"os"
)

// Expects piped input of Jira Markdown, outputs Github Markdown
func main() {

	str := ""
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		fmt.Println("The command is intended to work with pipes but didn't get one. Assuming empty input")
	} else {
		stdInBytes, _ := ioutil.ReadAll(os.Stdin)
		str = string(stdInBytes)
	}

	str = j2m.JiraToMD(str)

	fmt.Printf("%s\n", str)
}
