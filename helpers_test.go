package j2m

import (
	"regexp"
	"testing"
)

func TestReplaceAllStringSubmatchFuncUnmatchedOptionalGroup(t *testing.T) {
	re := regexp.MustCompile(`a(b)?(c)`)
	got := replaceAllStringSubmatchFunc(re, "ac abc", func(groups []string) string {
		return "[" + groups[1] + "|" + groups[2] + "]"
	})
	if want := "[|c] [b|c]"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
