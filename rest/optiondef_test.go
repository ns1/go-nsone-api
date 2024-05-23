package rest_test

import (
	"net/http"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dhcp"

	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func TestDHCPOptionDef(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	if err != nil {
		t.Fatalf("Error creating mock service: %v", err)
	}
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		defer mock.ClearTestCases()

		client.FollowPagination = true
		sgs := []dhcp.OptionDef{
			{
				Code:         1,
				Description:  "a",
				FriendlyName: "a",
				Schema: dhcp.OptionDefSchema{
					Type: dhcp.SchemaTypeString,
				},
			},
			{
				Code:         2,
				Description:  "b",
				FriendlyName: "b",
				Schema: dhcp.OptionDefSchema{
					Type: dhcp.SchemaTypeString,
				},
			},
		}
		err := mock.AddTestCase(http.MethodGet, "/dhcp/optiondef", http.StatusOK, nil, nil, "", sgs)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respSgs, _, err := client.OptionDef.List()
		if err != nil {
			t.Fatalf("error listing DHCP option definitions: %v", err)
		}
		if len(respSgs) != len(sgs) {
			t.Errorf("wrong length: want=%d, got=%d", len(sgs), len(respSgs))
		}

		for i, sg := range respSgs {
			if sg.FriendlyName != sgs[i].FriendlyName || sg.Code != sgs[i].Code ||
				sg.Description != sgs[i].Description || sg.Schema.Type != sgs[i].Schema.Type {
				t.Errorf("Incorrect data for option definition %d: want=%+v, got=%+v", i, sgs[i], sg)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		defer mock.ClearTestCases()

		client.FollowPagination = false
		od := dhcp.OptionDef{
			Code:         1,
			Description:  "a",
			FriendlyName: "a",
			Schema: dhcp.OptionDefSchema{
				Type: dhcp.SchemaTypeString,
			},
		}

		err := mock.AddTestCase(http.MethodGet, "/dhcp/optiondef/space/key", http.StatusOK, nil, nil, "", od)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respOD, _, err := client.OptionDef.Get("space", "key")
		if err != nil {
			t.Fatalf("error getting scpe group: %v", err)
		}
		if od.FriendlyName != respOD.FriendlyName || od.Code != respOD.Code ||
			od.Description != respOD.Description || od.Schema.Type != respOD.Schema.Type {
			t.Errorf("wrong option def returned, want=%+v, got=%+v", od, respOD)
		}
	})

	t.Run("Create", func(t *testing.T) {
		defer mock.ClearTestCases()

		t.Run("RequiredParams", func(t *testing.T) {
			od := &dhcp.OptionDef{}
			_, _, err = client.OptionDef.Create(od, "space", "key")
			if err == nil {
				t.Errorf("expected a missing code to result in an error")
			}

			od = &dhcp.OptionDef{
				Code: 1,
			}
			_, _, err = client.OptionDef.Create(od, "space", "key")
			if err == nil {
				t.Errorf("expected a missing friendly name to result in an error")
			}

			od = &dhcp.OptionDef{
				Code:        1,
				Description: "a",
			}
			_, _, err = client.OptionDef.Create(od, "space", "key")
			if err == nil {
				t.Errorf("expected a missing description to result in an error")
			}

			od = &dhcp.OptionDef{
				Code:         1,
				Description:  "a",
				FriendlyName: "a",
			}
			_, _, err = client.OptionDef.Create(od, "space", "key")
			if err == nil {
				t.Errorf("expected a missing schema type to result in an error")
			}
		})

		od := &dhcp.OptionDef{
			Code:         1,
			Description:  "a",
			FriendlyName: "a",
			Schema: dhcp.OptionDefSchema{
				Type: dhcp.SchemaTypeString,
			},
		}
		err := mock.AddTestCase(http.MethodPut, "/dhcp/optiondef/space/key", http.StatusCreated, nil, nil, od, od)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		respOD, _, err := client.OptionDef.Create(od, "space", "key")
		if err != nil {
			t.Fatalf("error creating option definition: %v", err)
		}
		if od.FriendlyName != respOD.FriendlyName || od.Code != respOD.Code ||
			od.Description != respOD.Description || od.Schema.Type != respOD.Schema.Type {
			t.Errorf("wrong option def returned, want=%+v, got=%+v", od, respOD)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		defer mock.ClearTestCases()

		err := mock.AddTestCase(http.MethodDelete, "/dhcp/optiondef/space/key", http.StatusNoContent, nil, nil, "", nil)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		_, err = client.OptionDef.Delete("space", "key")
		if err != nil {
			t.Fatalf("error deleting option definition: %v", err)
		}
	})
}
