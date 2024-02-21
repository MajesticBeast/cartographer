package cartographer

import (
	"fmt"
	"net/http"
	"net/url"
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

// checkStatusCode checks the status code of the response. If the status code is 429, it returns an error indicating
// that the request was rate limited. If the status code is not in the 200 range, it returns an error indicating
// the status code.
func checkStatusCode(res *http.Response) error {
	if res.StatusCode == 429 {
		return fmt.Errorf("rate limited - https://developer.hashicorp.com/terraform/cloud-docs/api-docs#rate-limiting")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

//
//func main() {
//	c := NewCartographer("thelostsons", os.Getenv("TFTOKEN"))
//
//	var moduleFilters []moduleFilter
//	moduleFilters = append(moduleFilters, moduleFilter{
//		Type:     moduleName,
//		Operator: Contains,
//		Value:    "iam",
//	})
//
//	var workspaceFilters []workspaceFilter
//	workspaceFilters = append(workspaceFilters, workspaceFilter{
//		Type:     workspaceName,
//		Operator: Contains,
//		Value:    "lostsons",
//	})
//
//	var providerFilters []providerFilter
//	providerFilters = append(providerFilters, providerFilter{
//		Type:     providerName,
//		Operator: Contains,
//		Value:    "aws",
//	})
//
//	var tfVersionFilters []tfVersionFilter
//	tfVersionFilters = append(tfVersionFilters, tfVersionFilter{
//		Type:     tfVersionVersion,
//		Operator: Contains,
//		Value:    "0.12",
//	})
//
//	mods, err := c.modules(moduleFilters)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	works, err := c.workspaces(workspaceFilters)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	provs, err := c.providers(providerFilters)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tfVers, err := c.tfVersions()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, mod := range mods {
//		log.Println(mod.Name, mod.Source, mod.Version, mod.RegistryType, mod.WorkspaceCount, mod.Workspaces)
//	}
//
//	for _, work := range works {
//		log.Println(work.WorkspaceName, work.WorkspaceCreatedAt, work.WorkspaceUpdatedAt, work.Modules, work.ModuleCount)
//	}
//
//	for _, prov := range provs {
//		log.Println(prov.Name, prov.Source, prov.Version, prov.RegistryType, prov.WorkspaceCount, prov.Workspaces)
//	}
//
//	for _, tfVer := range tfVers {
//		log.Println(tfVer.Version, tfVer.WorkspaceCount, tfVer.Workspaces)
//	}
//
//}
