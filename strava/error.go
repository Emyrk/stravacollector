package strava

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// {"message":"Rate Limit Exceeded","errors":[{"resource":"Application","field":"overall rate limit","code":"exceeded"}]}

type StravaAPIError struct {
	Response *http.Response
	Body     []byte

	// Authentication errors from oauth2 package
	oauthError *oauth2.RetrieveError
}

func (e StravaAPIError) Error() string {
	if e.oauthError != nil {
		return fmt.Sprintf("strava oauth error: %s", e.oauthError.Error())
	}

	if strings.Contains(string(e.Body), "<!DOCTYPE html>") {
		return fmt.Sprintf("status code: %d\nbody: %s", e.Response.StatusCode, "'HTML response'")
	}
	return fmt.Sprintf("status code: %d\nbody: %s", e.Response.StatusCode, string(e.Body))
}

func (e StravaAPIError) IsRefreshTokenError() bool {
	if e.oauthError == nil {
		return false
	}

	if e.oauthError.Response.StatusCode != http.StatusBadRequest {
		return false
	}

	body := e.oauthError.Body
	if !strings.Contains(string(body), `"field":"refresh_token"`) || !strings.Contains(string(body), `"code":"invalid"}`) {
		return false
	}

	return true
}

func IsAPIError(err error) *StravaAPIError {
	var e *StravaAPIError
	if errors.As(err, &e) {
		return e
	}

	var oauthError *oauth2.RetrieveError
	if errors.As(err, &oauthError) {
		return &StravaAPIError{
			Response:   oauthError.Response,
			Body:       oauthError.Body,
			oauthError: oauthError,
		}
	}
	return nil
}

func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	var e Error
	if errors.As(err, &e) {
		if e.Message == "Rate Limit Exceeded" {
			return true
		}
	}
	return false
}

type Error struct {
	Message string          `json:"message"`
	Errors  []DetailedError `json:"errors"`
}

func (e Error) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf(":%s-%s-%s", err.Code, err.Resource, err.Field))
	}
	return b.String()
}

func (e Error) ContainsCode(code string) bool {
	for _, err := range e.Errors {
		if err.Code == code {
			return true
		}
	}
	return false
}

type DetailedError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}
