package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/api/auth"
	"github.com/Emyrk/strava/api/auth/authkeys"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/gencache"
)

type refreshStore struct {
	database.Store
	login                         database.AthleteLogin
	loadErr, tokenErr, profileErr error
	savedToken                    *database.UpsertAthleteLoginParams
	savedProfile                  *database.UpsertAthleteParams
}

func (s *refreshStore) GetAthleteLogin(_ context.Context, id int64) (database.AthleteLogin, error) {
	if id != 67733138 {
		return database.AthleteLogin{}, sql.ErrNoRows
	}
	return s.login, s.loadErr
}
func (s *refreshStore) UpsertAthleteLogin(_ context.Context, p database.UpsertAthleteLoginParams) (database.AthleteLogin, error) {
	s.savedToken = &p
	return s.login, s.tokenErr
}
func (s *refreshStore) UpsertAthlete(_ context.Context, p database.UpsertAthleteParams) (database.Athlete, error) {
	s.savedProfile = &p
	return database.Athlete{}, s.profileErr
}

type refreshTransport func(*http.Request) (*http.Response, error)

func (f refreshTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRefreshAthlete(t *testing.T) {
	key, err := authkeys.GenerateKey()
	require.NoError(t, err)
	authentication, err := auth.New(auth.Options{SecretPEM: authkeys.MarshalPrivateKey(key), Issuer: "refresh-test"})
	require.NoError(t, err)
	for _, tc := range []struct {
		name                          string
		caller                        int64
		id                            string
		expired                       bool
		loadErr, tokenErr, profileErr error
		upstreamStatus                int
		profileStatus                 int
		mismatch                      bool
		want                          int
		wantProfile, wantToken        bool
	}{
		{name: "admin", caller: 2661162, want: 200, wantProfile: true},
		{name: "rotated token", caller: 2661162, expired: true, want: 200, wantProfile: true, wantToken: true},
		{name: "rotated token saved despite profile failure", caller: 2661162, expired: true, profileStatus: 503, want: 502, wantToken: true},
		{name: "unauthenticated", want: 401},
		{name: "other athlete", caller: 123, want: 401},
		{name: "other admin excluded", caller: 20563755, want: 401},
		{name: "invalid id", caller: 2661162, id: "bad", want: 400},
		{name: "negative id", caller: 2661162, id: "-1", want: 400},
		{name: "zero id", caller: 2661162, id: "0", want: 400},
		{name: "no authorization", caller: 2661162, loadErr: sql.ErrNoRows, want: 404},
		{name: "database unavailable", caller: 2661162, loadErr: errors.New("secret"), want: 500},
		{name: "revoked refresh token", caller: 2661162, expired: true, upstreamStatus: 401, want: 502},
		{name: "token persistence failure", caller: 2661162, expired: true, tokenErr: errors.New("secret"), want: 500, wantToken: true},
		{name: "strava unavailable", caller: 2661162, upstreamStatus: 503, want: 502},
		{name: "strava rate limit", caller: 2661162, upstreamStatus: 429, want: 429},
		{name: "different athlete", caller: 2661162, mismatch: true, want: 502},
		{name: "profile save failure", caller: 2661162, profileErr: errors.New("secret"), want: 500, wantProfile: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			store := &refreshStore{login: database.AthleteLogin{
				AthleteID: 67733138, ProviderID: "client", Summit: true,
				OauthAccessToken: "old-secret", OauthRefreshToken: "refresh-secret", OauthTokenType: "Bearer",
				OauthExpiry: database.Timestamptz(time.Now().Add(24 * time.Hour)),
			}, loadErr: tc.loadErr, tokenErr: tc.tokenErr, profileErr: tc.profileErr}
			if tc.expired {
				store.login.OauthExpiry = database.Timestamptz(time.Now().Add(-24 * time.Hour))
			}
			calls := 0
			client := &http.Client{Transport: refreshTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				status, body := 200, ""
				if r.URL.Path == "/token" {
					require.Equal(t, "POST", r.Method)
					require.NoError(t, r.ParseForm())
					require.Equal(t, "refresh-secret", r.Form.Get("refresh_token"))
					body = `{"access_token":"new-secret","refresh_token":"rotated-secret","token_type":"Bearer","expires_in":21600}`
				} else {
					require.Equal(t, "/api/v3/athlete", r.URL.Path)
					expected := "Bearer old-secret"
					if tc.expired {
						expected = "Bearer new-secret"
						require.NotNil(t, store.savedToken)
					}
					require.Equal(t, expected, r.Header.Get("Authorization"))
					body = `{"id":67733138,"username":"toddn_mcnaughton","firstname":"Todd","lastname":"McNaughton","sex":"M"}`
					if tc.mismatch {
						body = `{"id":99}`
					}
				}
				if tc.profileStatus != 0 && r.URL.Path == "/api/v3/athlete" {
					status = tc.profileStatus
					body = `{"message":"secret"}`
				}
				if tc.upstreamStatus != 0 {
					status = tc.upstreamStatus
					body = `{"message":"secret"}`
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
			fetches := 0
			cache := gencache.New(ctx, 15*time.Minute, func(context.Context) ([]database.HugelLeaderboardRow, error) { fetches++; return nil, nil })
			_, err := cache.Load(ctx)
			require.NoError(t, err)
			api := &API{Opts: &Options{DB: store}, Auth: authentication, Registry: prometheus.NewRegistry(),
				OAuthConfig:         &oauth2.Config{ClientID: "client", ClientSecret: "client-secret", Endpoint: oauth2.Endpoint{TokenURL: "https://strava.test/token", AuthStyle: oauth2.AuthStyleInParams}},
				HugelBoard2026Cache: cache,
			}
			id := tc.id
			if id == "" {
				id = "67733138"
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/refresh-athlete/"+id, nil).WithContext(ctx)
			if tc.caller != 0 {
				session, err := authentication.CreateSession(ctx, tc.caller)
				require.NoError(t, err)
				if tc.name == "admin" {
					req.AddCookie(&http.Cookie{Name: httpmw.StravaAuthJWTCookie, Value: session})
				} else {
					req.Header.Set(httpmw.StravaAuthJWTHeader, session)
				}
			}
			rec := httptest.NewRecorder()
			router := api.Routes()
			router.ServeHTTP(rec, req)
			require.Equal(t, tc.want, rec.Code, rec.Body.String())
			require.Equal(t, tc.wantProfile, store.savedProfile != nil)
			require.Equal(t, tc.wantToken, store.savedToken != nil)
			require.NotContains(t, rec.Body.String(), "secret")
			if tc.want == 200 {
				require.Equal(t, int64(67733138), store.savedProfile.ID)
				require.Equal(t, "Todd", store.savedProfile.Firstname)
				require.Equal(t, "M", store.savedProfile.Sex)
				require.JSONEq(t, `{}`, string(store.savedProfile.Clubs))
				var body map[string]any
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
				require.Equal(t, "Todd", body["firstname"])
			}
			if tc.wantToken {
				require.Equal(t, "rotated-secret", store.savedToken.OauthRefreshToken)
				require.Equal(t, "client", store.savedToken.ProviderID)
				require.True(t, store.savedToken.Summit)
			}
			_, err = cache.Load(ctx)
			require.NoError(t, err)
			expectedFetches := 1
			if tc.want == 200 {
				expectedFetches = 2
			}
			require.Equal(t, expectedFetches, fetches)
			if tc.want == 200 {
				req.Method = http.MethodGet
				methodRec := httptest.NewRecorder()
				router.ServeHTTP(methodRec, req)
				require.Equal(t, http.StatusMethodNotAllowed, methodRec.Code)
			}
			if tc.want == 401 || tc.want == 400 || tc.loadErr != nil {
				require.Zero(t, calls)
			}
		})
	}
}
