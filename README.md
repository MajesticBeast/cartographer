# Cartographer
Cartographer is a Go client for the Terraform Cloud Explorer API. The official Go client for the TF Cloud API does not seem to include the Explorer endpoint.

---
### Example

```go
package main

import (
	"fmt"
	carto "github.com/majesticbeast/cartographer"
	"os"
)

func main() {
	c := carto.NewCartographer(os.Getenv("ORG_NAME"), os.Getenv("TFTOKEN"))

	var moduleFilters = []carto.ModuleFilter{
		{
			Type:     carto.ModuleName,
			Operator: carto.Contains,
			Value:    "iam",
		},
		{
			Type:     carto.ModuleSource,
			Operator: carto.Contains,
			Value:    "partial-source-path",
		},
	}
	
	var workspaceFilters []carto.WorkspaceFilter
	var workspaceFilter = carto.WorkspaceFilter{
		Type:     carto.WorkspaceName,
		Operator: carto.Is,
		Value:    "a-workspace-name"
	}

	workspaceFilters = append(workspaceFilters, workspaceFilter)

	modules, err := c.Modules(moduleFilters)
	if err != nil {
		fmt.Println(err)
	}

	workspaces, err := c.Workspaces(workspaceFilters)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(modules)
	for _, workspace := range workspaces {
		fmt.Println(workspace.WorkspaceName, workspace.WorkspaceCreatedAt, workspace.WorkspaceUpdatedAt, workspace.Modules)
	}
}
```
