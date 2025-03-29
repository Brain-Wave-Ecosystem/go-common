package helpers

import (
	"strings"
	"unicode"
)

func ToTitleCase(name string) string {
	runes := []rune(name)

	runes[0] = unicode.ToUpper(runes[0])

	result := string(runes[0]) + strings.ToLower(string(runes[1:]))

	return strings.TrimSpace(result)
}

func ToTitleCaseArray(name string) string {
	parts := strings.Split(name, " ")

	for i := range parts {
		parts[i] = ToTitleCase(parts[i])
	}

	return strings.Join(parts, " ")
}
