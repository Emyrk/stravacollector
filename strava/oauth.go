package strava

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
)

func Oauth(ctx context.Context) {
	conf := oauth2.Config{
		ClientID:     "105702",
		ClientSecret: "5eff30500222d8764453038426cdd18edb43c3c2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://www.strava.com/oauth/authorize",
			TokenURL:  "https://www.strava.com/oauth/authorize",
			AuthStyle: 0,
		},
		RedirectURL: "http://localhost:8000",
		Scopes:      []string{},
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	client.Get("...")
	c := New("")
	c.Client = client
	resp, err := c.AthleteSegmentEfforts(ctx, 653262, 2)
	fmt.Println(resp, err)
}
