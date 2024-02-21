package cartographer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	moduleName ModuleFilterType = iota
	moduleSource
	moduleVersion
	moduleWorkspaceCount
	moduleInWorkspaces
)

type ModuleFilterType int

func (m ModuleFilterType) String() string {
	return [...]string{"name", "source", "version", "workspace-count", "workspaces"}[m]
}

type ModuleFilter struct {
	Type     ModuleFilterType
	Operator FilterOperator
	Value    string
}

// Modules Retrieve a list of modules across all workspaces in an organization. It takes an http.Client, the name of the
// organization, and a Terraform Cloud API token as arguments. If the request fails, it returns an error. If the request
// is successful, it returns a slice of module.
func (c *Cartographer) Modules(filters []ModuleFilter) (ModuleList, error) {
	var modules ModuleList

	baseUrl, err := buildUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "modules")

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

		var apiResponse modulesApiResponse
		if err = json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}

		for _, item := range apiResponse.Data {
			modules = append(modules, item.Attributes)
		}

		if apiResponse.Meta.Pagination.NextPage == nil {
			break
		}

		req.URL, err = url.Parse(apiResponse.Links.Next.(string))
		if err != nil {
			return nil, err
		}
		res.Body.Close()
	}

	return modules, nil
}

// ModuleList is a slice of module
type ModuleList []Module

// Module represents a module in Terraform Cloud
type Module struct {
	Name           string `json:"name"`
	Source         string `json:"source"`
	Version        string `json:"version"`
	RegistryType   string `json:"registry-type"`
	WorkspaceCount int    `json:"workspace-count"`
	Workspaces     string `json:"workspaces"`
}

// modulesApiResponse is the response from the Terraform Cloud API
type modulesApiResponse struct {
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
