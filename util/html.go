package util

import (
	"errors"
	"strings"
)

func AddHead(html string, header string) (string, error) {
	tag := `<head>`
	i := strings.Index(html, tag)
	if i == -1 {
		return "", errors.New("<head> element doesn't exist in html")
	}
	i += len(tag)
	return html[:i] + header + html[i:], nil
}
