package nymeria

import (
	"fmt"
	"net/http"
)

var (
	client = http.Client{}
)

func request(method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, endpoint), nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"X-Api-Key":    []string{apiKey},
		"Content-Type": []string{"application/json"},
	}

	return req, nil
}
