package kickass

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/xmlpath.v2"
)

// Kickass endpoint
const Endpoint = "https://kat.cr"

var (
	//ErrUnexpectedContent returned when addic7ed's website seem to have change
	ErrUnexpectedContent = errors.New("Unexpected content")
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
	Endpoint string
}

// New create client
func New() Client {
	return Client{
		Endpoint: Endpoint,
	}
}

// Search searches from a query
func (c *Client) Search(q *Query) ([]*Torrent, error) {
	baseURL := fmt.Sprintf("%s/usearch/%s", c.Endpoint, q.searchField())
	return c.getPages(q, baseURL)
}

// ListByUser returns the torrents for a specific user
func (c *Client) ListByUser(q *Query) ([]*Torrent, error) {
	baseURL := fmt.Sprintf("%s/user/%s/uploads", c.Endpoint, q.User)
	return c.getPages(q, baseURL)
}

// getPages downloads each page and merges the results
func (c *Client) getPages(q *Query, baseURL string) ([]*Torrent, error) {
	torrents := []*Torrent{}

	// Set default number of pages to 1
	if q.Pages == 0 {
		q.Pages = 1
	}

	for i := 1; i <= q.Pages; i++ {
		URL := fmt.Sprintf("%s/%d/%s", baseURL, i, q.urlParams())
		t, err := c.getPage(URL)
		if err != nil {
			return nil, err
		}

		torrents = append(torrents, t...)
	}

	return torrents, nil
}

// getPage downloads a page and parses its content
func (c *Client) getPage(URL string) ([]*Torrent, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseResult(root)
}
