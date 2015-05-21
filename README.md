Kickass
=======

A golang library for kickass torrent source


    package main

    import (
            "github.com/kr/pretty"
            "gitlab.quimbo.fr/nicolas/kickass"
    )

    func main() {
            k := kickass.New()
            torrents, err := k.SearchByUser("game of throne", "ettv")

            if err != nil {
                    pretty.Println(err)
            }

            pretty.Println(torrents)

    }
