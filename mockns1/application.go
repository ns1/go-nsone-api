package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
)

// AddApplicationTestCase sets up a test case for the api.Client.application.List()
// function
func (s *Service) AddApplicationTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*pulsar.Application,
) error {
	return s.AddTestCase(
		http.MethodGet, "/pulsar/apps", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddApplicationGetTestCase sets up a test case for the api.Client.application.Get()
// function
func (s *Service) AddApplicationGetTestCase(
	id string,
	requestHeaders, responseHeaders http.Header,
	response *pulsar.Application,
) error {
	return s.AddTestCase(
		http.MethodGet, "/pulsar/apps/"+id, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddApplicationCreateTestCase sets up a test case for the api.Client.application.Create()
// function
func (s *Service) AddApplicationCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	application, response *pulsar.Application,
) error {
	return s.AddTestCase(
		http.MethodPut, "/pulsar/apps", http.StatusCreated, requestHeaders,
		responseHeaders, application, response,
	)
}

// AddApplicationUpdateTestCase sets up a test case for the api.Client.application.Update()
// function
func (s *Service) AddApplicationUpdateTestCase(
	requestHeaders, responseHeaders http.Header,
	application, response *pulsar.Application,
) error {
	return s.AddTestCase(
		http.MethodPost, "/pulsar/apps/"+application.ID, http.StatusOK, requestHeaders,
		responseHeaders, application, response,
	)
}

// AddApplicationDeleteTestCase sets up a test case for the api.Client.application.Delete()
// function
func (s *Service) AddApplicationDeleteTestCase(
	id string, requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, "/pulsar/apps/"+id, http.StatusNoContent, requestHeaders,
		responseHeaders, "", "",
	)
}
