package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func TestCreateTeam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var tm account.Team
		require.NoError(t, json.Unmarshal(b, &tm))
		assert.Nil(t, tm.Permissions.Security)
		assert.False(t, tm.Permissions.Monitoring.ManageJobs)
		assert.False(t, tm.Permissions.Monitoring.CreateJobs)
		assert.False(t, tm.Permissions.Monitoring.UpdateJobs)
		assert.False(t, tm.Permissions.Monitoring.DeleteJobs)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	tm := &account.Team{
		ID:          "id-1",
		Name:        "team-1",
		Permissions: account.PermissionsMap{},
	}

	_, err := c.Teams.Create(tm)
	require.NoError(t, err)
}
