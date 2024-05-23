package rest_test

import (
	"net/http"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dhcp"

	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func TestDHCPReservation(t *testing.T) {
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
		sgs := []dhcp.Reservation{
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
			{IDAddress: &idAddr},
		}
		err := mock.AddTestCase(http.MethodGet, "/dhcp/reservation", http.StatusOK, nil, nil, "", sgs)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respSgs, _, err := client.Reservation.List()
		if err != nil {
			t.Fatalf("error listing reservations: %v", err)
		}
		if len(respSgs) != len(sgs) {
			t.Errorf("wrong length: want=%d, got=%d", len(sgs), len(respSgs))
		}

		for i, sg := range respSgs {
			if *sg.IDAddress != *sgs[i].IDAddress {
				t.Errorf("Incorrect name for reservation %d: want=%q, got=%q", i, *sgs[i].IDAddress, *sg.IDAddress)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		defer mock.ClearTestCases()

		idAddr := 1
		sg := dhcp.Reservation{IDAddress: &idAddr}

		err := mock.AddTestCase(http.MethodGet, "/dhcp/reservation/1", http.StatusOK, nil, nil, "", sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respAddr, _, err := client.Reservation.Get(1)
		if err != nil {
			t.Fatalf("error getting reservation: %v", err)
		}
		if *respAddr.IDAddress != *sg.IDAddress {
			t.Errorf("wrong reservation returned, want=%+v, got=%+v", sg, respAddr)
		}
	})

	t.Run("Create", func(t *testing.T) {
		defer mock.ClearTestCases()

		t.Run("RequiredParams", func(t *testing.T) {
			sg := &dhcp.Reservation{}
			_, _, err = client.Reservation.Create(sg)
			if err == nil {
				t.Errorf("expected a missing address id to result in an error")
			}
		})

		idAddr := 123
		sg := &dhcp.Reservation{
			IDAddress: &idAddr,
			Options:   make([]dhcp.Option, 0),
		}
		err := mock.AddTestCase(http.MethodPut, "/dhcp/reservation", http.StatusCreated, nil, nil, sg, sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		respSG, _, err := client.Reservation.Create(sg)
		if err != nil {
			t.Fatalf("error creating reservation: %v", err)
		}
		if *respSG.IDAddress != *sg.IDAddress {
			t.Errorf("wrong reservation returned: want=%+v, got=%+v", sg, respSG)
		}
	})

	t.Run("Edit", func(t *testing.T) {
		t.Run("RequiredParams", func(t *testing.T) {
			idAddr := 123
			sg := &dhcp.Reservation{
				IDAddress: &idAddr,
				Options:   make([]dhcp.Option, 0),
			}
			_, _, err = client.Reservation.Edit(sg)
			if err == nil {
				t.Errorf("expected a missing ID to result in an error")
			}

			sg = &dhcp.Reservation{}
			_, _, err = client.Reservation.Edit(sg)
			if err == nil {
				t.Errorf("expected a missing reservation ID to result in an error")
			}
		})

		defer mock.ClearTestCases()

		id := 1
		idAddr := 123
		sg := &dhcp.Reservation{
			ID:        &id,
			IDAddress: &idAddr,
			Options:   make([]dhcp.Option, 0),
		}
		err := mock.AddTestCase(http.MethodPost, "/dhcp/reservation/1", http.StatusOK, nil, nil, sg, sg)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respSG, _, err := client.Reservation.Edit(sg)
		if err != nil {
			t.Fatalf("error editing reservation: %v", err)
		}
		if respSG.IDAddress != sg.IDAddress {
			t.Errorf("wrong reservation returned: want=%+v, got=%+v", sg, respSG)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		defer mock.ClearTestCases()

		err := mock.AddTestCase(http.MethodDelete, "/dhcp/reservation/1", http.StatusNoContent, nil, nil, "", nil)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		_, err = client.Reservation.Delete(1)
		if err != nil {
			t.Fatalf("error deleting reservation: %v", err)
		}
	})
}
