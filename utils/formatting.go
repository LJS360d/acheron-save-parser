package utils

import (
	"strings"
	"unicode"
)

// ToTitleCase converts the first character to uppercase and the rest to lowercase.
func ToTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// ToCapitalized capitalizes the first letter of each word separated by spaces, underscores, or hyphens.
func ToCapitalized(s string) string {
	s = strings.ToLower(s)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// ToSnakeCase converts a string with spaces, camelCase, or kebab-case to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	s = strings.ReplaceAll(s, "-", " ") // Replace hyphens with spaces
	s = strings.ReplaceAll(s, "_", " ") // Replace underscores with spaces

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && s[i-1] != ' ' {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == ' ' {
			if i > 0 && s[i-1] != ' ' {
				result.WriteRune('_')
			}
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ToKebabCase converts the string to kebab-case by replacing underscores with hyphens and adding hyphens between transitions from lowercase to uppercase letters.
func ToKebabCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == '_' {
			result.WriteRune('-')
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
