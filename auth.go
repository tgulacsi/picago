// Copyright 2014 Tamás Gulácsi. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package picago

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const picasaScope = "https://picasaweb.google.com/data/"

var ErrCodeNeeded = errors.New("Authorization code is needed")

// Authorize authorizes using OAuth2
// the ID and secret strings can be acquired from Google for the application
// https://developers.google.com/accounts/docs/OAuth2#basicsteps
func Authorize(ID, secret string) error {
	return errors.New("Not implemented")
}

// NewClient returns an authorized http.Client usable for requests,
// caching tokens in the given file.
func NewClient(id, secret, code, tokenCacheFilename string) (*http.Client, error) {
	return NewClientCache(id, secret, code, oauth2.CacheFile(tokenCacheFilename))
}

// For redirect_uri, see https://developers.google.com/accounts/docs/OAuth2InstalledApp#choosingredirecturi .
//
// NewClientCache returns an authorized http.Client with the given oauth2.Cache implementation
func NewClientCache(id, secret, code string, cache oauth2.Cache) (*http.Client, error) {
	config, err := NewConfig(id, secret, cache)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(context.Background(), config.TokenSource())
	/*
				l, err := getListener()
				if err != nil {
					return nil, err
				}
				donech := make(chan struct{}, 1)
				transport.Config.RedirectURL = "http://" + l.Addr().String()
				// Get an authorization code from the data provider.
				// ("Please ask the user if I can access this resource.")
				url := transport.Config.AuthCodeURL("picago")
				fmt.Println("Visit this URL to allow access to your Picasa data:\n")
				fmt.Println(url)

				srv := &http.Server{Handler: NewAuthorizeHandler(transport, donech)}
				go srv.Serve(l)
				<-donech
				l.Close()

				if transport.Token == nil {
					return nil, ErrCodeNeeded
				}
			}
			if transport.Token == nil {
				// Exchange the authorization code for an access token.
				// ("Here's the code you gave the user, now give me a token!")
				transport.Token, err = transport.Exchange(code)
				if err != nil {
					return nil, fmt.Errorf("Exchange: %v", err)
				}
			}
		}
		return &http.Client{Transport: transport}, nil
	*/
}

func NewConfig(id, secret string, cache oauth2.Cache) (*oauth2.Config, error) {
	if id == "" || secret == "" {
		return nil, errors.New("Client ID and secret is needed!")
	}
	return &oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientId:     id,
		ClientSecret: secret,
		Scopes:       []string{picasaScope},
		TokenCache:   cache,
	}, nil
}

// NewAuthorizeHandler returns a http.HandlerFunc which will set the Token of
// the given oauth2.Transport and send a struct{} on the donech on success.
func NewAuthorizeHandler(config *oauth2.Config, donech chan<- struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := config.Exchange(r.FormValue("code"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error exchanging code: %v", err), http.StatusBadRequest)
			return
		}
		transport.Token = token
		// The Transport now has a valid Token.
		donech <- struct{}{}
	}
}

func getListener() (*net.TCPListener, error) {
	return net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
}
