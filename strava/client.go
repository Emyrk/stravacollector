package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Emyrk/strava/strava/stravalimit"
)

type Client struct {
	AccessToken string
	Client      *http.Client
}

func NewOAuthClient(cli *http.Client) *Client {
	return &Client{
		Client: cli,
	}
}

func New(accessToken string) *Client {
	return &Client{
		AccessToken: accessToken,
		Client:      http.DefaultClient,
	}
}

type GetActivitiesParams struct {
	Before  time.Time
	After   time.Time
	Page    int
	PerPage int
}

func (c *Client) GetActivities(ctx context.Context, params GetActivitiesParams) ([]ActivitySummary, error) {
	vals := url.Values{}
	if !params.Before.IsZero() && params.Before.Unix() > 0 {
		vals.Set("before", fmt.Sprintf("%d", params.Before.Unix()))
	}
	if !params.After.IsZero() && params.After.Unix() > 0 {
		vals.Set("after", fmt.Sprintf("%d", params.After.Unix()))
	}
	if params.Page > 0 {
		vals.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PerPage > 0 {
		vals.Set("per_page", fmt.Sprintf("%d", params.PerPage))
	}
	resp, err := c.Request(ctx, http.MethodGet, "/athlete/activities", nil, vals)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	var activities []ActivitySummary
	return activities, c.DecodeResponse(resp, &activities, http.StatusOK)
}

func (c *Client) GetActivity(ctx context.Context, activityID int64, includeEfforts bool) (DetailedActivity, error) {
	resp, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("/activities/%d", activityID), nil, url.Values{
		"include_all_efforts": []string{strconv.FormatBool(includeEfforts)},
	})
	if err != nil {
		return DetailedActivity{}, fmt.Errorf("request: %w", err)
	}

	var activity DetailedActivity
	return activity, c.DecodeResponse(resp, &activity, http.StatusOK)
}

func (c *Client) GetAuthenticatedAthelete(ctx context.Context) (Athlete, error) {
	resp, err := c.Request(ctx, http.MethodGet, "/athlete", nil, nil)
	if err != nil {
		return Athlete{}, fmt.Errorf("request: %w", err)
	}

	var athlete Athlete
	return athlete, c.DecodeResponse(resp, &athlete, http.StatusOK)
}

func (c *Client) AthleteSegmentEfforts(ctx context.Context, segmentID int, perPage int) ([]DetailedSegmentEffort, error) {
	var efforts []DetailedSegmentEffort
	resp, err := c.Request(ctx, http.MethodGet, "/segment_efforts", nil, url.Values{
		"segment_id": []string{fmt.Sprintf("%d", segmentID)},
		"per_page":   []string{fmt.Sprintf("%d", perPage)},
	})
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return efforts, c.DecodeResponse(resp, &efforts, http.StatusOK)
}

func (c *Client) DecodeResponse(res *http.Response, v any, expectedCode int) error {
	defer res.Body.Close()

	stravalimit.Update(res.Header)

	if res.StatusCode != expectedCode {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("status code: %d\nbody: %s", res.StatusCode, string(body))
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func (c *Client) Request(ctx context.Context, method string, path string, body any, values url.Values) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	u := fmt.Sprintf("https://strava.com/api/v3/%s", strings.TrimPrefix(path, "/"))
	if len(values) > 0 {
		u += "?" + values.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.AccessToken)

	return c.Client.Do(req)
}
