package util

import (
	"regexp"
	"strings"
)

// ParseURL takes a possible URL and splits it into host and path or returns empty strings if the parameter is not a URL
func ParseURL(possibleURL string) (string, string) {
	r := regexp.MustCompile(`(https?:\/\/.*?)\/([\w\.]+)(.*)`)
	match := r.FindStringSubmatch(possibleURL)

	if len(match) < 3 {
		// not a valid hastebin URL
		return "", ""
	}

	if match[2] == "" {
		// empty path
		return "", ""
	}

	if strings.Contains(match[3], "/") {
		// no a key
		return "", ""
	}

	return match[1], match[2]
}
