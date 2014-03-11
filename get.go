// Package picago implements a picasaweb walker
//
// https://developers.google.com/picasa-web/docs/2.0/developers_guide_protocol?hl=en
package picago

import (
	"net/http"
	"strings"
	"time"
)

const albumURL = "https://picasaweb.google.com/data/feed/api/user/"

// GetAlbums returns the list of albums of the given userID.
// If userID is empty, "default" is used.
func GetAlbums(client *http.Client, userID string) ([]Album, error) {
	resp, err := client.Get(albumURL + "/" + userID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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

type Album struct {
	ID, Title, Summary, Description, Location string
	AuthorName, AuthorURI                     string
	Keywords                                  []string
	Published, Updated                        time.Time
	URL                                       string
}
