package mockns1

import (
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// AddZoneListTestCase sets up a test case for the api.Client.Versions.List()
// function
func (s *Service) AddVersionListTestCase(
	zoneName string,
	requestHeaders, responseHeaders http.Header,
	response []*dns.Version,
) error {
	return s.AddTestCase(
		http.MethodGet, fmt.Sprintf("/zones/%s/versions", zoneName), http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddCreateVersionTestCase sets up a test case for the api.Client.Versions.Create()
// function
func (s *Service) AddCreateVersionTestCase(
	zoneName string,
	requestHeaders, responseHeaders http.Header,
	response *dns.Version,
) error {
	return s.AddTestCase(
		http.MethodPut, fmt.Sprintf("/zones/%s/versions?force=false", zoneName), http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDeleteVersionTestCase sets up a test case for the api.Client.Versions.Delete()
// function
func (s *Service) AddDeleteVersionTestCase(
	zoneName string,
	versionID int64,
	requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, fmt.Sprintf("/zones/%s/versions/%d", zoneName, versionID), http.StatusOK, requestHeaders,
		responseHeaders, "", nil,
	)
}

// AddActivateVersionTestCase sets up a test case for the api.Client.Versions.Activate()
// function
func (s *Service) AddActivateVersionTestCase(
	zoneName string,
	versionID int64,
	requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodPost, fmt.Sprintf("/zones/%s/versions/%d/activate", zoneName, versionID), http.StatusOK, requestHeaders,
		responseHeaders, "", nil,
	)
}
