package river

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestShouldRefetchZeroSegments(t *testing.T) {
	for _, tc := range []struct {
		name     string
		count    int32
		distance float64
		segments string
		want     bool
	}{
		{"first completed download", 1, 10000, "[]", true},
		{"not downloaded", 0, 10000, "[]", false},
		{"refetch completed", 2, 10000, "[]", false},
		{"later manual download", 3, 10000, "[]", false},
		{"short activity", 1, 1000, "[]", false},
		{"segments present", 1, 10000, "[{}]", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var activity strava.DetailedActivity
			require.NoError(t, json.Unmarshal([]byte(`{"segment_efforts":`+tc.segments+`}`), &activity))
			activity.Distance = tc.distance
			require.Equal(t, tc.want, shouldRefetchZeroSegments(activity, tc.count))
		})
	}
}

func TestRetryableActivityFetchError(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		want   bool
	}{
		{500, "internal error", true}, {501, "not implemented", true},
		{502, "bad gateway", true}, {503, "unavailable", true}, {504, "timeout", true},
		{599, "server error", true}, {600, "outside server error range", false},
		{400, "bad request", false}, {401, "unauthorized", false},
		{403, "forbidden", false}, {404, "not found", false},
		{429, "rate limited", false}, {597, "maintenance", false},
		{502, "Strava is temporarily unavailable", false},
	} {
		t.Run(fmt.Sprintf("%d/%s", tc.status, tc.body), func(t *testing.T) {
			err := &strava.StravaAPIError{Response: &http.Response{StatusCode: tc.status}, Body: []byte(tc.body)}
			require.Equal(t, tc.want, retryableActivityFetchError(err))
			require.Equal(t, tc.want, retryableActivityFetchError(fmt.Errorf("fetch: %w", err)))
		})
	}
	require.False(t, retryableActivityFetchError(nil))
	require.False(t, retryableActivityFetchError(errors.New("network error")))
}

type fetchActivityErrorStore struct{ database.Store }

func (fetchActivityErrorStore) GetAthleteLogin(context.Context, int64) (database.AthleteLogin, error) {
	// A zero expiry is a non-expiring test token; no OAuth token refresh needed.
	return database.AthleteLogin{OauthAccessToken: "test-token"}, nil
}
func (fetchActivityErrorStore) GetActivityDetail(context.Context, int64) (database.ActivityDetail, error) {
	return database.ActivityDetail{}, sql.ErrNoRows
}

type fetchActivityTransport func(*http.Request) (*http.Response, error)

func (f fetchActivityTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchActivityWorkRetriesServerErrors(t *testing.T) {
	for _, status := range []int{500, 502, 503, 504} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			client := &http.Client{Transport: fetchActivityTransport(func(req *http.Request) (*http.Response, error) {
				calls++
				return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"message":"temporary failure"}`)), Request: req}, nil
			})}
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
			worker := FetchActivityWorker{mgr: &Manager{db: fetchActivityErrorStore{}, oauthCfg: &oauth2.Config{}, logger: zerolog.Nop()}}
			err := worker.Work(ctx, &river.Job[FetchActivityArgs]{
				JobRow: &rivertype.JobRow{ID: 1, Priority: 1},
				Args:   FetchActivityArgs{Source: database.ActivityDetailSourceManual, ActivityID: 1, AthleteID: 1},
			})
			require.Error(t, err)
			apiErr := strava.IsAPIError(err)
			require.NotNil(t, apiErr)
			require.Equal(t, status, apiErr.Response.StatusCode)
			require.Equal(t, 1, calls)
		})
	}
}
