package utils

// Checks if a string is present in a slice of strings
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Returns all matches for a string from a slice of strings
func Find(source []string, match string) (matches []string) {
	if Contains(source, match) {
		matches = append(matches, match)
	}
	return
}

// Returns the input slice without the entry at the specified index.
func RemoveIndex[T any](slice []T, index int) (returnSlice []T) {
	returnSlice = append(returnSlice, slice[:index]...)
	return append(returnSlice, slice[index+1:]...)
}
