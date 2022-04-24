package utils

// Replaces a character in a string at the given index, returning the updated string.
func ReplaceCharAtStringIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}
