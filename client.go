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
	URL := fmt.Sprintf("%s/usearch/%s/%s", c.Endpoint, q.searchField(), q.urlParams())
	return c.getPage(URL)
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

// ListByUser returns the torrents for a specific user
func (c *Client) ListByUser(q *Query) ([]*Torrent, error) {
	URL := fmt.Sprintf("%s/user/%s/uploads/%s", c.Endpoint, q.User, q.urlParams())
	return c.getPage(URL)
}
