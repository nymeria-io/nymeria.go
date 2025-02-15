package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nymeria-io/nymeria.go"
)

type BulkVerifyParams struct {
	Email    string      `json:"email"`
	MetaData interface{} `json:"metadata"`
}

func Verify(email string) (*Verification, error) {
	email = nymeria.Normalize(email)

	if len(email) == 0 {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := nymeria.Request("GET", fmt.Sprintf("/email/verify?email=%s", url.QueryEscape(email)), nil)

	if err != nil {
		return nil, err
	}

	resp, err := nymeria.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if e, ok := nymeria.ErrMap[resp.StatusCode]; ok {
			return nil, e
		}

		return nil, nymeria.ErrServerError
	}

	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response struct {
		Status int          `json:"status"`
		Data   Verification `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func BulkVerify(params ...BulkVerifyParams) ([]Verification, error) {
	for i := range params {
		params[i].Email = nymeria.Normalize(params[i].Email)
	}

	if len(params) == 0 {
		return nil, nymeria.ErrInvalidParameters
	}

	requests := []map[string]interface{}{}

	for _, p := range params {
		requests = append(requests, map[string]interface{}{
			"params": map[string]interface{}{
				"email": p.Email,
			},
			"metadata": p.MetaData,
		})
	}

	bs, err := json.Marshal(map[string]interface{}{
		"requests": requests,
	})

	if err != nil {
		return nil, err
	}

	req, err := nymeria.Request("POST", "/email/verify/bulk", bytes.NewBuffer(bs))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := nymeria.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if e, ok := nymeria.ErrMap[resp.StatusCode]; ok {
			return nil, e
		}

		return nil, nymeria.ErrServerError
	}

	defer resp.Body.Close()

	bs, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response []struct {
		Status int          `json:"status"`
		Data   Verification `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	var records []Verification

	for _, v := range response {
		records = append(records, v.Data)
	}

	return records, nil
}
