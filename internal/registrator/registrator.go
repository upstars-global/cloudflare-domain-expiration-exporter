package registrator

import (
	"context"
	"github.com/upstars-global/domains-expiration-exporter/internal/types"
)

type Registrator interface {
	GetDomains(ctx context.Context) (types.DomainInfo, error)
	GetAPIKeyObfuscated() string
}
