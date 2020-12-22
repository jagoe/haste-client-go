package util

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		title  string
		url    string
		server string
		key    string
	}{
		{"No URL", "abcdef", "", ""},
		{"No HTTP URL", "ftp://hastebin/abcdef", "", ""},
		{"URL without path", "https://hastebin", "", ""},
		{"URL with empty path", "https://hastebin/", "", ""},
		{"URL with long path", "https://hastebin/path/abcdef", "", ""},
		{"HTTP URL with path", "http://hastebin/abcdef", "http://hastebin", "abcdef"},
		{"HTTPS URL with path", "https://hastebin/abcdef", "https://hastebin", "abcdef"},
		{"Valid URL with query", "https://hastebin/abcdef?q=s", "https://hastebin", "abcdef"},
		{"Valid URL with language-specific key", "https://hastebin/abcdef.yaml", "https://hastebin", "abcdef.yaml"},
	}

	for _, test := range tests {
		server, key := ParseURL(test.url)

		if server != test.server && key != test.key {
			t.Errorf(`%s: Parsing URL '%s' should have resulted in (server, key) = ("%s", "%s"), got ("%s", "%s")`,
				test.title, test.url, test.server, test.key, server, key)
		}
	}
}
