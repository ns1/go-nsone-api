package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dataset"
)

func TestDatasetsService(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	id := "87461586-c681-43b7-b283-f840be95e13c"

	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			dtList := []*dataset.Dataset{
				{
					ID:   "dt-1",
					Name: "My dataset 1",
					Datatype: &dataset.Datatype{
						Type:  dataset.DatatypeTypeNumQueries,
						Scope: dataset.DatatypeScopeAccount,
						Data:  nil,
					},
				},
				{
					ID:   "dt-2",
					Name: "My dataset 2",
					Datatype: &dataset.Datatype{
						Type:  dataset.DatatypeTypeNumQueries,
						Scope: dataset.DatatypeScopeAccount,
						Data:  nil,
					},
				},
			}

			require.Nil(t, mock.AddDatasetListTestCase(nil, nil, dtList))

			respDts, _, err := client.Datasets.List()
			require.Nil(t, err)
			require.NotNil(t, respDts)
			require.Equal(t, len(dtList), len(respDts))

			for i := range dtList {
				require.Equal(t, dtList[i].ID, respDts[i].ID, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/datasets", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				_, resp, err := client.Datasets.List()
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				_, resp, err := c.Datasets.List()
				require.Nil(t, resp)
				require.Error(t, err)
			})
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			dt := &dataset.Dataset{
				Name: "My dataset",
				Datatype: &dataset.Datatype{
					Type:  dataset.DatatypeTypeNumQueries,
					Scope: dataset.DatatypeScopeAccount,
					Data:  nil,
				},
				Repeat: nil,
				Timeframe: &dataset.Timeframe{
					Aggregation: dataset.TimeframeAggregationMontly,
					Cycles:      func() *int32 { i := int32(1); return &i }(),
				},
				ExportType:      dataset.ExportTypeCSV,
				RecipientEmails: nil,
			}
			require.Nil(t, mock.AddDatasetGetTestCase(id, nil, nil, dt))

			respDt, _, err := client.Datasets.Get(id)
			require.Nil(t, err)
			require.NotNil(t, respDt)
			require.Equal(t, dt.ID, respDt.ID)
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/datasets/"+id, http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				dt, resp, err := client.Datasets.Get(id)
				require.Nil(t, dt)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				dt, resp, err := c.Datasets.Get(id)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, dt)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		dt := &dataset.Dataset{
			Name: "My dataset",
			Datatype: &dataset.Datatype{
				Type:  dataset.DatatypeTypeNumQueries,
				Scope: dataset.DatatypeScopeAccount,
				Data:  nil,
			},
			Repeat: nil,
			Timeframe: &dataset.Timeframe{
				Aggregation: dataset.TimeframeAggregationMontly,
				Cycles:      func() *int32 { i := int32(1); return &i }(),
			},
			ExportType:      dataset.ExportTypeCSV,
			RecipientEmails: nil,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddDatasetCreateTestCase(nil, nil, dt, dt))

			_, _, err := client.Datasets.Create(dt)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/datasets", http.StatusConflict,
				nil, nil, dt, `{"message": "invalid parameters"}`,
			))

			_, _, err := client.Datasets.Create(dt)
			require.Contains(t, err.Error(), "invalid parameters")
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddDatasetDeleteTestCase(id, nil, nil))

			_, err := client.Datasets.Delete(id)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodDelete, "/datasets/"+id, http.StatusNotFound,
				nil, nil, "", `{"message": "dataset not found"}`,
			))

			_, err := client.Datasets.Delete(id)
			require.Equal(t, api.ErrDatasetNotFound, err)
		})
	})

	t.Run("Get Report", func(t *testing.T) {
		reportId := "f840be95e13c"

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()
			fileContents := []byte(`foo,bar`)
			require.Nil(t, mock.AddDatasetGetReportTestCase(id, reportId, nil, nil, fileContents))

			buf, _, err := client.Datasets.GetReport(id, reportId)
			require.Nil(t, err)
			assert.Equal(t, fileContents, buf.Bytes())
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodGet, "/datasets/"+id+"/reports/"+reportId, http.StatusNotFound,
				nil, nil, "", `{"message": "dataset not found"}`,
			))

			_, _, err := client.Datasets.GetReport(id, reportId)
			require.Equal(t, api.ErrDatasetNotFound, err)
		})
	})
}
