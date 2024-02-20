package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func buildUrl(orgName string) string {
	return fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/explorer?type=modules", orgName)
}

type Module struct {
	Name           string `json:"name"`
	Source         string `json:"source"`
	Version        string `json:"version"`
	RegistryType   string `json:"registry-type"`
	WorkspaceCount int    `json:"workspace-count"`
	Workspaces     string `json:"workspaces"`
}

// Modules Retrieve a list of modules across all workspaces
func Modules(client *http.Client, orgName string, token string) ([]Module, error) {
	var modules []Module

	req, err := http.NewRequest("GET", buildUrl(orgName), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	for {
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if err := CheckStatusCode(res); err != nil {
			return nil, err
		}

		var apiResponse apiResponse
		if err = json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}

		for _, item := range apiResponse.Data {
			modules = append(modules, item.Attributes)
		}

		if apiResponse.Meta.Pagination.CurrentPage == apiResponse.Meta.Pagination.TotalPages {
			break
		}

		req.URL, err = url.Parse(apiResponse.Links.Next.(string))
		res.Body.Close()
	}

	return modules, nil
}

func CheckStatusCode(res *http.Response) error {
	if res.StatusCode == 429 {
		return fmt.Errorf("rate limited - https://developer.hashicorp.com/terraform/cloud-docs/api-docs#rate-limiting")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return nil
}

type apiResponse struct {
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
