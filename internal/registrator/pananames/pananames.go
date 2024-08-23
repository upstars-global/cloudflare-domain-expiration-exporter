package pananames

import (
	"context"
	pananames "github.com/pananames/go-api-client"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator"
	"github.com/upstars-global/domains-expiration-exporter/internal/types"
)

type pnImpl struct {
	api              *pananames.Client
	obfuscatedAPIKey string
}

func New(key string) (registrator.Registrator, error) {
	client, err := pananames.NewClient(key)
	if err != nil {
		return nil, err
	}

	return &pnImpl{
		api:              client,
		obfuscatedAPIKey: "key[0:6] + \"...\" + key[len(key)-6:]",
	}, nil
}

func (p *pnImpl) GetDomains(ctx context.Context) (types.DomainInfo, error) {
	domainsList := make(types.DomainInfo)

	listOptions := &pananames.GetDomainsOptions{ListOptions: pananames.ListOptions{Limit: 30, Page: 1}}
	for {
		domains, page, err := p.api.GetDomains(listOptions)
		if err != nil {
			return nil, err
		}
		for _, d := range domains {
			domainsList[d.Domain] = d.ExpirationDate.Time
		}
		if listOptions.Page = page.NextPage(); listOptions.Page == 0 {
			break
		}
	}

	return domainsList, nil
}

func (p *pnImpl) GetAPIKeyObfuscated() string {
	return p.obfuscatedAPIKey
}
