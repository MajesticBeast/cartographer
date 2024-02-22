package cartographer

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWorkspaces(t *testing.T) {
	jsonResponse := `{
		"data": [
			{
				"attributes": {
					"all-checks-succeeded": true,
					"checks-errored": 0,
					"checks-failed": 0,
					"checks-passed": 5,
					"checks-unknown": 1,
					"current-run-applied-at": "2019-07-23T15:00:00Z", 
					"current-run-external-id": "testID",
					"current-run-status": "applied",
					"drifted": false,
					"external-id": "testID",
					"module-count": 3,
					"modules": "iam,s3,tags",
					"organization-name": "myOrgName",
					"project-external-id": "testID",
					"project-name": "testProjectName",
					"provider-count": 2,
					"providers": "aws,github",
					"resources-drifted": 0,
					"resources-undrifted": 3,
					"state-version-terraform-version": "0.12.0", 
					"vcs-repo-identifier": "github.com/majesticbeast/cartographer",
					"workspace-created-at": "2019-07-23T15:00:00Z",
					"workspace-name": "testWorkspace",
					"workspace-terraform-version": "0.12.0",
					"workspace-updated-at": "2019-07-23T15:00:00Z"
				},
				"id": "test",
				"type": "test"
			}
		],
		"links": {
			"self": "test",
			"first": "test",
			"last": "test",
			"prev": null,
			"next": null
		},
		"meta": {
			"pagination": {
				"current-page": 1,
				"page-size": 1,
				"next-page": null,
				"prev-page": null,
				"total-pages": 1,
				"total-count": 1
			}
		}
	}`

	mockClient := &MockClient{
		MockDo: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(jsonResponse)),
			}, nil
		},
	}

	c := &Cartographer{
		client:  mockClient,
		orgName: "test",
		token:   "test",
	}

	workspaces, err := c.Workspaces([]WorkspaceFilter{})
	if err != nil {
		t.Errorf("Workspaces() returned an error: %v", err)
	}

	if len(workspaces) != 1 {
		t.Errorf("Workspaces() returned %v tfVersions, expected 1", len(workspaces))
	}

	workspace := workspaces[0]
	if workspace.WorkspaceName != "testWorkspace" {
		t.Errorf("Workspaces() returned workspace with name %v, expected 'testWorkspace'", workspace.WorkspaceName)
	}
	if workspace.OrganizationName != "myOrgName" {
		t.Errorf("Workspaces() returned workspace with organization name %v, expected 'myOrgName'", workspace.OrganizationName)
	}
	if workspace.ProjectName != "testProjectName" {
		t.Errorf("Workspaces() returned workspace with project name %v, expected 'testProjectName'", workspace.ProjectName)
	}
	if workspace.ModuleCount != 3 {
		t.Errorf("Workspaces() returned workspace with module count %v, expected '3'", workspace.ModuleCount)
	}
}
