package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func TestVersion(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()
	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		t.Run("list versions", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			versions := []*dns.Version{
				&dns.Version{
					Id:          15,
					Name:        "version1",
					Active:      false,
					ActivatedAt: 2,
					CreatedAt:   3,
				},
				&dns.Version{
					Id:          16,
					Name:        "version2",
					Active:      true,
					ActivatedAt: 3,
					CreatedAt:   4,
				},
			}
			require.Nil(t, mock.AddVersionListTestCase("versioned.zone", nil, nil, versions))

			respVersions, _, err := client.Versions.List("versioned.zone")
			require.Nil(t, err)
			require.NotNil(t, respVersions)
			require.Equal(t, len(versions), len(respVersions))

			for i := range versions {
				require.Equal(t, versions[i].Name, respVersions[i].Name, i)
				require.Equal(t, versions[i].Id, respVersions[i].Id, i)
				require.Equal(t, versions[i].Active, respVersions[i].Active, i)
				require.Equal(t, versions[i].ActivatedAt, respVersions[i].ActivatedAt, i)
				require.Equal(t, versions[i].CreatedAt, respVersions[i].CreatedAt, i)

			}
		})

		t.Run("error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/zones/versioned.zone/versions", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				versions, resp, err := client.Versions.List("versioned.zone")
				require.Nil(t, versions)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("create versions", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			version := dns.Version{
				Id:          15,
				Name:        "version1",
				Active:      false,
				ActivatedAt: 2,
				CreatedAt:   3,
			}
			require.Nil(t, mock.AddCreateVersionTestCase("versioned.zone", nil, nil, &version))

			boolFlag := false
			respVersion, _, err := client.Versions.Create("versioned.zone", boolFlag)
			require.Nil(t, err)
			require.NotNil(t, respVersion)

			require.Equal(t, version.Name, respVersion.Name)
			require.Equal(t, version.Id, respVersion.Id)
			require.Equal(t, version.Active, respVersion.Active)
			require.Equal(t, version.ActivatedAt, respVersion.ActivatedAt)
			require.Equal(t, version.CreatedAt, respVersion.CreatedAt)

		})

		t.Run("error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, "/zones/versioned.zone/versions?force=false", http.StatusInternalServerError,
					nil, nil, "", `{"message": "test error"}`,
				))

				boolFlag := false
				version, resp, err := client.Versions.Create("versioned.zone", boolFlag)
				require.Nil(t, version)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("delete versions", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddDeleteVersionTestCase("versioned.zone", 15, nil, nil))
			_, err := client.Versions.Delete("versioned.zone", 15)
			require.Nil(t, err)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodDelete, "/zones/versioned.zone/versions/15", http.StatusInternalServerError,
					nil, nil, "", `{"message": "test error"}`,
				))
				_, err := client.Versions.Delete("versioned.zone", 15)
				require.NotNil(t, err)
			})
		})
	})

	t.Run("Activate", func(t *testing.T) {
		t.Run("delete versions", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddActivateVersionTestCase("versioned.zone", 15, nil, nil))
			_, err := client.Versions.Activate("versioned.zone", 15)
			require.Nil(t, err)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, "/zones/versioned.zone/versions/15", http.StatusInternalServerError,
					nil, nil, "", `{"message": "test error"}`,
				))
				_, err := client.Versions.Activate("versioned.zone", 15)
				require.NotNil(t, err)
			})
		})
	})

}
