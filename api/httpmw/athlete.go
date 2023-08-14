package httpmw

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/go-chi/chi/v5"

	"github.com/Emyrk/strava/database"
)

type athleteIDStateKey struct{}

func Athlete(r *http.Request) database.GetAthleteFullRow {
	ath, ok := r.Context().Value(athleteIDStateKey{}).(database.GetAthleteFullRow)
	if !ok {
		panic("developer error: athlete middleware not provided")
	}
	return ath
}

func ExtractAthlete(db database.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			athID, err := strconv.ParseInt(chi.URLParam(r, "athlete_id"), 10, 64)
			if err != nil {
				httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
					Message: "Invalid athlete ID",
					Detail:  err.Error(),
				})
				return
			}

			row, err := db.GetAthleteFull(r.Context(), athID)
			if err != nil {
				httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
					Message: "Failed to fetch athlete",
					Detail:  err.Error(),
				})
				return
			}
			athlete := row.Athlete

			ctx = context.WithValue(ctx, athleteIDStateKey{}, athlete)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
