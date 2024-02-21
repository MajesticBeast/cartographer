package main

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

type ProviderList []Provider

type ProviderFilterType int

func (p ProviderFilterType) String() string {
	return [...]string{"name", "source", "version", "registry-type", "workspace-count", "workspaces"}[p]
}

type ProviderFilter struct {
	Type     ProviderFilterType
	Operator FilterOperator
	Value    string
}

// Providers Retrieve a list of providers across all workspaces in an organization. It takes an http.Client, the name of
// the organization, and a Terraform Cloud API token as arguments. If the request fails, it returns an error. If the
// request is successful, it returns a slice of Provider.
func (c *Cartographer) Providers(filters []ProviderFilter) (ProviderList, error) {
	var providers ProviderList

	baseUrl, err := buildUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "providers")

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

		if err := CheckStatusCode(res); err != nil {
			return nil, err
		}

		var response providerApiResponse
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, err
		}

		for _, provider := range response.Data {
			providers = append(providers, Provider{
				Name:           provider.Attributes.Name,
				Source:         provider.Attributes.Source,
				Version:        provider.Attributes.Version,
				RegistryType:   provider.Attributes.RegistryType,
				WorkspaceCount: provider.Attributes.WorkspaceCount,
				Workspaces:     provider.Attributes.Workspaces,
			})
		}

		if response.Links.Next == nil {
			break
		}

		req, err = http.NewRequest("GET", response.Links.Next.(string), nil)
		if err != nil {
			return nil, err
		}
	}

	return providers, nil
}

// Provider represents a Terraform Cloud provider.
type Provider struct {
	Name           string
	Source         string
	Version        string
	RegistryType   string
	WorkspaceCount int
	Workspaces     string
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
