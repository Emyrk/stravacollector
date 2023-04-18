package strava

import (
	"context"
	"fmt"
	"io"

	"github.com/Emyrk/strava/strava/stravalib"
	"github.com/antihax/optional"
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
func (c *Client) GetSegmentEfforts(ctx context.Context, id int64, perPage int) ([]stravalib.DetailedSegmentEffort, error) {
	perPage32 := int32(perPage)
	segment, resp, err := c.API.SegmentEffortsApi.GetEffortsBySegmentId(c.WithAccess(ctx), int32(id), &stravalib.SegmentEffortsApiGetEffortsBySegmentIdOpts{
		//StartDateLocal: nil,
		//EndDateLocal:   nil,
		PerPage: optional.NewInt32(perPage32),
	})
	if err != nil {
		d, _ := io.ReadAll(resp.Body)
		fmt.Println(string(d), resp.ContentLength, resp.StatusCode, resp)
		return nil, err
	}
	return segment, nil
}

func (c *Client) GetSegmentByID(ctx context.Context, id int64) (stravalib.DetailedSegment, error) {
	segment, _, err := c.API.SegmentsApi.GetSegmentById(c.WithAccess(ctx), id)
	if err != nil {
		return stravalib.DetailedSegment{}, err
	}
	return segment, nil
}
