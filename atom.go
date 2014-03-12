// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package picago

import (
	"encoding/xml"
	"io"
	"time"
)

type Atom struct {
	ID           string    `xml:"id"`
	Updated      time.Time `xml:"updated"`
	Title        string    `xml:"title"`
	Subtitle     string    `xml:"subtitle"`
	Icon         string    `xml:"icon"`
	Author       Author    `xml:"author"`
	TotalResults int       `xml:"totalResults"`
	StartIndex   int       `xml:"startIndex"`
	ItemsPerPage int       `xml:"itemsPerPage"`
	Entries      []Entry   `xml:"entry"`
}

type Entry struct {
	ETag      string    `xml:"etag,attr"`
	EntryID   string    `xml:"http://www.w3.org/2005/Atom id"`
	ID        string    `xml:"http://schemas.google.com/photos/2007 id"`
	Published time.Time `xml:"published"`
	Updated   time.Time `xml:"updated"`
	Title     string    `xml:"title"`
	Summary   string    `xml:"summary"`
	Links     []Link    `xml:"link"`
	Author    Author    `xml:"author"`
	Location  string    `xml:"http://schemas.google.com/photos/2007 location"`
	NumPhotos int       `xml:"numphotos"`
	Media     Media     `xml:"group"`
	ExifUID   string    `xml:"tags>imageUniqueID"`
	Point     string    `xml:"where>Point>pos"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	URL  string `xml:"href,attr"`
}

type Media struct {
	Title       string `xml:"http://search.yahoo.com/mrss title"`
	Description string `xml:"description"`
	Keywords    string `xml:"keywords"`
	Content     `xml:"content"`
}

type Content struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

func ParseAtom(r io.Reader) (*Atom, error) {
	result := new(Atom)
	if err := xml.NewDecoder(r).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}
