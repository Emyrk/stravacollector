package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	AccessToken string
	Client      *http.Client
}

func New(accessToken string) *Client {
	return &Client{
		AccessToken: accessToken,
		Client:      http.DefaultClient,
	}
}

func (c *Client) AthleteSegmentEfforts(ctx context.Context, segmentID int, perPage int) ([]SegmentEffort, error) {
	var efforts []SegmentEffort
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
	fmt.Println(req.URL.String())

	return c.Client.Do(req)
}
