package utils

import (
	"fmt"
)

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
