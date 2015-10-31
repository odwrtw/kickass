Kickass
=======

[![GoDoc](https://godoc.org/github.com/odwrtw/kickass?status.svg)](http://godoc.org/github.com/odwrtw/kickass)

Golang library to get torrents from kickass

## Features

* Search with paramaters
* List the uploads of a user

## Example

```go
    package main

    import (
            "github.com/kr/pretty"
            "github.com/odwrtw/kickass"
    )

    func main() {
            // New kickass client
            k := kickass.New()

            // Search query
            query := &kickass.Query{
                  Search:   "star wars",
                  User:     "YIFY",
                  OrderBy:  "seeders",
                  Order:    "desc",
                  Category: "movies",
            }

            // Search
            torrents, err := k.Search(query)
            if err != nil {
                    pretty.Println(err)
                    return
            }

            pretty.Println(torrents)
    }
```
