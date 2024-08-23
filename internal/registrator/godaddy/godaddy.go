package godaddy

import (
	"context"
	"github.com/oze4/godaddygo"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator"
	"github.com/upstars-global/domains-expiration-exporter/internal/types"
	"net/http"
	"time"
)

const (
	MaxHTTPClientTimeout = time.Second * 60
)

type gdImpl struct {
	api              godaddygo.V1
	obfuscatedAPIKey string
}

func New(key, secret, env string) (registrator.Registrator, error) {
	var gdEnv godaddygo.APIEnv

	if env == "dev" {
		gdEnv = godaddygo.APIDevEnv
	} else {
		gdEnv = godaddygo.APIProdEnv
	}
	gdConfig := godaddygo.NewConfig(key, secret, gdEnv)
	gdClient := &http.Client{Timeout: MaxHTTPClientTimeout}

	api, err := godaddygo.WithClient(gdClient, gdConfig)
	if err != nil {
		return nil, err
	}

	return &gdImpl{
		api:              api.V1(),
		obfuscatedAPIKey: secret[0:6] + "..." + secret[len(secret)-6:],
	}, nil
}

func (g *gdImpl) GetDomains(ctx context.Context) (types.DomainInfo, error) {
	domainsInfo, err := g.api.ListDomains(ctx)
	if err != nil {
		return nil, err
	}

	domains := make(types.DomainInfo)
	for _, d := range domainsInfo {
		domains[d.Domain] = d.Expires
	}
	return domains, nil
}

func (g *gdImpl) GetAPIKeyObfuscated() string {
	return g.obfuscatedAPIKey
}
