// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// pica-dl implements a simple Picasa Web downloader.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/tgulacsi/picago"
)

// See https://developers.google.com/accounts/docs/OAuth2InstalledApp .
func main() {
	flagID := flag.String("id", os.Getenv("CLIENT_ID"), "application client ID")
	flagSecret := flag.String("secret", os.Getenv("CLIENT_SECRET"), "application client secret")
	flagCode := flag.String("code", os.Getenv("AUTH_CODE"), "authorization code")
	flagTokenCache := flag.String("cache", "token-cache.json", "token cache filename")

	flag.Parse()
	userid := flag.Arg(0)

	pica, err := picago.NewClient(*flagID, *flagSecret, *flagCode, *flagTokenCache)
	if err != nil {
		log.Fatalf("error with authorization: %v", err)
	}
	albums, err := picago.GetAlbums(pica, userid)
	if err != nil {
		log.Fatalf("error listing albums: %v", err)
	}

	for _, album := range albums {
		albumJ, err := json.Marshal(album)
		if err != nil {
			log.Fatalf("error marshaling %#v: %v", album, err)
		}
		log.Printf("downloading album %s.", albumJ)
		photos, err := picago.GetPhotos(pica, userid, album.ID)
		if err != nil {
			log.Printf("error listing photos of %s: %v", album.ID, err)
			continue
		}
		log.Printf("album %s contains %d photos.", album.ID, len(photos))
		for _, photo := range photos {
			photoJ, err := json.Marshal(photo)
			if err != nil {
				log.Fatalf("error marshaling %#v: %v", photo, err)
			}
			log.Printf("Photo: %s", photoJ)
		}
	}
}
