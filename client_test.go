package kickass

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gopkg.in/xmlpath.v2"
)

var mockDefaultClient = New()

func TestMissingParams(t *testing.T) {
	for expectedErr, f := range map[error]func(*Query) ([]*Torrent, error){
		ErrMissingSearchParam: mockDefaultClient.Search,
		ErrMissingUserParam:   mockDefaultClient.ListByUser,
	} {
		_, err := f(&Query{})
		if err != expectedErr {
			t.Errorf("expected error %q, got %q", expectedErr, err)
		}
	}
}

func TestGetResults(t *testing.T) {
	parseFunc = func(root *xmlpath.Node) ([]*Torrent, error) {
		return nil, nil
	}

	// Request URL
	var reqURL string

	// Fake server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `test`)
		reqURL = fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
	}))
	defer ts.Close()

	// Client
	c := New()
	c.Endpoint = ts.URL

	type mock struct {
		f           func(q *Query) ([]*Torrent, error)
		baseURLFunc func(q *Query) string
	}

	for _, m := range []mock{
		{
			f:           c.ListByUser,
			baseURLFunc: c.listByUserBaseURL,
		},
		{
			f:           c.Search,
			baseURLFunc: c.searchBaseURL,
		},
	} {
		_, err := m.f(mockQuery)
		if err != nil {
			t.Fatalf("expected no error got %q", err)
		}

		url, err := url.Parse(fmt.Sprintf("%s/%d/%s", m.baseURLFunc(mockQuery), 1, mockQuery.urlParams()))
		if err != nil {
			t.Fatalf("failed to parse URL: %q", err)
		}

		if reqURL != url.String() {
			t.Errorf("invalid search url: expected %q, got %q", url.String(), reqURL)
		}
	}
}
