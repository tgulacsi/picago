// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// pica-dl implements a simple Picasa Web downloader.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tgulacsi/picago"
)

func main() {
	flag.Parse()
	userid := flag.Arg(0)

	client := &http.Client{}
	albums, err := picago.GetAlbums(client, userid)
	if err != nil {
		log.Fatalf("error listing albums: %v", albums)
	}

	for _, album := range albums {
		log.Printf("downloading album %s.", album)
		photos, err := picago.GetPhotos(client, userid, album.ID)
		if err != nil {
			log.Printf("error listing photos of %s: %v", album.ID, err)
			continue
		}
		log.Printf("album %s contains %d photos.", album.ID, len(photos))
		for _, photo := range photos {
			log.Printf("Photo: %s", photo)
		}
	}
}
