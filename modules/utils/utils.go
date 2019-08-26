package utils

import (
	"fmt"
	"time"
)

// ShortenCommit is used for shortening a commit to 8 characters.
func ShortenCommit(str string) string {
	if len(str) > 8 {
		return str[:8]
	}
	return str
}

// FormatDate formats Unix time into a readable format.
func FormatDate(i int64) string {
	return time.Unix(i, 0).Format("Jan 2, 2006 at 3:04 PM")
}

// Spacify adds a space before every capital letter.
func Spacify(str string) string {
	var final string
	hasSpace := false
	for i, ch := range str {
		if i != 0 && !hasSpace && ch >= 'A' && ch <= 'Z' {
			final += fmt.Sprintf(" %c", ch)
		} else {
			final += string(ch)
		}
		hasSpace = ch == ' '
	}
	return final
}
