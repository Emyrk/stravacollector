package httpmw

import (
	"context"
	"net/http"

	"github.com/Emyrk/strava/api/modelsdk"

	"github.com/Emyrk/strava/api/httpapi"

	"github.com/Emyrk/strava/api/auth"
)

const (
	StravaAuthJWTCookie = "strava_auth_jwt"
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

func AuthenticatedAthleteIDOptional(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(authAthIDStateKey{}).(int64)
	if !ok {
		return -1, false
	}
	return id, true
}

func Authenticated(a *auth.Authentication, optional bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			payload := r.Header.Get(StravaAuthJWTHeader)
			if payload == "" {
				cookie, err := r.Cookie(StravaAuthJWTCookie)
				if err != nil {
					if optional {
						next.ServeHTTP(rw, r.WithContext(ctx))
						return
					}
					httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
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
				if optional {
					next.ServeHTTP(rw, r.WithContext(ctx))
					return
				}
				httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
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
