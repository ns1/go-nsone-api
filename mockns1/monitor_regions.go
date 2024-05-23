package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

// AddMonitorRegionsListTestCase sets up a test case for the api.Client.MonitorRegionsService.List() function.
func (s *Service) AddMonitorRegionsListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*monitor.Region,
) error {
	return s.AddTestCase(
		http.MethodGet, "/monitoring/regions", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}
