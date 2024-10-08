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

func TestCreateUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var u account.User
		require.NoError(t, json.Unmarshal(b, &u))
		assert.Nil(t, u.Permissions.Security)
		assert.False(t, u.Permissions.Monitoring.ManageJobs)
		assert.False(t, u.Permissions.Monitoring.CreateJobs)
		assert.False(t, u.Permissions.Monitoring.UpdateJobs)
		assert.False(t, u.Permissions.Monitoring.DeleteJobs)

		w.Write(b)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	u := &account.User{
		Name:        "name-1",
		Username:    "user-1",
		Email:       "email-1",
		Permissions: account.PermissionsMap{},
	}

	_, err := c.Users.Create(u)
	require.NoError(t, err)
}
