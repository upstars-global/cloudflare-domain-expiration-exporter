package types

import (
	"github.com/upstars-global/domains-expiration-exporter/internal/checker"
	"time"
)

type DomainInfo map[string]time.Time

func (d *DomainInfo) GetExpirations() map[string]checker.CheckResult {
	now := time.Now()
	results := make(map[string]checker.CheckResult)

	for domain, expDate := range *d {

		result := checker.CheckResult{
			ExpiresIn: int64(expDate.Sub(now).Seconds()),
			Status:    checker.CheckResultStatusOK,
		}

		results[domain] = result
	}
	return results
}
