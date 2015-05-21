package kickass

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/xmlpath.v2"
)

var (
	defaultAddress      = "http://kat.cr"
	xpathTorrentResults = xmlpath.MustCompile("//tr[contains(@id, 'torrent_')]")
	xpathTorrentName    = xmlpath.MustCompile(".//a[@class=\"cellMainLink\"]")
	xpathTorrentURL     = xmlpath.MustCompile(".//a[contains(@title,'Download torrent file')]/@href")
	xpathMagnetURL      = xmlpath.MustCompile(".//a[contains(@class, 'imagnet')]/@href")
	xpathSeed           = xmlpath.MustCompile(".//td[5]")
	xpathLeech          = xmlpath.MustCompile(".//td[6]")
	xpathAge            = xmlpath.MustCompile(".//td[4]")
	xpathSize           = xmlpath.MustCompile(".//td[2]")
	xpathFileCount      = xmlpath.MustCompile(".//td[3]")
	xpathVerify         = xmlpath.MustCompile(".//a[contains(@class,'iverify')]")
	xpathUser           = xmlpath.MustCompile(".//a[contains(@href, '/user/')]")
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

// Client kickass client
type Client struct {
	Address    string
	HTTPClient *http.Client
}

// New create client
func New() Client {
	return Client{
		Address:    defaultAddress,
		HTTPClient: &http.Client{},
	}
}

// Search torrent
func (c *Client) Search(query string) (*[]Torrent, error) {
	query = url.QueryEscape(query)
	URL := c.Address + "/usearch/" + query

	resp, err := c.HTTPClient.Get(URL)
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

// SearchByUser search torrent by user
func (c *Client) SearchByUser(query, user string) (*[]Torrent, error) {
	q := fmt.Sprintf("%s user:%s", query, user)
	return c.Search(q)
}

func parseResult(root *xmlpath.Node) (*[]Torrent, error) {
	torrents := []Torrent{}
	iter := xpathTorrentResults.Iter(root)
	for iter.Next() {
		name, ok := xpathTorrentName.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		torrentURL, ok := xpathTorrentURL.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		magnet, ok := xpathMagnetURL.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		verify := xpathVerify.Exists(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}

		seedStr, ok := xpathSeed.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		seed, err := strconv.Atoi(seedStr)
		if err != nil {
			return nil, err
		}

		leechStr, ok := xpathLeech.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		leech, err := strconv.Atoi(leechStr)
		if err != nil {
			return nil, err
		}

		age, ok := xpathAge.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		fileCountStr, ok := xpathFileCount.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		fileCount, err := strconv.Atoi(fileCountStr)
		if err != nil {
			return nil, err
		}

		size, ok := xpathSize.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		user, ok := xpathUser.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}

		t := Torrent{
			Name:       name,
			TorrentURL: torrentURL,
			MagnetURL:  magnet,
			Seed:       seed,
			Leech:      leech,
			Age:        age,
			FileCount:  fileCount,
			Size:       size,
			Verified:   verify,
			User:       user,
		}

		torrents = append(torrents, t)
	}

	return &torrents, nil
}
