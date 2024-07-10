package checker

import (
	"context"
	"github.com/upstars-global/cloudflare-domain-expiration-exporter/internal/cf"
	"github.com/upstars-global/cloudflare-domain-expiration-exporter/internal/expiration"
	"go.uber.org/zap"
	"sync"
	"time"
)

// checkInterval is the interval to check the domains
const checkInterval = 1 * time.Hour

// retriesCount is the number of retries to get the expiration
const retriesCount = 3

// backoffMultiplier is used to increase the backoff time between retries
const backoffMultiplier = 2

// expiresInUnknown is the value to indicate that the expiration is unknown
const expiresInUnknown = -999999

type Checker interface {
	Start()
	GetExpirations() map[string]CheckResult
}

// CheckResultStatus is the status of the check result (ok: ttl detected, unknown: ttl not detected)
type CheckResultStatus string

const (
	CheckResultStatusOK      CheckResultStatus = "ok"
	CheckResultStatusUnknown CheckResultStatus = "unknown"
)

// CheckResult is the result of the check
type CheckResult struct {
	ExpiresIn int
	Status    CheckResultStatus
}

type checkerImpl struct {
	log        *zap.Logger
	apis       []*apiInfo
	expiration expiration.Expiration
	mux        sync.Mutex
}

type apiInfo struct {
	api          cf.CF
	checkResults map[string]CheckResult
}

// New creates a new checker
func New(log *zap.Logger, apis []cf.CF, exp expiration.Expiration) Checker {
	infos := make([]*apiInfo, 0)
	for _, api := range apis {
		infos = append(infos, &apiInfo{
			api: api,
		})
	}

	return &checkerImpl{
		log:        log.With(zap.String("component", "checker")),
		apis:       infos,
		expiration: exp,
	}
}

// GetExpirations returns the domains expiration
func (c *checkerImpl) GetExpirations() map[string]CheckResult {
	c.mux.Lock()
	defer c.mux.Unlock()

	res := make(map[string]CheckResult)
	for _, api := range c.apis {
		for d, e := range api.checkResults {
			res[d] = e
		}
	}

	return res
}

// Start starts the checker
func (c *checkerImpl) Start() {
	c.check()

	t := time.NewTicker(checkInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			c.check()
		}
	}
}

// deleteUnused deletes the unused domains from the check results
func (c *checkerImpl) deleteUnused(api *apiInfo, domains []string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for d := range api.checkResults {
		if !contains(domains, d) {
			delete(api.checkResults, d)
		}
	}
}

// checkAPI checks the expiration for the given API
func (c *checkerImpl) checkAPI(ctx context.Context, api *apiInfo) {
	c.log.Info("checking api", zap.String("api", api.api.GetAPIKeyObfuscated()))

	// initialize check results
	c.mux.Lock()
	if api.checkResults == nil {
		api.checkResults = make(map[string]CheckResult)
	}
	c.mux.Unlock()

	domains, err := api.api.GetDomains(ctx)
	if err != nil {
		c.log.Error("failed to list domains", zap.Error(err), zap.String("api", api.api.GetAPIKeyObfuscated()))
		return
	}

	// delete unused domains (removed from CloudFlare zones)
	// so we don't need to check them
	c.deleteUnused(api, domains)

	// check expiration for each domain
	for _, d := range domains {
		// closure to check expiration
		check := func() (expiresIn int, e error) {
			expiresIn, _, e = c.expiration.DaysTillExpiration(d)
			return
		}

		// get expiration with retries
		expiresIn := 0
		backoff := time.Duration(1)
		for i := 0; i < retriesCount; i++ {
			var e error
			expiresIn, e = check()
			if err != nil {
				c.log.Error("failed to get expiration", zap.Error(e), zap.String("domain", d), zap.String("api", api.api.GetAPIKeyObfuscated()), zap.Int("attempt", i+1))
				time.Sleep(backoff * time.Second)
				backoff *= backoffMultiplier
				continue
			}
			break
		}

		c.log.Info("zone result", zap.String("domain", d), zap.Int("expires_in", expiresIn), zap.String("api", api.api.GetAPIKeyObfuscated()))

		checkResult := CheckResult{
			ExpiresIn: expiresIn,
			Status:    CheckResultStatusOK,
		}

		if expiresIn == 0 {
			checkResult.Status = CheckResultStatusUnknown
		}

		c.mux.Lock()
		api.checkResults[d] = checkResult
		c.mux.Unlock()
	}
}

// check checks the domains expiration
func (c *checkerImpl) check() {
	c.log.Info("checking domains")

	ctx := context.Background()

	for _, api := range c.apis {
		go c.checkAPI(ctx, api)
	}
}

func contains(domains []string, domain string) bool {
	for _, d := range domains {
		if d == domain {
			return true
		}
	}
	return false
}
