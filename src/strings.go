package main

// contains checks if a string is present in a list of strings.
func contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
