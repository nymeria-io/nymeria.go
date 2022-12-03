package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nymeriaio/nymeria.go"
	"github.com/nymeriaio/nymeria.go/internal/api"
)

type SearchParams struct {
	Query string
	Size  int /* how many records to retrieve starting at From (default: 0) */
	From  int /* from which record to start */
}

func (s SearchParams) Invalid() bool {
	return len(s.Query) == 0
}

func (s SearchParams) URL() string {
	return fmt.Sprintf(
		"query=%s&size=%s&form=%s",
		url.QueryEscape(s.Query),
		url.QueryEscape(fmt.Sprintf("%d", s.Size)),
		url.QueryEscape(fmt.Sprintf("%d", s.From)),
	)
}

func Search(params SearchParams) ([]Person, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/person/search?%s", params.URL()), nil)

	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)

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

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response struct {
		Data   []Person `json:"data"`
		Status int      `json:"status"`
		Total  int      `json:"total"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
