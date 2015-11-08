package kickass

import (
	"strconv"

	"gopkg.in/xmlpath.v2"
)

// XPATH
var (
	xpathNoResult       = xmlpath.MustCompile("//text()[contains(.,'did not match any documents')]")
	xpathTorrentResults = xmlpath.MustCompile("//tr[contains(@id, 'torrent_')]")
	xpathTorrentName    = xmlpath.MustCompile(".//a[@class=\"cellMainLink\"]")
	xpathTorrentURL     = xmlpath.MustCompile(".//a[contains(@title,'Download torrent file')]/@href")
	xpathMagnetURL      = xmlpath.MustCompile(".//a[contains(@title,'Torrent magnet link')]/@href")
	xpathSeed           = xmlpath.MustCompile(".//td[5]")
	xpathLeech          = xmlpath.MustCompile(".//td[6]")
	xpathAge            = xmlpath.MustCompile(".//td[4]")
	xpathSize           = xmlpath.MustCompile(".//td[2]")
	xpathFileCount      = xmlpath.MustCompile(".//td[3]")
	xpathVerify         = xmlpath.MustCompile(".//a[contains(@class,'iverify')]")
	xpathUser           = xmlpath.MustCompile(".//a[contains(@href, '/user/')]")
)

// Default parse function, to be overwritten during the tests
var parseFunc = parseResult

func parseResult(root *xmlpath.Node) ([]*Torrent, error) {
	torrents := []*Torrent{}

	// Don't go further if there is no results
	if xpathNoResult.Exists(root) {
		return torrents, nil
	}

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
			// The user name is not always present
			user = ""
		}

		t := &Torrent{
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

	return torrents, nil
}
