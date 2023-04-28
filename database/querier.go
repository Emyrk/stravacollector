// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package database

import (
	"context"
)

type sqlcQuerier interface {
	// BestRouteEfforts returns all activities that have efforts on all the provided segments.
	// The returned activities include the best effort for each segment.
	BestRouteEfforts(ctx context.Context, expectedSegments []int64) ([]BestRouteEffortsRow, error)
	DeleteActivity(ctx context.Context, id int64) (ActivitySummary, error)
	GetActivitySummary(ctx context.Context, id int64) (ActivitySummary, error)
	GetAthlete(ctx context.Context, athleteID int64) (Athlete, error)
	GetAthleteLoad(ctx context.Context, athleteID int64) (AthleteLoad, error)
	GetAthleteLogin(ctx context.Context, athleteID int64) (AthleteLogin, error)
	GetAthleteNeedsLoad(ctx context.Context) (GetAthleteNeedsLoadRow, error)
	InsertWebhookDump(ctx context.Context, rawJson string) (WebhookDump, error)
	UpdateActivityName(ctx context.Context, arg UpdateActivityNameParams) error
	UpsertActivityDetail(ctx context.Context, arg UpsertActivityDetailParams) (ActivityDetail, error)
	UpsertActivitySummary(ctx context.Context, arg UpsertActivitySummaryParams) (ActivitySummary, error)
	UpsertAthlete(ctx context.Context, arg UpsertAthleteParams) (Athlete, error)
	UpsertAthleteLoad(ctx context.Context, arg UpsertAthleteLoadParams) (AthleteLoad, error)
	UpsertAthleteLogin(ctx context.Context, arg UpsertAthleteLoginParams) (AthleteLogin, error)
	UpsertMap(ctx context.Context, arg UpsertMapParams) (Map, error)
	UpsertMapSummary(ctx context.Context, arg UpsertMapSummaryParams) (Map, error)
	UpsertSegmentEffort(ctx context.Context, arg UpsertSegmentEffortParams) (SegmentEffort, error)
}

var _ sqlcQuerier = (*sqlQuerier)(nil)
