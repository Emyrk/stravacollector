package httpmw

import (
	"context"
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"

	"github.com/Emyrk/strava/api/auth"
)

const (
	StravaAuthJWTCookie = "strava-auth-jwt"
	StravaAuthJWTHeader = "Strava-Auth-JWT"
)

type authAthIDStateKey struct{}

func AuthenticatedAthleteID(r *http.Request) int64 {
	id, ok := r.Context().Value(authAthIDStateKey{}).(int64)
	if !ok {
		panic("developer error: authenticated middleware not provided")
	}
	return id
}

func Authenticated(a *auth.Authentication) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			payload := r.Header.Get(StravaAuthJWTHeader)
			if payload == "" {
				cookie, err := r.Cookie(StravaAuthJWTCookie)
				if err != nil {
					httpapi.Write(ctx, rw, http.StatusUnauthorized, httpapi.Response{
						Message: "No valid authentication provided",
						Detail:  err.Error(),
					})
					return
				}
				payload = cookie.Value
			}

			athleteID, err := a.ValidateSession(payload)
			if err != nil {
				// Delete expired cookies
				http.SetCookie(rw, &http.Cookie{
					Name:   StravaAuthJWTCookie,
					MaxAge: -1,
				})
				httpapi.Write(ctx, rw, http.StatusUnauthorized, httpapi.Response{
					Message: "Authentication failed",
					Detail:  err.Error(),
				})
				return
			}

			ctx = context.WithValue(ctx, authAthIDStateKey{}, athleteID)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
