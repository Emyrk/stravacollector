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
