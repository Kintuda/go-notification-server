package exception

import (
	"strings"
	"unicode"
)

func LowerFirstChar(s string) string {
	if len(s) <= 2 {
		return strings.ToLower(s)
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func FormatMessage(tag string) string {
	switch tag {
	case "uuid":
		return "uuid is invalid"
	case "required":
		return "field is required, and should not null or empty"
	default:
		return ""
	}
}
