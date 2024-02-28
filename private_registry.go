package cartographer

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// PrivateRegistryModules retrieves a list of all modules in the given organization's private registry. It only returns
// the latest version of the module. This is currently used by indexing to 0 which as of now is the latest version.
// No comparison is done to check if the version is the latest.
func (c *Cartographer) PrivateRegistryModules() ([]PrivateRegistryModule, error) {
	var modules []PrivateRegistryModule

	baseUrl, err := buildRegistryUrl(c.orgName)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("page[size]", "100")

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

		var apiResponse privateRegistryApiResponse
		err = json.NewDecoder(res.Body).Decode(&apiResponse)
		if err != nil {
			return nil, err
		}

		for _, registry := range apiResponse.Data {
			modules = append(modules, PrivateRegistryModule{
				Id:            registry.Id,
				Type:          registry.Type,
				Name:          registry.Attributes.Name,
				Status:        registry.Attributes.Status,
				LatestVersion: registry.Attributes.VersionStatuses[0].Version,
				UpdatedAt:     registry.Attributes.UpdatedAt,
				CreatedAt:     registry.Attributes.CreatedAt,
			})
		}

		if apiResponse.Meta.Pagination.NextPage == nil {
			break
		}

		if apiResponse.Links.Next == nil {
			break
		}

		req.URL, err = url.Parse(apiResponse.Links.Next.(string))
		if err != nil {
			return nil, err
		}
		res.Body.Close()

		preventRateLimiting(apiResponse.Meta.Pagination.TotalPages)
	}

	return modules, nil
}

type PrivateRegistryModule struct {
	Id            string    `json:"id"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	LatestVersion string    `json:"latest_version"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type privateRegistryApiResponse struct {
	Data []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			Provider        string `json:"provider"`
			Status          string `json:"status"`
			VersionStatuses []struct {
				Version string `json:"version"`
				Status  string `json:"status"`
			} `json:"version-statuses"`
			CreatedAt           time.Time `json:"created-at"`
			UpdatedAt           time.Time `json:"updated-at"`
			RegistryName        string    `json:"registry-name"`
			NoCode              bool      `json:"no-code"`
			PublishingMechanism string    `json:"publishing-mechanism"`
			VcsRepo             struct {
				Branch                  string      `json:"branch"`
				IngressSubmodules       bool        `json:"ingress-submodules"`
				TagsRegex               interface{} `json:"tags-regex"`
				Identifier              string      `json:"identifier"`
				DisplayIdentifier       string      `json:"display-identifier"`
				GithubAppInstallationId string      `json:"github-app-installation-id"`
				RepositoryHttpUrl       string      `json:"repository-http-url"`
				ServiceProvider         string      `json:"service-provider"`
				Tags                    bool        `json:"tags"`
			} `json:"vcs-repo"`
			Permissions struct {
				CanDelete bool `json:"can-delete"`
				CanResync bool `json:"can-resync"`
				CanRetry  bool `json:"can-retry"`
			} `json:"permissions"`
		} `json:"attributes"`
		Relationships struct {
			Organization struct {
				Data struct {
					Id   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"organization"`
			NoCodeModules struct {
				Data []interface{} `json:"data"`
			} `json:"no-code-modules"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
	Links struct {
		Self  string      `json:"self"`
		First string      `json:"first"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
		Last  string      `json:"last"`
	} `json:"links"`
	Meta struct {
		Pagination struct {
			CurrentPage int         `json:"current-page"`
			PageSize    int         `json:"page-size"`
			PrevPage    interface{} `json:"prev-page"`
			NextPage    interface{} `json:"next-page"`
			TotalPages  int         `json:"total-pages"`
			TotalCount  int         `json:"total-count"`
		} `json:"pagination"`
	} `json:"meta"`
}
