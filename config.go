package nymeria

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

var (
	Client = http.Client{
		Timeout: Timeout,
	}
)

func Request(method, endpoint string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", BaseURL, endpoint), data)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"X-Api-Key":  []string{ApiKey},
		"User-Agent": []string{UserAgent},
	}

	return req, nil
}

func Normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
