package account

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalUsers(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		out  User
	}{
		{
			"basic",
			[]byte(`{
    "permissions": {},
    "teams": [],
    "email": "support@nsone.net",
    "last_access": 1376325771.0,
    "created": 1376325771.0,
    "notify": {
      "billing": true
    },
    "name": "API Example",
    "username": "apiexample",
	"ip_whitelist": ["1.1.1.1", "2.2.2.2"],
	"ip_whitelist_strict": true,
	"2fa_enabled": true,
	"shared_auth": {
        "saml": {
            "sso": true,
            "idp": {
                "use_metadata_url": true,
                "metadata_url": null,
                "metadata_file": null,
                "provider": "okta"
            }
        }
    }
  }`),
			User{
				Username:             "apiexample",
				Name:                 "API Example",
				Email:                "support@nsone.net",
				LastAccess:           1376325771.0,
				Created:              1376325771.0,
				Notify:               NotificationSettings{true},
				TeamIDs:              []string{},
				IPWhitelist:          []string{"1.1.1.1", "2.2.2.2"},
				IPWhitelistStrict:    true,
				TwoFactorAuthEnabled: true,
				SharedAuth: SharedAuth{
					SAML: SAML{
						SSO: true,
						IDP: IDP{
							UseMetadataURL: pointerBool(true),
							MetadataURL:    nil,
							MetadataFile:   nil,
							Provider:       pointerString("okta"),
						},
					},
				},
			},
		},
		{
			"perms",
			[]byte(`{
    "permissions": {
      "dns": {
        "view_zones": true,
        "manage_zones": true,
        "zones_allow_by_default": false,
        "zones_deny": [],
        "zones_allow": ["example.com"]
      },
      "data": {
        "push_to_datafeeds": false,
        "manage_datasources": false,
        "manage_datafeeds": false
      },
      "account": {
        "manage_payment_methods": false,
        "manage_plan": false,
        "manage_teams": false,
        "manage_apikeys": false,
        "manage_account_settings": false,
        "view_activity_log": false,
        "view_invoices": false,
        "manage_users": false,
        "manage_ip_whitelist": true
      },
      "monitoring": {
        "manage_lists": false,
        "manage_jobs": false,
        "view_jobs": false
      }
    },
    "teams": ["520422919f782d37dffb588a"],
    "email": "newuser@example.com",
    "last_access": null,
    "notify": {
      "billing": true
    },
    "name": "New User",
    "username": "newuser"
  }`),
			User{
				Username:   "newuser",
				Name:       "New User",
				Email:      "newuser@example.com",
				TeamIDs:    []string{"520422919f782d37dffb588a"},
				LastAccess: 0.0,
				Notify:     NotificationSettings{true},
				Permissions: PermissionsMap{
					DNS: PermissionsDNS{
						ViewZones:           true,
						ManageZones:         true,
						ZonesAllowByDefault: false,
						ZonesDeny:           []string{},
						ZonesAllow:          []string{"example.com"},
					},
					Data: PermissionsData{
						PushToDatafeeds:   false,
						ManageDatasources: false,
						ManageDatafeeds:   false,
					},
					Account: PermissionsAccount{
						ManagePaymentMethods:  false,
						ManagePlan:            false,
						ManageTeams:           false,
						ManageApikeys:         false,
						ManageAccountSettings: false,
						ViewActivityLog:       false,
						ViewInvoices:          false,
						ManageUsers:           false,
						ManageIPWhitelist:     true,
					},
					Monitoring: PermissionsMonitoring{
						ManageLists: false,
						ManageJobs:  false,
						ViewJobs:    false,
						CreateJobs:  false,
						UpdateJobs:  false,
						DeleteJobs:  false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u User
			require.NoError(t, json.Unmarshal(tt.in, &u))
			require.Equal(t, tt.out, u)
		})
	}
}

func pointerBool(b bool) *bool {
	return &b
}

func pointerString(s string) *string {
	return &s
}
