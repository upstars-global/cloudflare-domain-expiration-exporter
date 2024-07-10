package cf

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
)

type cfImpl struct {
	api              *cloudflare.API
	obfuscatedAPIKey string
}

type CF interface {
	GetDomains(ctx context.Context) ([]string, error)
	GetAPIKeyObfuscated() string
}

func New(apiKey string) (CF, error) {
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

func (c *cfImpl) GetDomains(ctx context.Context) ([]string, error) {
	zones, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	var domains []string
	for _, z := range zones {
		domains = append(domains, z.Name)
	}

	return domains, nil
}

func (c *cfImpl) GetAPIKeyObfuscated() string {
	return c.obfuscatedAPIKey
}
