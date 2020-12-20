package util

import "regexp"

// ParseURL takes a possible URL and splits it into host and path or returns empty strings if the parameter is not a URL
func ParseURL(possibleURL string) (string, string) {
	r := regexp.MustCompile(`(.*?//.*?)/(.*?)[\?$]`)
	match := r.FindStringSubmatch(possibleURL)

	if len(match) < 3 {
		// no match
		return "", ""
	}

	return match[1], match[2]
}
