package cartographer

import (
	"net/http"
	"testing"
)

func TestNewCartographer(t *testing.T) {
	orgName := "testOrg"
	token := "testToken"

	cartographer := NewCartographer(orgName, token)

	if cartographer.orgName != orgName {
		t.Errorf("Expected orgName to be %s, but got %s", orgName, cartographer.orgName)
	}

	if cartographer.token != token {
		t.Errorf("Expected token to be %s, but got %s", token, cartographer.token)
	}
}

func TestBuildUrl(t *testing.T) {
	orgName := "testOrg"
	expectedUrl := "https://app.terraform.io/api/v2/organizations/testOrg/explorer"

	url, err := buildUrl(orgName)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if url.String() != expectedUrl {
		t.Errorf("Expected URL to be %s, but got %s", expectedUrl, url.String())
	}
}

func TestCheckStatusCode(t *testing.T) {
	res := &http.Response{
		StatusCode: 200,
	}

	err := checkStatusCode(res)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}
