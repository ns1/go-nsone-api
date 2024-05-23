package rest_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/monitor"
)

func TestMonitorRegionsService(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			mockRegions := []*monitor.Region{
				{
					Code:    "lga",
					Name:    "New York",
					Subnets: []string{"1.2.3.4/24"},
				},
				{
					Code:    "sjc",
					Name:    "San Jose",
					Subnets: []string{"5.6.7.8/24"},
				},
			}

			require.Nil(t, mock.AddMonitorRegionsListTestCase(nil, nil, mockRegions))

			regions, resp, err := client.MonitorRegions.List()

			require.Nil(t, err)
			require.Equal(t, 200, resp.StatusCode)
			require.Len(t, regions, len(mockRegions))
		})
	})
}
