package segment

import (
	"context"
	"github.com/Emyrk/strava/strava"
)

type Segment struct {
	ID int64
}

func NewSegment(id int64) *Segment {
	return &Segment{
		ID: id,
	}
}

func (s *Segment) LoadLeaderboard(ctx context.Context, cli *strava.Client)  {
	cli.GetSegmentById()
}