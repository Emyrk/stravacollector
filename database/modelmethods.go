package database

import "golang.org/x/oauth2"

func (a *AthleteLogin) OAuthToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  a.OauthAccessToken,
		TokenType:    a.OauthTokenType,
		RefreshToken: a.OauthRefreshToken,
		Expiry:       a.OauthExpiry,
	}
}

func DistanceToMiles(distance float64) float64 {
	return distance / 1609.34
}

func DistanceToFeet(distance float64) float64 {
	return distance / 3.28084
}
