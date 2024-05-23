package mockns1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/testify/assert"
	api "gopkg.in/ns1/ns1-go.v2/rest"
)

type testCase struct {
	status  int
	request struct {
		headers http.Header
		body    []byte
		json    bool
	}
	response struct {
		headers http.Header
		body    []byte
	}
}

// AddTestCase adds a new test case to the mock service. Test cases are
// unique based on the method, uri, request headers, and request body.
func (s *Service) AddTestCase(
	method, uri string, returnStatus int,
	requestHeaders, responseHeaders http.Header,
	requestBody, responseBody interface{},
	params ...api.Param,
) error {
	s.stopTimer()
	defer s.startTimer()

	if !strings.HasPrefix(uri, "/v1/") {
		uri = "/v1/" + uri
	}

	if len(params) > 0 {
		uri = fmt.Sprintf("%s?%s=%s", uri, params[0].Key, params[0].Value)

		for _, p := range params[1:] {
			uri = fmt.Sprintf("%s&%s=%s", uri, p.Key, p.Value)
		}
	}

	uri = strings.Replace(uri, "//", "/", -1)

	tc := &testCase{
		status: returnStatus,
	}
	tc.request.headers = requestHeaders
	tc.response.headers = responseHeaders

	var err error
	if tc.request.body, tc.request.json, err = convertBody(requestBody); err != nil {
		return fmt.Errorf("unable to convert request body to []byte: %s", err)
	}
	if tc.response.body, _, err = convertBody(responseBody); err != nil {
		return fmt.Errorf("unable to convert response body to []byte: %s", err)
	}

	if _, exists := s.tests[method]; !exists {
		s.tests[method] = map[string][]*testCase{}
	}
	if _, exists := s.tests[method][uri]; !exists {
		s.tests[method][uri] = []*testCase{}
	}

	t := new(testifyT)
	for _, test := range s.tests[method][uri] {
		header := assert.Equal(t, tc.request.headers, test.request.headers)
		body := assert.Equal(t, tc.request.body, test.request.body)

		if header && body {
			return errors.New("test case already registered")
		}
	}

	s.tests[method][uri] = append(s.tests[method][uri], tc)

	return nil
}

// ClearTestCases removes all previously added test cases
func (s *Service) ClearTestCases() {
	s.tests = map[string]map[string][]*testCase{}
}

func convertBody(body interface{}) ([]byte, bool, error) {
	switch b := body.(type) {
	case []byte:
		return b, false, nil
	case string:
		return []byte(b), false, nil
	}

	data, err := json.Marshal(body)
	return data, true, err
}
