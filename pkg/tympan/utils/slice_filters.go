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

// Looks in a slice of strings to see if any of them match; if so, returns the index of the match in the slice.
// If no members of the slice match, returns -1.
func FindIndex(source []string, match string) (index int) {
	index = -1

	for i, v := range source {
		if v == match {
			index = i
		}
	}

	return index
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
