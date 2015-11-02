package kickass

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/xmlpath.v2"
)

// Kickass default endpoint
const DefaultEndpoint = "https://kat.cr"

// MaxElementsPerPage represents the max number of elements per page
const MaxElementsPerPage = 25

// Custom errors
var (
	ErrUnexpectedContent  = errors.New("kickass: unexpected content")
	ErrMissingUserParam   = errors.New("kickass: missing user param")
	ErrMissingSearchParam = errors.New("kickass: missing search param")
)

// Torrent represents a torrent from kickass
type Torrent struct {
	Name       string
	TorrentURL string
	MagnetURL  string
	Seed       int
	Leech      int
	Age        string
	Size       string
	FileCount  int
	Verified   bool
	User       string
}

// Client represents the kickass client
type Client struct {
	Endpoint   string
	HTTPClient *http.Client
}

// New creates a new client
func New() Client {
	return Client{
		Endpoint:   DefaultEndpoint,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) searchBaseURL(q *Query) string {
	return fmt.Sprintf("%s/usearch/%s", c.Endpoint, q.searchField())
}

// Search searches from a query
func (c *Client) Search(q *Query) ([]*Torrent, error) {
	// The only required param is the search
	if q.Search == "" {
		return nil, ErrMissingSearchParam
	}

	return c.getPages(q, c.searchBaseURL(q))
}

func (c *Client) listByUserBaseURL(q *Query) string {
	return fmt.Sprintf("%s/user/%s/uploads", c.Endpoint, q.User)
}

// ListByUser returns the torrents for a specific user
func (c *Client) ListByUser(q *Query) ([]*Torrent, error) {
	// The only required param is the user
	if q.User == "" {
		return nil, ErrMissingUserParam
	}

	return c.getPages(q, c.listByUserBaseURL(q))
}

// getPages downloads each page and merges the results
func (c *Client) getPages(q *Query, baseURL string) ([]*Torrent, error) {
	torrents := []*Torrent{}

	// Set default number of pages to 1
	if q.Pages == 0 {
		q.Pages = 1
	}

	for i := 1; i <= q.Pages; i++ {
		URL := fmt.Sprintf("%s/%s", baseURL, q.urlParams(i))

		t, err := c.getPage(URL)
		if err != nil {
			return nil, err
		}

		torrents = append(torrents, t...)

		// If the number of results is lower than the max number of elements
		// per page that means there is no need to continue
		if len(t) < MaxElementsPerPage {
			break
		}
	}

	return torrents, nil
}

// getPage downloads a page and parses its content
func (c *Client) getPage(URL string) ([]*Torrent, error) {
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseFunc(root)
}
