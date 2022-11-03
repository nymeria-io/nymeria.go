package nymeria

import (
	"fmt"
	"net/http"
	"time"
)

const (
	ApiVersion = "4"
	BaseURL    = "https://www.nymeria.io/api/v" + ApiVersion
	UserAgent  = "nymeria.go/" + ApiVersion
)

var (
	// Timeout is the default timeout used by the HTTP client (default: 5 seconds).
	Timeout = 30 * time.Second

	// The API key that will be used for all authenticated requests.
	ApiKey string

	ErrInvalidParameters      = fmt.Errorf(`error: invalid parameter(s)`)
	ErrBadRequest             = fmt.Errorf(`error: bad request; perhaps your parameters were wrong`)
	ErrAuthenticationRequired = fmt.Errorf(`error: invalid or unauthorized api key detected`)
	ErrPaymentRequired        = fmt.Errorf(`error: payment required; perhaps your plan expired or was exhausted`)
	ErrNotFound               = fmt.Errorf(`error: resource not found`)
	ErrServerError            = fmt.Errorf(`error: server error detected`)

	ErrMap = map[int]error{
		http.StatusBadRequest:      ErrBadRequest,
		http.StatusNotFound:        ErrNotFound,
		http.StatusPaymentRequired: ErrPaymentRequired,
	}
)
