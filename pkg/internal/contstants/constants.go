//pkg/internal/constants/constants.go
package constants

import "time"

const (
    BaseURL        = "https://api.anthropic.com"
    DefaultTimeout = 10 * time.Second
    MaxRetries     = 3
    MinRetryDelay  = 1 * time.Second
    MaxRetryDelay  = 30 * time.Second
)
