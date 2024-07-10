package expiration

import (
	"fmt"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"time"
)

type expirationImpl struct {
	manualExpirations map[string]time.Time
}

type Expiration interface {
	DaysTillExpiration(domain string) (daysTillExpiration int, expiringAt time.Time, err error)
}

func New(manualExpirations map[string]time.Time) Expiration {
	return &expirationImpl{
		manualExpirations: manualExpirations,
	}
}

func (e *expirationImpl) DaysTillExpiration(domain string) (daysTillExpiration int, expiringAt time.Time, err error) {
	wh, _ := whois.Whois(domain)
	parsed, err := whoisparser.Parse(wh)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to parse whois information for domain %s: %w", domain, err)
	}

	if parsed.Domain.ExpirationDateInTime == nil {
		if manualExpirations, ok := e.manualExpirations[domain]; ok {
			daysTillExpiration = int(manualExpirations.Sub(time.Now()).Hours() / 24)
			expiringAt = manualExpirations
			return
		}

		return 0, time.Time{}, fmt.Errorf("failed to get expiration date for domain %s", domain)
	}

	daysTillExpiration = int(parsed.Domain.ExpirationDateInTime.Sub(time.Now()).Hours() / 24)
	expiringAt = *parsed.Domain.ExpirationDateInTime

	return
}
