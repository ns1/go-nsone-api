package mockns1

import (
	"net/http"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

// AddActivityListTestCase sets up a test case for the api.Client.Activity.List()
// function
func (s *Service) AddActivityListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*account.Activity,
	params ...api.Param,
) error {
	return s.AddTestCase(
		http.MethodGet, "/account/activity", http.StatusOK, requestHeaders,
		responseHeaders, "", response, params...,
	)
}
