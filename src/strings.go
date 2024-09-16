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

// convertMap converts a map[string]*string to a map[string]string
func convertMap(src map[string]*string) map[string]string {
	dst := make(map[string]string)

	for key, value := range src {
		if value != nil {
			dst[key] = *value
		} else {
			dst[key] = ""
		}
	}

	return dst
}
