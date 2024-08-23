package domainsList

import (
	"context"
	"fmt"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator"
	"github.com/upstars-global/domains-expiration-exporter/internal/types"
)

type domainsImpl struct {
	domains          []string
	obfuscatedAPIKey string
}

func New(domains []string) (registrator.Registrator, error) {

	return &domainsImpl{
		domains:          domains,
		obfuscatedAPIKey: "manual",
	}, nil
}

func (r *domainsImpl) GetDomains(ctx context.Context) (types.DomainInfo, error) {
	domains := make(types.DomainInfo)
	for _, domain := range r.domains {
		wh, _ := whois.Whois(domain)
		parsed, err := whoisparser.Parse(wh)
		if err != nil {
			return nil, fmt.Errorf("failed to parse whois information for domain %s: %w", domain, err)
		}

		//if parsed.Domain.ExpirationDateInTime == nil {
		//  if manualExpirations, ok := r.domains[domain]; ok {
		//    daysTillExpiration = int(manualExpirations.Sub(time.Now()).Hours() / 24)
		//    expiringAt = manualExpirations
		//    return nil, err
		//  }
		//}

		domains[parsed.Domain.Domain] = *parsed.Domain.ExpirationDateInTime
	}
	return domains, nil
}

func (r *domainsImpl) GetAPIKeyObfuscated() string {
	return r.obfuscatedAPIKey
}
