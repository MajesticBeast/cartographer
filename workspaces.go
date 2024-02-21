package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	WorkspaceAllChecksSucceeded WorkspaceFilterType = iota
	WorkspaceChecksErrored
	WorkspaceChecksFailed
	WorkspaceChecksPassed
	WorkspaceChecksUnknown
	WorkspaceCurrentRunAppliedAt
	WorkspaceCurrentRunExternalId
	WorkspaceCurrentRunStatus
	WorkspaceDrifted
	WorkspaceExternalId
	WorkspaceModuleCount
	WorkspaceModulesInWorkspace
	WorkspaceOrganizationName
	WorkspaceProjectExternalId
	WorkspaceProjectName
	WorkspaceProviderCount
	WorkspaceProviders
	WorkspaceResourcesDrifted
	WorkspaceResourcesUndrifted
	WorkspaceStateVersionTerraformVersion
	WorkspaceVcsRepoIdentifier
	WorkspaceCreatedAt
	WorkspaceName
	WorkspaceTerraformVersion
	WorkspaceUpdatedAt
)

type WorkspaceFilterType int

type WorkspaceFilter struct {
	Type     WorkspaceFilterType
	Operator FilterOperator
	Value    string
}

func (w WorkspaceFilterType) String() string {
	return [...]string{
		"all-checks-succeeded",
		"checks-errored",
		"checks-failed",
		"checks-passed",
		"checks-unknown",
		"current-run-applied-at",
		"current-run-external-id",
		"current-run-status",
		"drifted",
		"external-id",
		"module-count",
		"modules",
		"organization-name",
		"project-external-id",
		"project-name",
		"provider-count",
		"providers",
		"resources-drifted",
		"resources-undrifted",
		"state-version-terraform-version",
		"vcs-repo-identifier",
		"workspace-created-at",
		"workspace-name",
		"workspace-terraform-version",
		"workspace-updated-at",
	}[w]
}

type WorkspaceList []Workspace

// Workspaces Retrieve a list of workspaces in an organization. It takes an http.Client, the name of the organization,
// and a Terraform Cloud API token as arguments. If the request fails, it returns an error. If the request is successful,
// it returns a slice of Workspace.
func (c *Cartographer) Workspaces(filters []WorkspaceFilter) (WorkspaceList, error) {
	var workspaces WorkspaceList

	baseUrl, err := buildUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("type", "workspaces")

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

		var apiResponse workspacesApiResponse
		if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, err
		}

		for _, item := range apiResponse.Data {
			workspaces = append(workspaces, item.Attributes)
		}

		if apiResponse.Meta.Pagination.NextPage == nil {
			break
		}

		req.URL, err = url.Parse(apiResponse.Links.Next.(string))
		res.Body.Close()
	}

	return workspaces, nil
}

// Workspace represents a workspace in Terraform Cloud
type Workspace struct {
	AllChecksSucceeded           bool        `json:"all-checks-succeeded"`
	ChecksErrored                int         `json:"checks-errored"`
	ChecksFailed                 int         `json:"checks-failed"`
	ChecksPassed                 int         `json:"checks-passed"`
	ChecksUnknown                int         `json:"checks-unknown"`
	CurrentRunAppliedAt          *time.Time  `json:"current-run-applied-at"`
	CurrentRunExternalId         string      `json:"current-run-external-id"`
	CurrentRunStatus             string      `json:"current-run-status"`
	Drifted                      bool        `json:"drifted"`
	ExternalId                   string      `json:"external-id"`
	ModuleCount                  int         `json:"module-count"`
	Modules                      interface{} `json:"modules"`
	OrganizationName             string      `json:"organization-name"`
	ProjectExternalId            string      `json:"project-external-id"`
	ProjectName                  string      `json:"project-name"`
	ProviderCount                int         `json:"provider-count"`
	Providers                    interface{} `json:"providers"`
	ResourcesDrifted             int         `json:"resources-drifted"`
	ResourcesUndrifted           int         `json:"resources-undrifted"`
	StateVersionTerraformVersion string      `json:"state-version-terraform-version"`
	VcsRepoIdentifier            *string     `json:"vcs-repo-identifier"`
	WorkspaceCreatedAt           time.Time   `json:"workspace-created-at"`
	WorkspaceName                string      `json:"workspace-name"`
	WorkspaceTerraformVersion    string      `json:"workspace-terraform-version"`
	WorkspaceUpdatedAt           time.Time   `json:"workspace-updated-at"`
}

// workspacesApiResponse is the response from the Terraform Cloud API
type workspacesApiResponse struct {
	Data []struct {
		Attributes struct {
			AllChecksSucceeded           bool        `json:"all-checks-succeeded"`
			ChecksErrored                int         `json:"checks-errored"`
			ChecksFailed                 int         `json:"checks-failed"`
			ChecksPassed                 int         `json:"checks-passed"`
			ChecksUnknown                int         `json:"checks-unknown"`
			CurrentRunAppliedAt          *time.Time  `json:"current-run-applied-at"`
			CurrentRunExternalId         string      `json:"current-run-external-id"`
			CurrentRunStatus             string      `json:"current-run-status"`
			Drifted                      bool        `json:"drifted"`
			ExternalId                   string      `json:"external-id"`
			ModuleCount                  int         `json:"module-count"`
			Modules                      interface{} `json:"modules"`
			OrganizationName             string      `json:"organization-name"`
			ProjectExternalId            string      `json:"project-external-id"`
			ProjectName                  string      `json:"project-name"`
			ProviderCount                int         `json:"provider-count"`
			Providers                    interface{} `json:"providers"`
			ResourcesDrifted             int         `json:"resources-drifted"`
			ResourcesUndrifted           int         `json:"resources-undrifted"`
			StateVersionTerraformVersion string      `json:"state-version-terraform-version"`
			VcsRepoIdentifier            *string     `json:"vcs-repo-identifier"`
			WorkspaceCreatedAt           time.Time   `json:"workspace-created-at"`
			WorkspaceName                string      `json:"workspace-name"`
			WorkspaceTerraformVersion    string      `json:"workspace-terraform-version"`
			WorkspaceUpdatedAt           time.Time   `json:"workspace-updated-at"`
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
