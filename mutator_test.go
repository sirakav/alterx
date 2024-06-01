package alterx

import (
	"bytes"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testConfig = Config{
	Patterns: []string{
		"{{sub}}-{{word}}.{{root}}",     // ex: api-prod.scanme.sh
		"{{word}}-{{sub}}.{{root}}",     // ex: prod-api.scanme.sh
		"{{word}}.{{sub}}.{{root}}",     // ex: prod.api.scanme.sh
		"{{sub}}.{{word}}.{{root}}",     // ex: api.prod.scanme.sh
		"{{word}}-{{sub}}.{{suffix}}",   // ex: dev-api
		"{{sub}}-{{word}}.{{suffix}}",   // ex: api-dev
		"{{word}}.{{sub}}.{{suffix}}",   // ex: dev.api.scanme.sh
		"{{sub}}.{{word}}.{{suffix}}",   // ex: api.dev.scanme.sh
		"{{sub}}{{number}}.{{suffix}}",  // ex: www123.scanme.sh
		"{{word}}.{{suffix}}",           // ex: prod.scanme.sh
		"{{sub}}{{word}}.{{suffix}}",    // ex: devtest.scanme.sh
		"{{region}}.{{sub}}.{{suffix}}", // ex: us-west.www.scanme.sh
		"{{word}}{{number}}.{{suffix}}", // ex: prod123.scanme.sh
	},
	Payloads: map[string][]string{
		"word":   {"dev", "lib", "prod"},
		"sub":    {"www", "mail", "ftp"},
		"tld":    {"com", "net", "org"},
		"number": {"123", "456", "789"},
		"suffix": {"scanme", "sh"},
		"region": {"us-west", "us-east", "eu-west"},
	},
}

func TestMutatorCount(t *testing.T) {
	opts := &Options{
		Domains: []string{"api.scanme.sh", "chaos.scanme.sh", "nuclei.scanme.sh", "cloud.nuclei.scanme.sh"},
	}
	opts.Patterns = testConfig.Patterns
	opts.Payloads = testConfig.Payloads

	m, err := New(opts)
	require.Nil(t, err)
	require.EqualValues(t, 180, m.EstimateCount())
}

func TestMutatorResults(t *testing.T) {
	opts := &Options{
		Domains: []string{"api.scanme.sh", "chaos.scanme.sh", "nuclei.scanme.sh", "cloud.nuclei.scanme.sh"},
	}
	opts.Patterns = testConfig.Patterns
	opts.Payloads = testConfig.Payloads
	opts.MaxSize = math.MaxInt
	m, err := New(opts)
	require.Nil(t, err)
	var buff bytes.Buffer
	err = m.ExecuteWithWriter(&buff)
	require.Nil(t, err)
	count := strings.Split(strings.TrimSpace(buff.String()), "\n")
	require.EqualValues(t, 180, len(count), buff.String())
}

func BenchmarkNew(b *testing.B) {
	opts := &Options{
		Domains: []string{"api.scanme.sh", "chaos.scanme.sh", "nuclei.scanme.sh", "cloud.nuclei.scanme.sh"},
	}
	opts.Patterns = testConfig.Patterns
	opts.Payloads = testConfig.Payloads

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(opts)
	}
}

func BenchmarkExecuteWithWriter(b *testing.B) {
	categories := []string{"word", "number", "sub", "tld", "region"}
	numEntries := 10000
	payloads := generateRandomPayloads(categories, numEntries)

	tests := []struct {
		name          string
		domains       []string
		dedupeResults bool
	}{
		{
			name:          "1k with dedupe",
			domains:       generateNRandomDomains(1000),
			dedupeResults: true,
		},
		{
			name:          "1k no dedupe",
			domains:       generateNRandomDomains(1000),
			dedupeResults: false,
		},
		{
			name:          "10k with dedupe",
			domains:       generateNRandomDomains(10000),
			dedupeResults: true,
		},
		{
			name:          "10k no dedupe",
			domains:       generateNRandomDomains(10000),
			dedupeResults: false,
		},
		{
			name:          "100k with dedupe",
			domains:       generateNRandomDomains(100000),
			dedupeResults: true,
		},
		{
			name:          "100k no dedupe",
			domains:       generateNRandomDomains(100000),
			dedupeResults: false,
		},
		{
			name:          "1M with dedupe",
			domains:       generateNRandomDomains(1000000),
			dedupeResults: true,
		},
		{
			name:          "1M no dedupe",
			domains:       generateNRandomDomains(1000000),
			dedupeResults: false,
		},
	}

	for _, tt := range tests {
		opts := &Options{
			Domains:       tt.domains,
			Patterns:      testConfig.Patterns,
			Payloads:      payloads,
			DedupeResults: tt.dedupeResults,
			MaxSize:       math.MaxInt,
		}

		m, err := New(opts)
		require.Nil(b, err)

		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buff bytes.Buffer
				err = m.ExecuteWithWriter(&buff)
				// _ = m.Execute(context.Background())
			}
		})
	}
}

func generateNRandomDomains(n int) []string {
	var domains []string
	for i := 0; i < n; i++ {
		domains = append(domains, randomString(5)+".scanme.sh")
	}
	return domains
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}

// generateRandomPayloads generates random payloads for the specified categories.
func generateRandomPayloads(categories []string, numEntries int) map[string][]string {
	rand.Seed(time.Now().UnixNano())
	payloads := make(map[string][]string)
	for _, category := range categories {
		for i := 0; i < numEntries; i++ {
			payloads[category] = append(payloads[category], randomString(8)) // random string of length 8
		}
	}
	return payloads
}

func BenchmarkEstimateCount(b *testing.B) {
	categories := []string{"word", "number", "sub", "tld"}
	numEntries := 100000
	payloads := generateRandomPayloads(categories, numEntries)

	tests := []struct {
		name    string
		domains []string
	}{
		{
			name:    "1k domains",
			domains: generateNRandomDomains(1000),
		},
		{
			name:    "10k domains",
			domains: generateNRandomDomains(10000),
		},
		{
			name:    "100k domains",
			domains: generateNRandomDomains(100000),
		},
		{
			name:    "1M domains",
			domains: generateNRandomDomains(1000000),
		},
	}

	for _, tt := range tests {
		opts := &Options{
			Domains: tt.domains,
		}
		opts.Patterns = testConfig.Patterns
		opts.Payloads = payloads

		b.Run(tt.name, func(b *testing.B) {
			m, err := New(opts)
			require.Nil(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.EstimateCount()
			}
		})
	}
}

func TestNewWithPatternDetection(t *testing.T) {
	opts := &Options{
		Domains:          []string{"api-dev.scanme.sh", "api-1.scanme.sh", "prod123.scanme.sh"},
		PatternDetection: true,
		Payloads: map[string][]string{
			"word":    {"api", "dev", "test", "prod"},
			"number":  {"1", "2", "3", "1234", "123"},
			"country": {"us", "uk", "ca"},
		},
		Patterns: testConfig.Patterns,
		Enrich:   true,
	}

	m, err := New(opts)
	require.Nil(t, err)
	require.EqualValues(t, 393, m.EstimateCount())
}
