package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSearchDomainsList(t *testing.T) {
	_ = os.Setenv("DOMAINS_LIST", "google.com")

	val, ok, _ := searchDomainsList()
	assert.Equal(t, true, ok)
	assert.Equal(t, "google.com", val.List[0], "there should be 'google.com' in the list")

	_ = os.Unsetenv("DOMAINS_LIST")
	_, ok, _ = searchDomainsList()
	assert.Equal(t, false, ok, "should not find DOMAINS_LIST variable")

	_ = os.Setenv("DOMAINS_LIST", "")
	_, _, err := searchDomainsList()
	assert.NotNil(t, err, "there should be an error 'DOMAINS_LIST env var found but is empty'")
}
