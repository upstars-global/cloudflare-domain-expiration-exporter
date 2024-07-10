package main

import (
	"github.com/ghodss/yaml"
	"os"
	"strings"
	"time"
)

// argvGetApiKeys is a helper function to get the API keys from the environment variable
func argvGetApiKeys() []string {
	keys := make([]string, 0)

	parts := strings.Split(os.Getenv("CF_API_KEYS"), ",")

	for _, part := range parts {
		if part != "" {
			keys = append(keys, part)
		}
	}

	return keys
}

// parseManualExpirations is a helper function to parse the manual expirations
// reads file called `manual_expiration.yaml` (if exists) and return a map of domain -> expiration
func parseManualExpirations() (map[string]time.Time, error) {
	fn := "manual_expiration.yaml"
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return make(map[string]time.Time), nil
	}

	y, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var m struct {
		Domains map[string]time.Time `json:"domains"`
	}
	if err = yaml.Unmarshal(y, &m); err != nil {
		return nil, err
	}

	return m.Domains, nil
}
