package nymeria

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

var (
	// ErrInvalidRequest is returned any time a request looks malformed or inavlid.
	// This can be due to a bad parameter, or bad encoding.
	ErrInvalidRequest = errors.New(`error: the request looks invalid or malformed`)
)

// Verification is the result of an email verify request. The "result" can either be
// valid or invalid. The tags can include things like:
//
//		has_dns_mx, smtp_connectable, accepts_all
//
type Verification struct {
	Status string `json:"status"`

	Meta struct {
		Email string `json:"email"`
	} `json:"meta"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data struct {
		Result string   `json:"result"`
		Tags   []string `json:"tags"`
	} `json:"data"`
}

// Verify takes a professional email address and tries to verify its validity.
func Verify(email string) (*Verification, error) {
	req, err := request("GET", fmt.Sprintf("/verify?email=%s", url.QueryEscape(email)))

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response Verification

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}
