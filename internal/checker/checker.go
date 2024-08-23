package checker

import (
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
	ExpiresIn int64
	Status    CheckResultStatus
}
