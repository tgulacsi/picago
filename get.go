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
	"strconv"
	"strings"
	"time"
)

const albumURL = "https://picasaweb.google.com/data/feed/api/user/%(userID)s"
const photoURL = "https://picasaweb.google.com/data/feed/api/user/%(userID)s/albumid/%(albumID)s"

type Album struct {
	ID, Title, Summary, Description, Location string
	AuthorName, AuthorURI                     string
	Keywords                                  []string
	Published, Updated                        time.Time
	URL                                       string
}

type Photo struct {
	ID, ExifUID, Title, Summary, Description string
	Keywords                                 []string
	Published, Updated                       time.Time
	Latitude, Longitude                      float64
	URL                                      string
}

// GetAlbums returns the list of albums of the given userID.
// If userID is empty, "default" is used.
func GetAlbums(client *http.Client, userID string) ([]Album, error) {
	if userID == "" {
		userID = "default"
	}
	resp, err := client.Get(strings.Replace(albumURL, "%(userID)s", userID, 1))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GetAlbums(%s)=%s %s\n%s", userID, resp.Request.URL, resp.Status, buf)
	}
	feed, err := ParseAtom(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(feed.Entries) == 0 {
		return nil, nil
	}
	albums := make([]Album, 0, len(feed.Entries))
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
	return albums, nil
}

func GetPhotos(client *http.Client, userID, albumID string) ([]Photo, error) {
	if userID == "" {
		userID = "default"
	}
	url := strings.Replace(photoURL, "%(userID)s", userID, 1)
	resp, err := client.Get(strings.Replace(url, "%(albumID)s", albumID, 1))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		buf, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GetPhotos(%s, %s)=%s %s\n%s", userID, albumID, resp.Request.URL, resp.Status, buf)
	}
	feed, err := ParseAtom(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(feed.Entries) == 0 {
		return nil, nil
	}
	photos := make([]Photo, 0, len(feed.Entries))
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
		photos = append(photos, Photo{
			ID:          entry.ID,
			ExifUID:     entry.ExifUID,
			Summary:     entry.Summary,
			Title:       entry.Media.Title,
			Description: entry.Media.Description,
			//AuthorName:  entry.Author.Name,
			//AuthorURI:   entry.Author.URI,
			Keywords:  strings.Split(entry.Media.Keywords, ","),
			Published: entry.Published,
			Updated:   entry.Updated,
			URL:       entry.Media.Content.URL,
			Latitude:  lat,
			Longitude: long,
		})
	}
	return photos, nil
}

// DownloadPhoto returns an io.ReadCloser for reading the photo bytes
func DownloadPhoto(client *http.Client, url string) (io.ReadCloser, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
