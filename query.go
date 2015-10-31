package kickass

import (
	"net/url"
	"strings"
)

// Query represents a kickass query
type Query struct {
	Search   string
	User     string
	OrderBy  string
	Order    string
	Category string
	ImdbID   string
}

// urlParams extracts the relevant params and builds a query string
func (q *Query) urlParams() string {
	urlValues := &url.Values{}

	if q.OrderBy != "" {
		urlValues.Add("field", q.OrderBy)
	}

	if q.Order != "" {
		urlValues.Add("sorder", q.Order)
	}

	str := urlValues.Encode()
	if str != "" {
		str = "?" + str
	}

	return str
}

// searchField return a string with the relevant query params
func (q *Query) searchField() string {
	search := q.Search

	if q.User != "" {
		search += " user:" + q.User
	}

	if q.Category != "" {
		search += " category:" + q.Category
	}

	if q.ImdbID != "" {
		// Remove the "tt" part from the ID
		search += " imdb:" + strings.Replace(q.ImdbID, "tt", "", -1)
	}

	return search
}