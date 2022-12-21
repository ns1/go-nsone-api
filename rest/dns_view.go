package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

// DNSViewService handles 'views/' endpoint.
type DNSViewService service

// List returns all DNS Views
//
// NS1 API docs: https://ns1.com/api#getlist-all-dns-views
func (s *DNSViewService) List() ([]*dns.DNSView, *http.Response, error) {
	return s.ListWithContext(context.Background())
}

// ListWithContext is the same as List, but takes a context.
func (s *DNSViewService) ListWithContext(ctx context.Context) ([]*dns.DNSView, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "views", nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var vl []*dns.DNSView
	resp, err := s.client.Do(req, &vl)
	if err != nil {
		return nil, resp, err
	}

	return vl, resp, nil
}

// Create takes a *dns.DNSView and creates a new DNS View.
//
// The given DNSView must have at least the name
// NS1 API docs: https://ns1.com/api#putcreate-a-dns-view
func (s *DNSViewService) Create(v *dns.DNSView) (*http.Response, error) {
	return s.CreateWithContext(context.Background(), v)
}

// CreateWithContext is the same as Create, but takes a context.
func (s *DNSViewService) CreateWithContext(ctx context.Context, v *dns.DNSView) (*http.Response, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("/v1/views/%s", v.Name), v)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == http.StatusConflict {
				return nil, ErrViewExists
			}
		}

		return resp, err
	}

	return resp, nil
}

// Get takes a DNS view name and returns DNSView struct.
//
// NS1 API docs: https://ns1.com/api#getview-dns-view-details
func (s *DNSViewService) Get(viewName string) (*dns.DNSView, *http.Response, error) {
	return s.GetWithContext(context.Background(), viewName)
}

// GetWithContext is the same as Get, but takes a context.
func (s *DNSViewService) GetWithContext(ctx context.Context, viewName string) (*dns.DNSView, *http.Response, error) {
	path := fmt.Sprintf("views/%s", viewName)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	var v dns.DNSView
	resp, err := s.client.Do(req, &v)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == http.StatusNotFound {
				return nil, resp, ErrViewMissing
			}
		}
		return nil, resp, err
	}

	return &v, resp, nil
}

// Update takes a *dns.DNSView and updates the DNS view with same name on NS1.
//
// NS1 API docs: https://ns1.com/api#postedit-a-dns-view
func (s *DNSViewService) Update(v *dns.DNSView) (*http.Response, error) {
	return s.UpdateWithContext(context.Background(), v)
}

// UpdateWithContext is the same as Update, but takes a context.
func (s *DNSViewService) UpdateWithContext(ctx context.Context, v *dns.DNSView) (*http.Response, error) {
	path := fmt.Sprintf("views/%s", v.Name)

	req, err := s.client.NewRequest("POST", path, &v)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req, &v)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == http.StatusNotFound {
				return resp, ErrViewMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// Delete takes a DNS view name, and removes an existing DNS view
//
// NS1 API docs: https://ns1.com/api#deletedelete-a-dns-view
func (s *DNSViewService) Delete(viewName string) (*http.Response, error) {
	return s.DeleteWithContext(context.Background(), viewName)
}

// DeleteWithContext is the same as Delete, but takes a context.
func (s *DNSViewService) DeleteWithContext(ctx context.Context, viewName string) (*http.Response, error) {
	path := fmt.Sprintf("views/%s", viewName)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req, nil)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == http.StatusNotFound {
				return resp, ErrViewMissing
			}
		}
		return resp, err
	}

	return resp, nil
}

// GetPreference returns a map[string]int of preferences.
//
// NS1 API docs: https://ns1.com/api#getget-dns-view-preference
func (s *DNSViewService) GetPreferences() (map[string]int, *http.Response, error) {
	return s.GetPreferencesWithContext(context.Background())
}

// GetPreferencesWithContext is the same as GetPreferences, but takes a context.
func (s *DNSViewService) GetPreferencesWithContext(ctx context.Context) (map[string]int, *http.Response, error) {
	path := "config/views/preference"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	m := make(map[string]int)
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// UpdatePreference takes a map[string]int and returns a map[string]int of preferences.
//
// NS1 API docs: https://ns1.com/api#postedit-dns-view-preference
func (s *DNSViewService) UpdatePreferences(m map[string]int) (map[string]int, *http.Response, error) {
	return s.UpdatePreferencesWithContext(context.Background(), m)
}

// UpdatePreferencesWithContext is the same as UpdatePreferences, but takes a context.
func (s *DNSViewService) UpdatePreferencesWithContext(ctx context.Context, m map[string]int) (map[string]int, *http.Response, error) {
	path := "config/views/preference"

	req, err := s.client.NewRequest("POST", path, m)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	mapUpdated := make(map[string]int)
	resp, err := s.client.Do(req, &mapUpdated)
	if err != nil {
		switch errType := err.(type) {
		case *Error:
			if errType.Resp.StatusCode == http.StatusNotFound {
				return nil, resp, ErrViewMissing
			}
		}
		return nil, resp, err
	}

	return mapUpdated, resp, nil
}

var (
	// ErrViewExists bundles CREATE error.
	ErrViewExists = errors.New("DNS view already exists")

	// ErrViewExists bundles GET error.
	ErrViewMissing = errors.New("DNS view not found")
)
