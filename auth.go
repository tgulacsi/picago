// Package picago implements a picasaweb walker
//
// https://developers.google.com/picasa-web/docs/2.0/developers_guide_protocol?hl=en
package picago

import (
	"errors"
)

const oAuthScope = "https://picasaweb.google.com/data/"

// Authorize authorizes using OAuth2
// the ID and secret strings can be acquired from Google for the application
// https://developers.google.com/accounts/docs/OAuth2#basicsteps
func Authorize(ID, secret string) error {
	return errors.New("Not implemented")
}
