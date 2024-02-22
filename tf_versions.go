package cartographer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	TFVersionVersion TFVersionFilterType = iota
	TFVersionWorkspaceCount
	TFVersionWorkspaces
)

type TFVersionFilterType int

type TFVersionList []TFVersion

type TFVersionFilter struct {
	Type     TFVersionFilterType
	Operator FilterOperator
	Value    string
}

// TFVersionFilterType Stringer
func (c TFVersionFilterType) String() string {
	return [...]string{"version", "workspace-count", "workspaces"}[c]
}

// TFVersions Retrieve a list of Terraform versions across all workspaces in an organization. It takes an http.Client,
// the name of the organization, and a Terraform Cloud API token as arguments. If the request fails, it returns an error.
// If the request is successful, it returns a slice of TFVersion.
func (c *Cartographer) TFVersions(filters []TFVersionFilter) (TFVersionList, error) {
	var tfVersions TFVersionList

	baseUrl, err := buildUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "tf_versions")
	q.Add("page[size]", "100")

	for i, filter := range filters {
		key := fmt.Sprintf("filter[%d][%s][%s][0]", i, filter.Type.String(), filter.Operator.String())
		q.Add(key, filter.Value)
	}

	baseUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	for {
		res, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}

		if err := checkStatusCode(res); err != nil {
			return nil, err
		}

		var apiResponse tfVersionsApiResponse
		if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}

		for _, item := range apiResponse.Data {
			tfVersions = append(tfVersions, item.Attributes)
		}

		if apiResponse.Meta.Pagination.NextPage == nil {
			break
		}

		if apiResponse.Links.Next == nil {
			break
		}

		req.URL, err = url.Parse(apiResponse.Links.Next.(string))
		if err != nil {
			return nil, err
		}
		res.Body.Close()
	}

	return tfVersions, nil
}

// TFVersion represents a Terraform version.
type TFVersion struct {
	Version        string `json:"version"`
	WorkspaceCount int    `json:"workspace-count"`
	Workspaces     string `json:"workspaces"`
}

// tfVersionsApiResponse is the response from the Terraform Cloud API for the tf_versions endpoint.
type tfVersionsApiResponse struct {
	Data []struct {
		Attributes struct {
			Version        string `json:"version"`
			WorkspaceCount int    `json:"workspace-count"`
			Workspaces     string `json:"workspaces"`
		} `json:"attributes"`
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Links struct {
		Self  string      `json:"self"`
		First string      `json:"first"`
		Last  string      `json:"last"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
	} `json:"links"`
	Meta struct {
		Pagination struct {
			CurrentPage int         `json:"current-page"`
			PageSize    int         `json:"page-size"`
			NextPage    interface{} `json:"next-page"`
			PrevPage    interface{} `json:"prev-page"`
			TotalPages  int         `json:"total-pages"`
			TotalCount  int         `json:"total-count"`
		} `json:"pagination"`
	} `json:"meta"`
}
