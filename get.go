// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package picago

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const albumURL = "https://picasaweb.google.com/data/feed/api/user/{userID}?start-index={startIndex}"
const photoURL = "https://picasaweb.google.com/data/feed/api/user/{userID}/albumid/{albumID}?start-index={startIndex}"

var DebugDir string

type Album struct {
	ID, Title, Summary, Description, Location string
	AuthorName, AuthorURI                     string
	Keywords                                  []string
	Published, Updated                        time.Time
	URL                                       string
}

type Photo struct {
	ID, ExifUID, Title, Summary, Description, Location string
	Keywords                                           []string
	Published, Updated                                 time.Time
	Latitude, Longitude                                float64
	URL, Type                                          string
}

// GetAlbums returns the list of albums of the given userID.
// If userID is empty, "default" is used.
func GetAlbums(client *http.Client, userID string) ([]Album, error) {
	if userID == "" {
		userID = "default"
	}
	url := strings.Replace(albumURL, "{userID}", userID, 1)

	var albums []Album
	var err error
	hasMore, startIndex := true, 1
	for hasMore {
		albums, hasMore, err = getAlbums(albums, client, url, startIndex)
		if !hasMore {
			break
		}
		startIndex = len(albums) + 1
	}
	return albums, err
}

func getAlbums(albums []Album, client *http.Client, url string, startIndex int) ([]Album, bool, error) {
	if startIndex <= 0 {
		startIndex = 1
	}
	feed, err := downloadAndParse(client,
		strings.Replace(url, "{startIndex}", strconv.Itoa(startIndex), 1))
	if err != nil {
		return nil, false, err
	}
	if len(feed.Entries) == 0 {
		return nil, false, nil
	}
	if cap(albums)-len(albums) < len(feed.Entries) {
		albums = append(albums, make([]Album, 0, len(feed.Entries))...)
	}
	for _, entry := range feed.Entries {
		albumURL := ""
		for _, link := range entry.Links {
			if link.Rel == "http://schemas.google.com/g/2005#feed" {
				albumURL = link.URL
				break
			}
		}
		albums = append(albums, Album{
			ID:          entry.ID,
			Summary:     entry.Summary,
			Title:       entry.Media.Title,
			Description: entry.Media.Description,
			Location:    entry.Location,
			AuthorName:  entry.Author.Name,
			AuthorURI:   entry.Author.URI,
			Keywords:    strings.Split(entry.Media.Keywords, ","),
			Published:   entry.Published,
			Updated:     entry.Updated,
			URL:         albumURL,
		})
	}
	return albums, startIndex+len(feed.Entries) < feed.TotalResults, nil
}

func GetPhotos(client *http.Client, userID, albumID string) ([]Photo, error) {
	if userID == "" {
		userID = "default"
	}
	url := strings.Replace(photoURL, "{userID}", userID, 1)
	url = strings.Replace(url, "{albumID}", albumID, 1)

	var photos []Photo
	var err error
	hasMore, startIndex := true, 1
	for hasMore {
		photos, hasMore, err = getPhotos(photos, client, url, startIndex)
		if !hasMore {
			break
		}
		startIndex = len(photos) + 1
	}
	return photos, err
}

func getPhotos(photos []Photo, client *http.Client, url string, startIndex int) ([]Photo, bool, error) {
	if startIndex <= 0 {
		startIndex = 1
	}
	feed, err := downloadAndParse(client,
		strings.Replace(url, "{startIndex}", strconv.Itoa(startIndex), 1))
	if err != nil {
		return nil, false, err
	}
	if len(feed.Entries) == 0 {
		return nil, false, nil
	}
	if cap(photos)-len(photos) < len(feed.Entries) {
		photos = append(photos, make([]Photo, 0, len(feed.Entries))...)
	}
	for _, entry := range feed.Entries {
		var lat, long float64
		i := strings.Index(entry.Point, " ")
		if i >= 1 {
			lat, err = strconv.ParseFloat(entry.Point[:i], 64)
			if err != nil {
				log.Printf("cannot parse %q as latitude: %v", entry.Point[:i], err)
			}
			long, err = strconv.ParseFloat(entry.Point[i+1:], 64)
			if err != nil {
				log.Printf("cannot parse %q as longitude: %v", entry.Point[i+1:], err)
			}
		}
		if entry.Point != "" && lat == 0 && long == 0 {
			log.Fatalf("point=%q but couldn't parse it as lat/long", entry.Point)
		}
		url, typ := entry.Content.URL, entry.Content.Type
		if url == "" {
			url, typ = entry.Media.Content.URL, entry.Media.Content.Type
		}
		photos = append(photos, Photo{
			ID:          entry.ID,
			ExifUID:     entry.ExifUID,
			Summary:     entry.Summary,
			Title:       entry.Media.Title,
			Description: entry.Media.Description,
			Location:    entry.Location,
			//AuthorName:  entry.Author.Name,
			//AuthorURI:   entry.Author.URI,
			Keywords:  strings.Split(entry.Media.Keywords, ","),
			Published: entry.Published,
			Updated:   entry.Updated,
			URL:       url,
			Type:      typ,
			Latitude:  lat,
			Longitude: long,
		})
	}
	return photos, len(photos) < feed.NumPhotos, nil
}

func downloadAndParse(client *http.Client, url string) (*Atom, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("downloadAndParse=%s %s\n%s", resp.Request.URL, resp.Status, buf)
	}
	var r io.Reader = resp.Body
	if DebugDir != "" {
		fn := filepath.Join(DebugDir, neturl.QueryEscape(url)+".xml")
		xmlfh, err := os.Create(fn)
		if err != nil {
			return nil, fmt.Errorf("error creating debug filx %s: %v", fn, err)
		}
		defer xmlfh.Close()
		r = io.TeeReader(resp.Body, xmlfh)
	}
	return ParseAtom(r)
}

// DownloadPhoto returns an io.ReadCloser for reading the photo bytes
func DownloadPhoto(client *http.Client, url string) (io.ReadCloser, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("downloading %s: %s: %s", url, resp.Status, buf)
	}
	return resp.Body, nil
}
