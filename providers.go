package cartographer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ProviderName ProviderFilterType = iota
	ProviderSource
	ProviderVersion
	ProviderRegistryType
	ProviderWorkspaceCount
	ProviderWorkspaces
)

type ProviderFilterType int

func (p ProviderFilterType) String() string {
	return [...]string{"name", "source", "version", "registry-type", "workspace-count", "workspaces"}[p]
}

type ProviderFilter struct {
	Type     ProviderFilterType
	Operator FilterOperator
	Value    string
}

// Providers Retrieve a list of providers across all workspaces in an organization.
func (c *Cartographer) Providers(filters []ProviderFilter) ([]Provider, error) {
	var providers []Provider

	baseUrl, err := buildExplorerUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "providers")
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

		var apiResponse providerApiResponse
		if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}

		for _, item := range apiResponse.Data {
			providers = append(providers, item.Attributes)
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

		preventRateLimiting(apiResponse.Meta.Pagination.TotalPages)
	}

	return providers, nil
}

// Provider represents a Terraform Cloud provider.
type Provider struct {
	Name           string `json:"name"`
	Source         string `json:"source"`
	Version        string `json:"version"`
	RegistryType   string `json:"registry-type"`
	WorkspaceCount int    `json:"workspace-count"`
	Workspaces     string `json:"workspaces"`
}

// providerApiResponse is the response from the Terraform Cloud API for the providers endpoint.
type providerApiResponse struct {
	Data []struct {
		Attributes struct {
			Name           string `json:"name"`
			Source         string `json:"source"`
			Version        string `json:"version"`
			RegistryType   string `json:"registry-type"`
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
