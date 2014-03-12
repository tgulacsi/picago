// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

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
