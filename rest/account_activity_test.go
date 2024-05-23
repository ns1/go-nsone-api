package rest_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestActivity(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()
	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		activity := []*account.Activity{
			{
				UserID:       "user-1",
				ResourceID:   "resource-id-1",
				Timestamp:    567,
				UserType:     "apikey",
				Action:       "create",
				UserName:     "username-1",
				ID:           "id-1",
				ResourceType: "record",
			},
			{
				UserID:       "user-2",
				ResourceID:   "resource-id-2",
				Timestamp:    567,
				UserType:     "apikey",
				Action:       "delete",
				UserName:     "username-2",
				ID:           "id-2",
				ResourceType: "record",
			},
		}

		t.Run("list default activity", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddActivityListTestCase(nil, nil, activity))

			respActivity, _, err := client.Activity.List()
			require.Nil(t, err)
			require.NotNil(t, respActivity)
			require.Equal(t, len(activity), len(respActivity))

			for i := range activity {
				require.Equal(t, activity[i], respActivity[i], i)
			}
		})

		t.Run("list most recent 1 activity", func(t *testing.T) {
			defer mock.ClearTestCases()

			limit := 1
			params := []api.Param{{Key: "limit", Value: fmt.Sprintf("%d", limit)}}

			require.Nil(t, mock.AddActivityListTestCase(nil, nil, []*account.Activity{activity[0]}, params...))

			respActivity, _, err := client.Activity.List(params...)
			require.Nil(t, err)
			require.NotNil(t, respActivity)
			require.Equal(t, limit, len(respActivity))
			require.Equal(t, activity[0], respActivity[0])
		})

		t.Run("list all dns record activity, multiple params", func(t *testing.T) {
			defer mock.ClearTestCases()

			limit := 1000
			params := []api.Param{{Key: "limit", Value: fmt.Sprintf("%d", limit)}, {Key: "resource_type", Value: "record"}}

			require.Nil(t, mock.AddActivityListTestCase(nil, nil, activity, params...))

			respActivity, _, err := client.Activity.List(params...)
			require.Nil(t, err)
			require.NotNil(t, respActivity)
			require.Equal(t, len(activity), len(respActivity))

			for i := range activity {
				require.Equal(t, activity[i], respActivity[i], i)
			}
		})
	})
}
