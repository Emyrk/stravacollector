package httpmw

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/lib/slice"

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

func AuthenticatedAsAdmins() func(next http.Handler) http.Handler {
	return AuthenticatedAs(
		2661162,  // Steven
		20563755, // Bia
	)
}

func RequestAuthenticatedAsAdminsOrMe(rw http.ResponseWriter, r *http.Request, cmp int64) bool {
	var (
		id = AuthenticatedAthleteID(r)
	)

	if id == cmp {
		return true
	}

	if slice.Contains([]int64{2661162, 20563755}, id) {
		return true
	}

	httpapi.Write(r.Context(), rw, http.StatusUnauthorized, modelsdk.Response{
		Message: "Not authorized",
		Detail:  fmt.Sprintf("id %d is not allowed to access this resource", id),
	})
	return false
}

func AuthenticatedAs(allowedIDs ...int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				ctx = r.Context()
				id  = AuthenticatedAthleteID(r)
			)

			if !slice.Contains(allowedIDs, id) {
				httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
					Message: "Not authorized",
					Detail:  fmt.Sprintf("id %d is not allowed to access this resource", id),
				})
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
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
