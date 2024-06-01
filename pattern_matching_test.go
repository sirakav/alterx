package alterx

import "testing"

// import (
// 	"testing"
// )

// func TestCategorizeToken(t *testing.T) {
// 	payloads := map[string][]string{
// 		"word":    {"api", "dev", "test"},
// 		"number":  {"1", "2", "3"},
// 		"country": {"us", "uk", "ca"},
// 		"suffix":  {"com", "net", "org"},
// 	}

// 	tests := []struct {
// 		token    string
// 		expected string
// 	}{
// 		{"api", "{{word}}"},
// 		{"uk", "{{country}}"},
// 		{"3", "{{number}}"},
// 		{"xyz", "{{sub}}"},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.token, func(t *testing.T) {
// 			result := categorizeToken(test.token, payloads)
// 			if result != test.expected {
// 				t.Errorf("Expected %s, but got %s", test.expected, result)
// 			}
// 		})
// 	}
// }

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	tests := []struct {
		item     string
		expected bool
	}{
		{"a", true},
		{"d", false},
	}

	for _, test := range tests {
		t.Run(test.item, func(t *testing.T) {
			result := contains(slice, test.item)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestTokenizeSubDomain(t *testing.T) {
	payloads := map[string][]string{
		"word":    {"api", "dev", "test"},
		"number":  {"1", "2", "3", "1234"},
		"country": {"us", "uk", "ca"},
	}

	tests := []struct {
		domain   string
		expected []string
	}{
		{"uk1234", []string{"uk", "1234"}},
		{"api-dev", []string{"api", "-", "dev"}},
		{"ca1", []string{"ca", "1"}},
		{"xyz", []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			result := tokenizeSubDomain(test.domain, payloads)
			if len(result) != len(test.expected) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
				return
			}
			for i, token := range result {
				if token != test.expected[i] {
					t.Errorf("Expected %v, but got %v", test.expected, result)
					break
				}
			}
		})
	}
}

func TestDetectPatterns(t *testing.T) {
	payloads := map[string][]string{
		"word":    {"api", "dev", "test", "example"},
		"number":  {"1", "2", "3", "1234"},
		"country": {"us", "uk", "ca"},
	}

	tests := []struct {
		domain   Input
		expected string
	}{
		{Input{TLD: ".com", ETLD: "example.com", Root: "example", Sub: "uk1234", Suffix: "com", MultiLevel: []string{}}, "{{country}}{{number}}.{{suffix}}"},
		{Input{TLD: ".com", ETLD: "example.com", Root: "example", Sub: "api-dev", Suffix: "com", MultiLevel: []string{}}, "{{word}}-{{word}}.{{suffix}}"},
		{Input{TLD: ".net", ETLD: "example.net", Root: "example", Sub: "ca1", Suffix: "net", MultiLevel: []string{}}, "{{country}}{{number}}.{{suffix}}"},
		{Input{TLD: ".org", ETLD: "example.org", Root: "example", Sub: "dev123", Suffix: "org", MultiLevel: []string{}}, "{{word}}{{number}}.{{suffix}}"},
		{Input{TLD: ".org", ETLD: "example.org", Root: "example", Sub: "api", Suffix: "org", MultiLevel: []string{"dev"}}, "{{word}}.{{word}}.{{suffix}}"},
		{Input{TLD: ".org", ETLD: "example.org", Root: "example", Sub: "api", Suffix: "org", MultiLevel: []string{"api", "api"}}, "{{word}}.{{word}}.{{word}}.{{suffix}}"},
		{Input{TLD: ".com", ETLD: "example.com", Root: "example", Sub: "uk-1234", Suffix: "com", MultiLevel: []string{}}, "{{country}}-{{number}}.{{suffix}}"},
		{Input{TLD: ".com", ETLD: "example.com", Root: "example", Sub: "dev1234", Suffix: "com", MultiLevel: []string{}}, "{{word}}{{number}}.{{suffix}}"},
		{Input{TLD: ".com", ETLD: "example.com", Root: "example", Sub: "unknown", Suffix: "com", MultiLevel: []string{"unknown"}}, "{{suffix}}"},
	}

	for _, test := range tests {
		t.Run(test.domain.Sub, func(t *testing.T) {
			result := detectPatterns(test.domain, payloads)
			if len(result) != len(test.expected) {
				t.Errorf("Expected %v, but got %v, subdomain %v", test.expected, result, test.domain.Sub)
				return
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v, subdomain %v", test.expected, result, test.domain.Sub)
			}
		})
	}
}
