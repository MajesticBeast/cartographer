package cartographer

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestModules(t *testing.T) {
	jsonResponse := `{
		"data": [
			{
				"attributes": {
					"name": "testname",
					"source": "testsource",
					"version": "testversion",
					"registry-type": "testtype",
					"workspace-count": 2,
					"workspaces": "testworkspaces"
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

	modules, err := c.Modules([]ModuleFilter{})
	if err != nil {
		t.Errorf("Modules() returned an error: %v", err)
	}

	if len(modules) != 1 {
		t.Errorf("Modules() returned %v modules, expected 1", len(modules))
	}

	module := modules[0]
	if module.Name != "testname" {
		t.Errorf("Modules() returned module with name %v, expected 'test'", module.Name)
	}
	if module.Source != "testsource" {
		t.Errorf("Modules() returned module with source %v, expected 'test'", module.Source)
	}
	if module.Version != "testversion" {
		t.Errorf("Modules() returned module with version %v, expected 'test'", module.Version)
	}
	if module.RegistryType != "testtype" {
		t.Errorf("Modules() returned module with registry type %v, expected 'test'", module.RegistryType)
	}
	if module.WorkspaceCount != 2 {
		t.Errorf("Modules() returned module with workspace count %v, expected 2", module.WorkspaceCount)
	}
	if module.Workspaces != "testworkspaces" {
		t.Errorf("Modules() returned module with workspaces %v, expected 'test'", module.Workspaces)
	}
}
