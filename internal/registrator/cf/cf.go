package cf

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator"
	"github.com/upstars-global/domains-expiration-exporter/internal/types"
	"time"
)

type cfImpl struct {
	api              *cloudflare.API
	obfuscatedAPIKey string
}

func New(apiKey string) (registrator.Registrator, error) {
	api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		return nil, err
	}

	impl := &cfImpl{
		api:              api,
		obfuscatedAPIKey: apiKey[0:6] + "..." + apiKey[len(apiKey)-6:],
	}

	// Check permissions
	err = impl.checkPermissions(context.Background())
	if err != nil {
		return nil, err
	}

	return impl, nil
}

func (c *cfImpl) checkPermissions(ctx context.Context) error {
	_, err := c.GetDomains(ctx)
	if err != nil {
		return fmt.Errorf("invalid access token=%s: %w", c.GetAPIKeyObfuscated(), err)
	}
	return nil
}

func (c *cfImpl) GetDomains(ctx context.Context) (types.DomainInfo, error) {
	var info types.DomainInfo

	zones, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	for _, z := range zones {
		t, err := c.getExpirationDate(ctx, z.Name)
		if err != nil {
			return nil, err
		}
		info[z.Name] = t
	}

	return info, nil
}

func (c *cfImpl) GetAPIKeyObfuscated() string {
	return c.obfuscatedAPIKey
}

func (c *cfImpl) getExpirationDate(ctx context.Context, domain string) (time.Time, error) {
	//TODO: write real function that return expiration date for certain domain
	return time.Time{}, nil
}
