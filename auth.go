// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package picago

import (
    "net/http"
	"errors"

    "code.google.com/p/goauth2/oauth"
)

const oAuthScope = "https://picasaweb.google.com/data/"
const oAuthEndpoint = " https://accounts.google.com/o/oauth2/auth"

// Authorize authorizes using OAuth2
// the ID and secret strings can be acquired from Google for the application
// https://developers.google.com/accounts/docs/OAuth2#basicsteps
func Authorize(ID, secret string) error {
	return errors.New("Not implemented")
}


type PicasaClient struct {
    *http.Client
}

// For redirect_uri, see https://developers.google.com/accounts/docs/OAuth2InstalledApp#choosingredirecturi .
//
func NewClient(id, secret string) (*PicasaClient, error) {
    client := &http.Client{}
    url := oAuthEndpount + "/?response_type=code&include_granted_scopes=true&" +
    "client_id="+id+"&redirect_uri=http://localhost:"+strconv.Itoa(port)

    config := &oauth.Config{
        ClientId: id,
        ClientSecret: secret,
        RedirectURL: "http://localhost:" + strconv.Itoa(port)
}

func AuthorizationHandle(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, fmt.Sprintf("cannot parse as form: %v", err), http.StatusBadRequest)
        return
    }
    if r.F1G
}
