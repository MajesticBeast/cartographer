package cartographer

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTFVersions(t *testing.T) {
	jsonResponse := `{
		"data": [
			{
				"attributes": {
					"version": "0.12.0",
					"workspace-count": 2,
					"workspaces": "workspace1,workspace2"
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

	tfVersions, err := c.TFVersions([]TFVersionFilter{})
	if err != nil {
		t.Errorf("TFVersions() returned an error: %v", err)
	}

	if len(tfVersions) != 1 {
		t.Errorf("TFVersions() returned %v tfVersions, expected 1", len(tfVersions))
	}

	tfVersion := tfVersions[0]
	if tfVersion.Version != "0.12.0" {
		t.Errorf("TFVersions() returned tfVersion with version %v, expected 'test'", tfVersion.Version)
	}
	if tfVersion.WorkspaceCount != 2 {
		t.Errorf("TFVersions() returned tfVersion with workspace count %v, expected 2", tfVersion.WorkspaceCount)
	}
	if tfVersion.Workspaces != "workspace1,workspace2" {
		t.Errorf("TFVersions() returned tfVersion with workspaces %v, expected 'test'", tfVersion.Workspaces)
	}
}
