package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	Is FilterOperator = iota
	IsNot
	Contains
	DoesNotContain
	IsEmpty
	IsNotEmpty
	Gt
	Lt
	Gteq
	Lteq
	IsBefore
	IsAfter
)

type FilterOperator int

func (f FilterOperator) String() string {
	return [...]string{"is", "is-not", "contains", "does-not-contain", "is-empty", "is-not-empty", "gt", "lt", "gteq", "lteq", "is-before", "is-after"}[f]
}

type Cartographer struct {
	client  *http.Client
	orgName string
	token   string
}

func NewCartographer(orgName string, token string) *Cartographer {
	return &Cartographer{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		orgName: orgName,
		token:   token,
	}
}

// buildUrl Builds the URL for the Terraform Cloud API. It takes the organization name as an argument and returns the
// formatted URL string.
func buildUrl(orgName string) (*url.URL, error) {
	baseURL, err := url.Parse(fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/explorer", orgName))
	if err != nil {
		return nil, err
	}
	return baseURL, nil
}

// CheckStatusCode checks the status code of the response. If the status code is 429, it returns an error indicating
// that the request was rate limited. If the status code is not in the 200 range, it returns an error indicating
// the status code.
func CheckStatusCode(res *http.Response) error {
	if res.StatusCode == 429 {
		return fmt.Errorf("rate limited - https://developer.hashicorp.com/terraform/cloud-docs/api-docs#rate-limiting")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func main() {
	c := NewCartographer("thelostsons", os.Getenv("TFTOKEN"))

	var moduleFilters []ModuleFilter
	moduleFilters = append(moduleFilters, ModuleFilter{
		Type:     ModuleName,
		Operator: Contains,
		Value:    "iam",
	})

	mods, err := c.Modules(moduleFilters)
	if err != nil {
		log.Fatal(err)
	}

	for _, mod := range mods {
		log.Println(mod.Name, mod.Source, mod.Version, mod.RegistryType, mod.WorkspaceCount, mod.Workspaces)
	}

	var workspaceFilters []WorkspaceFilter
	workspaceFilters = append(workspaceFilters, WorkspaceFilter{
		Type:     WorkspaceName,
		Operator: Contains,
		Value:    "lostsons",
	})

	var providerFilters []ProviderFilter
	providerFilters = append(providerFilters, ProviderFilter{
		Type:     ProviderName,
		Operator: Contains,
		Value:    "aws",
	})

	works, err := c.Workspaces(workspaceFilters)
	if err != nil {
		log.Fatal(err)
	}

	provs, err := c.Providers(providerFilters)
	if err != nil {
		log.Fatal(err)
	}

	for _, work := range works {
		log.Println(work.WorkspaceName, work.WorkspaceCreatedAt, work.WorkspaceUpdatedAt, work.Modules, work.ModuleCount)
	}

	for _, prov := range provs {
		log.Println(prov.Name, prov.Source, prov.Version, prov.RegistryType, prov.WorkspaceCount, prov.Workspaces)
	}

}
