// package alterx

// import (
// 	"strconv"
// 	"strings"
// )

// const endToken = "{{suffix}}"

// // categorizeToken categorizes a given token based on payloads
// func categorizeToken(token string, payloads map[string][]string) string {
// 	for category, values := range payloads {
// 		if contains(values, token) {
// 			return "{{" + category + "}}"
// 		}
// 	}
// 	// If the token is not in the payloads, check if it's a number
// 	_, err := strconv.Atoi(token)
// 	if err == nil {
// 		return "{{number}}"
// 	}

// 	// Dash is a special case
// 	if token == "-" {
// 		return "-"
// 	}

// 	// If it's not a number, return it as a subdomain
// 	return endToken
// }

// // contains checks if a slice contains a string
// func contains(slice []string, item string) bool {
// 	for _, v := range slice {
// 		if v == item {
// 			return true
// 		}
// 	}
// 	return false
// }

// // tokenizeDomain tokenizes a domain string into meaningful parts based on the payloads
// func tokenizeDomain(domain string, payloads map[string][]string) []string {
// 	var tokens []string
// 	i := 0
// 	for i < len(domain) {
// 		// handle dashes as a special case
// 		if domain[i] == '-' {
// 			tokens = append(tokens, "-")
// 			i++
// 			continue
// 		}

// 		longestMatch := ""
// 		matchLen := 0
// 		// Check each substring starting from the current position
// 		for j := i + 1; j <= len(domain); j++ {
// 			substr := domain[i:j]
// 			category := categorizeToken(substr, payloads)
// 			if category != endToken {
// 				longestMatch = substr
// 				matchLen = j - i
// 			}
// 		}
// 		// If a match was found, add it to tokens and advance the index
// 		if longestMatch != "" {
// 			tokens = append(tokens, longestMatch)
// 			i += matchLen
// 		} else {
// 			// If no match, add the current character and move to the next
// 			tokens = append(tokens, string(domain[i]))
// 			i++
// 		}
// 	}
// 	return tokens
// }

// // detectPatterns detects patterns in a domain string based on payloads
// func detectPatterns(domain string, payloads map[string][]string) string {
// 	domainParts := strings.Split(domain, ".")
// 	var pattern strings.Builder
// 	for i, part := range domainParts {
// 		tokens := tokenizeDomain(part, payloads)
// 		for _, token := range tokens {
// 			tokenCategory := categorizeToken(token, payloads)

// 			pattern.WriteString(tokenCategory)
// 			if tokenCategory == endToken {
// 				if pattern.String() == endToken {
// 					return ""
// 				}
// 				return pattern.String()
// 			}
// 		}
// 		if i < len(domainParts)-1 {
// 			pattern.WriteString(".")
// 		}
// 	}

// 	return pattern.String()
// }

package alterx

import (
	"strconv"
	"strings"
)

// categorizeToken categorizes a given token based on payloads
func categorizeToken(token string, payloads map[string][]string) string {
	for category, values := range payloads {
		if contains(values, token) {
			return "{{" + category + "}}"
		}
	}
	// If the token is not in the payloads, check if it's a number
	_, err := strconv.Atoi(token)
	if err == nil {
		return "{{number}}"
	}

	// Dash is a special case
	if token == "-" {
		return "-"
	}

	// If it's not a number, return it as a subdomain
	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// tokenizeSubDomain tokenizes a domain string into meaningful parts based on the payloads
func tokenizeSubDomain(domain string, payloads map[string][]string) []string {
	var tokens []string
	i := 0
	for i < len(domain) {
		// handle dashes as a special case
		if domain[i] == '-' {
			tokens = append(tokens, "-")
			i++
			continue
		}

		longestMatch := ""
		matchLen := 0
		// Check each substring starting from the current position
		for j := i + 1; j <= len(domain); j++ {
			substr := domain[i:j]
			category := categorizeToken(substr, payloads)
			if category != "" {
				longestMatch = substr
				matchLen = j - i
			}
		}
		// If a match was found, add it to tokens and advance the index
		if longestMatch != "" {
			tokens = append(tokens, longestMatch)
			i += matchLen
		} else {
			// If no match, add the current character and move to the next
			tokens = append(tokens, string(domain[i]))
			i++
		}
	}
	return tokens
}

// detectPatterns detects patterns in a domain string based on payloads
func detectPatterns(input Input, payloads map[string][]string) string {
	var pattern strings.Builder
	tokens := tokenizeSubDomain(input.Sub, payloads)
	for _, token := range tokens {
		tokenCategory := categorizeToken(token, payloads)
		pattern.WriteString(tokenCategory)
	}

	// Tokenize multi-level subdomains
	for _, subdomain := range input.MultiLevel {
		pattern.WriteString(".")
		subTokens := tokenizeSubDomain(subdomain, payloads)
		for _, token := range subTokens {
			tokenCategory := categorizeToken(token, payloads)
			pattern.WriteString(tokenCategory)
		}
	}

	pattern.WriteString(".{{suffix}}")
	return pattern.String()
}
