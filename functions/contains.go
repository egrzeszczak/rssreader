package rssfunctions

// contains function checks if a specific word is present in a slice of strings.
func Contains(slice []string, word string) bool {
	for _, s := range slice {
		if s == word {
			return true
		}
	}
	return false
}
