package api

import (
	"fmt"
	"io"
	"net/http"

	"git.nymeria.io/nymeria.go"
)

var (
	Client = http.Client{
		Timeout: nymeria.Timeout,
	}
)

func Request(method, endpoint string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", nymeria.BaseURL, endpoint), data)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"X-Api-Key":    []string{nymeria.ApiKey},
		"Content-Type": []string{"application/json"},
		"User-Agent":   []string{nymeria.UserAgent},
	}

	return req, nil
}
