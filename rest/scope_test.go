package rest_test

import (
	"net/http"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dhcp"

	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func TestDHCPScope(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	if err != nil {
		t.Fatalf("Error creating mock service: %v", err)
	}
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		defer mock.ClearTestCases()

		client.FollowPagination = true
		idAddr := 1
		sgs := []dhcp.Scope{
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
		}
		err := mock.AddTestCase(http.MethodGet, "/dhcp/scope", http.StatusOK, nil, nil, "", sgs)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respSgs, _, err := client.Scope.List()
		if err != nil {
			t.Fatalf("error listing scopes: %v", err)
		}
		if len(respSgs) != len(sgs) {
			t.Errorf("wrong length: want=%d, got=%d", len(sgs), len(respSgs))
		}

		for i, sg := range respSgs {
			if *sg.IDAddress != *sgs[i].IDAddress {
				t.Errorf("Incorrect name for scope %d: want=%q, got=%q", i, *sgs[i].IDAddress, *sg.IDAddress)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		defer mock.ClearTestCases()

		idAddr := 1
		sg := dhcp.Scope{IDAddress: &idAddr}

		err := mock.AddTestCase(http.MethodGet, "/dhcp/scope/1", http.StatusOK, nil, nil, "", sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respAddr, _, err := client.Scope.Get(1)
		if err != nil {
			t.Fatalf("error getting scope: %v", err)
		}
		if *respAddr.IDAddress != *sg.IDAddress {
			t.Errorf("wrong scope returned, want=%+v, got=%+v", sg, respAddr)
		}
	})

	t.Run("Create", func(t *testing.T) {
		defer mock.ClearTestCases()

		t.Run("RequiredParams", func(t *testing.T) {
			sg := &dhcp.Scope{}
			_, _, err = client.Scope.Create(sg)
			if err == nil {
				t.Errorf("expected a missing address id to result in an error")
			}
		})

		idAddr := 123
		sg := &dhcp.Scope{
			IDAddress: &idAddr,
		}
		err := mock.AddTestCase(http.MethodPut, "/dhcp/scope", http.StatusCreated, nil, nil, sg, sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		respSG, _, err := client.Scope.Create(sg)
		if err != nil {
			t.Fatalf("error creating scope: %v", err)
		}
		if *respSG.IDAddress != *sg.IDAddress {
			t.Errorf("wrong scope returned: want=%+v, got=%+v", sg, respSG)
		}
	})

	t.Run("Edit", func(t *testing.T) {
		t.Run("RequiredParams", func(t *testing.T) {
			idAddr := 123
			sg := &dhcp.Scope{IDAddress: &idAddr}
			_, _, err = client.Scope.Edit(sg)
			if err == nil {
				t.Errorf("expected a missing ID to result in an error")
			}

			sg = &dhcp.Scope{}
			_, _, err = client.Scope.Edit(sg)
			if err == nil {
				t.Errorf("expected a missing address ID to result in an error")
			}
		})

		defer mock.ClearTestCases()

		idAddr := 123
		sg := &dhcp.Scope{
			ID:        1,
			IDAddress: &idAddr,
		}
		err := mock.AddTestCase(http.MethodPost, "/dhcp/scope/1", http.StatusOK, nil, nil, sg, sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respSG, _, err := client.Scope.Edit(sg)
		if err != nil {
			t.Fatalf("error editing scope: %v", err)
		}
		if respSG.IDAddress != sg.IDAddress {
			t.Errorf("wrong address returned: want=%+v, got=%+v", sg, respSG)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		defer mock.ClearTestCases()

		err := mock.AddTestCase(http.MethodDelete, "/dhcp/scope/1", http.StatusNoContent, nil, nil, "", nil)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		_, err = client.Scope.Delete(1)
		if err != nil {
			t.Fatalf("error deleting scope: %v", err)
		}
	})
}
