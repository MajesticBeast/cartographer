package main

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

type TFVersionFilter struct {
	Type     TFVersionFilterType
	Operator FilterOperator
	Value    string
}

func (c TFVersionFilterType) String() string {
	return [...]string{"version", "workspace-count", "workspaces"}[c]
}

func (c *Cartographer) TFVersions() (TFVersionList, error) {
	var versions TFVersionList

	baseUrl, err := buildUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "tf_versions")

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

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		var apiResponse tfVersionsApiResponse
		err = json.NewDecoder(res.Body).Decode(&apiResponse)
		if err != nil {
			return nil, err
		}

		for _, tfVersion := range apiResponse.Data {
			versions = append(versions, TFVersion{
				Version:        tfVersion.Attributes.Version,
				WorkspaceCount: tfVersion.Attributes.WorkspaceCount,
				Workspaces:     tfVersion.Attributes.Workspaces,
			})
		}

		if apiResponse.Links.Next == nil {
			break
		}

		req, err = http.NewRequest("GET", apiResponse.Links.Next.(string), nil)
		if err != nil {
			return nil, err
		}
		res.Body.Close()
	}

	return versions, nil
}

type TFVersion struct {
	Version        string `json:"version"`
	WorkspaceCount int    `json:"workspace-count"`
	Workspaces     string `json:"workspaces"`
}

type TFVersionList []TFVersion

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
