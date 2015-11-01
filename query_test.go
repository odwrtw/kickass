package kickass

import "testing"

var mockQuery = &Query{
	Search:   "star wars",
	User:     "YIFY",
	OrderBy:  "seeders",
	Order:    "desc",
	Category: "movies",
	ImdbID:   "tt2488496",
	Pages:    2,
}

func TestQueryURLParams(t *testing.T) {
	// Test with mock query
	expectedURLParams := "?field=seeders&sorder=desc"
	if mockQuery.urlParams() != expectedURLParams {
		t.Errorf("invalid URL params, expected %q, got %q", expectedURLParams, mockQuery.urlParams())
	}

	// Test with empty query
	q := &Query{}
	if q.urlParams() != "" {
		t.Errorf("invalid URL params, expected nothing, got %q", q.urlParams())
	}
}

func TestQuerySearchField(t *testing.T) {
	// Test with mock query
	expectedSearch := "star wars user:YIFY category:movies imdb:2488496"
	if mockQuery.searchField() != expectedSearch {
		t.Errorf("invalid search field, expected %q, got %q", expectedSearch, mockQuery.searchField())
	}

	// Test with simple search
	q := &Query{Search: mockQuery.Search}
	if q.searchField() != q.Search {
		t.Errorf("invalid search field, expected %q, got %q", mockQuery.Search, q.urlParams())
	}
}
