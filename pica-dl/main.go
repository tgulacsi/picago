// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

// pica-dl implements a simple Picasa Web downloader.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tgulacsi/picago"
)

// See https://developers.google.com/accounts/docs/OAuth2InstalledApp .
func main() {
	flagID := flag.String("id", os.Getenv("CLIENT_ID"), "application client ID")
	flagSecret := flag.String("secret", os.Getenv("CLIENT_SECRET"), "application client secret")
	flagCode := flag.String("code", os.Getenv("AUTH_CODE"), "authorization code")
	flagTokenCache := flag.String("cache", "token-cache.json", "token cache filename")
	flagDir := flag.String("dir", "", "directory to download images to")
	flagDebugDir := flag.String("debug", "", "set to a valid path to save the response XMLs there")

	flagBBox := flag.String("bbox", "", "bounding box for the query (west,south,east,north coordinates)")
	flagLocation := flag.String("location", "", "location for the query (e.g. London)")
	flagText := flag.String("text", "", "text to search for in tags, comments and title")
	flagTag := flag.String("tag", "", "tag to search for")

	flag.Parse()
	qp := picago.QueryParams{Location: *flagLocation, Text: *flagText, Tag: *flagTag}
	if *flagBBox != "" {
		if _, err := fmt.Sscanf(*flagBBox, "%f,%f,%f,%f",
			&qp.BBox.West, &qp.BBox.South, &qp.BBox.East, &qp.BBox.North,
		); err != nil {
			log.Fatalf("error parsing bbox=%q as four floats: %v", *flagBBox, err)
		}
	}
	picago.DebugDir = *flagDebugDir
	userid := flag.Arg(0)

	client, err := picago.NewClient(*flagID, *flagSecret, *flagCode, *flagTokenCache)
	if err != nil {
		log.Fatalf("error with authorization: %v", err)
	}
	user, err := picago.GetUser(client, "")
	log.Printf("user=%#v err=%v", user, err)

	albums, err := picago.GetAlbums(client, userid)
	if err != nil {
		log.Fatalf("error listing albums: %v", err)
	}
	log.Printf("user %s has %d albums.", userid, len(albums))

	download := *flagDir != ""
	dir, fn := "", ""
	for _, album := range albums {
		albumJ, err := json.Marshal(album)
		if err != nil {
			log.Fatalf("error marshaling %#v: %v", album, err)
		}
		if download {
			dir = filepath.Join(*flagDir, album.Name)
			if err = os.MkdirAll(dir, 0750); err != nil {
				log.Fatalf("cannot create directory %s: %v", dir, err)
			}
			fn = filepath.Join(dir, "album-"+album.Name+".json")
			if err = ioutil.WriteFile(fn, albumJ, 0750); err != nil {
				log.Fatalf("error writing %s: %v", fn, err)
			}
		}
		log.Printf("downloading album %s.", albumJ)
		photos, err := picago.GetPhotosSpec(client, userid, album.ID, qp)
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

			if !download {
				continue
			}
			fn = filepath.Join(dir, photo.Filename)
			if err = ioutil.WriteFile(fn+".json", photoJ, 0750); err != nil {
				log.Fatalf("error writing %s.json: %v", fn, err)
			}
			if err = downloadTo(fn, client, photo.URL); err != nil {
				log.Fatalf("downloading %s: %v", photo.URL, err)
			}
		}
	}
}

func downloadTo(fn string, client *http.Client, url string) error {
	body, err := picago.DownloadPhoto(client, url)
	if err != nil {
		return err
	}
	defer body.Close()
	fh, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fh, body)
	return err
}
