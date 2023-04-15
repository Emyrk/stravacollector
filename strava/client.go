package strava

import (
	"context"

	"github.com/Emyrk/strava/strava/stravalib"
)

type Client struct {
	API         *stravalib.APIClient
	AccessToken string
}

func New(accessToken string) *Client {
	cfg := stravalib.NewConfiguration()
	api := stravalib.NewAPIClient(cfg)

	return &Client{
		API:         api,
		AccessToken: accessToken,
	}
}

func (c *Client) WithAccess(ctx context.Context) context.Context {
	return context.WithValue(ctx, stravalib.ContextAccessToken, c.AccessToken)
}

func (c *Client) GetSegmentById(ctx context.Context, id int64) (stravalib.DetailedSegment, error) {
	segment, _, err := c.API.SegmentsApi.GetSegmentById(c.WithAccess(ctx), id)
	if err != nil {
		return stravalib.DetailedSegment{}, err
	}
	return segment, nil
}
