package river

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Emyrk/strava/api/hugelhelp"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/internal/hugeldate"
	"github.com/Emyrk/strava/strava"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverlog"
	"github.com/riverqueue/river/rivertype"
)

func (m *Manager) EnqueueFetchActivity(ctx context.Context, source database.ActivityDetailSource, athleteID int64, activityID int64, hugelPotential bool, onHugelDates bool, priority int, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{
		Priority: priority,
		Tags:     []string{fmt.Sprintf("%d", athleteID), fmt.Sprintf("%d", activityID)},
	}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(ctx, FetchActivityArgs{
		ActivityID:     activityID,
		AthleteID:      athleteID,
		Source:         source,
		HugelPotential: hugelPotential,
		OnHugelDates:   onHugelDates,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type FetchActivityArgs struct {
	Source     database.ActivityDetailSource `json:"source"`
	ActivityID int64                         `json:"activity_id"`
	AthleteID  int64                         `json:"athlete_id"`
	// HugelPotential is a boolean that helps filter which events to
	// sync during the hugel event.
	// TODO: Remove this after november.
	HugelPotential bool `json:"can_be_hugel"`
	OnHugelDates   bool `json:"on_hugel_dates"`
}

func (FetchActivityArgs) Kind() string { return "fetch_activity" }
func (FetchActivityArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: riverStravaQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 5,
		},
	}
}

type FetchActivityWorker struct {
	mgr *Manager
	river.WorkerDefaults[FetchActivityArgs]
}

func (*FetchActivityWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w *FetchActivityWorker) Work(ctx context.Context, job *river.Job[FetchActivityArgs]) error {
	now := time.Now().In(hugeldate.CentralTimeZone)
	logger := jobLogFields(w.mgr.logger, job)

	args := job.Args
	if args.Source == "" {
		args.Source = database.ActivityDetailSourceUnknown
	}

	logger = logger.With().
		Int64("activity_id", args.ActivityID).
		Int64("athlete_id", args.AthleteID).
		Str("source", string(args.Source)).
		Logger()

	adjustInt := int64(0)
	adjustDaily := int64(0)
	if args.HugelPotential {
		adjustInt = 5
		adjustDaily = 50
	}
	if args.Source == database.ActivityDetailSourceManual {
		adjustInt = 10
		adjustDaily = 115
	}

	// Snooze the job if not relevant during the hugel event.
	if hugelhelp.HugelOngoing(now) {
		// Manyal events are ok
		if args.Source != database.ActivityDetailSourceManual {
			// Only sync hugel potential activities during the event.
			if !args.HugelPotential && !args.OnHugelDates {
				// Wait a day before trying again.
				_ = river.RecordOutput(ctx, "hugel event is ongoing, and this activity is not relevant, snoozing job")
				return river.JobSnooze(time.Hour * 24)
			}
		}
	}

	err := w.mgr.jobStravaCheck(logger, 1, adjustInt, adjustDaily)
	if err != nil {
		return w.mgr.StravaSnooze(ctx)
	}

	// Only track athletes we have in our database
	athlete, err := w.mgr.db.GetAthleteLogin(ctx, args.AthleteID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Error().Msg("athlete not found, job abandoned")
		return river.RecordOutput(ctx, "athlete not found, job abandoned")
	}
	if err != nil {
		return err
	}

	// First check if we just fetched this from another source.
	act, err := w.mgr.db.GetActivityDetail(ctx, args.ActivityID)
	if err == nil {
		// Already fetched this activity.
		if !(args.Source == database.ActivityDetailSourceManual || args.Source == database.ActivityDetailSourceZeroSegmentRefetch) {
			// Manual and zero segment refetches are always allowed to refetch.
			// Others are aborted if the activity was updated in the last 24 hours.
			if time.Since(act.UpdatedAt.Time) < time.Hour*24 {
				return river.RecordOutput(ctx, "activity already fetched, skipping")
			}
		}
	}

	cli := strava.NewOAuthClient(w.mgr.oauthCfg.Client(ctx, athlete.OAuthToken()))
	activity, err := cli.GetActivity(ctx, args.ActivityID, true)
	if err != nil {
		se := strava.IsAPIError(err)
		if se != nil && se.Response.StatusCode != http.StatusTooManyRequests {
			// Kill the job, since we can't fetch this activity due to some other error.
			// Insert the error to review later.
			jobData, _ := json.Marshal(failedJob[FetchActivityArgs]{
				Job:   job,
				Args:  args,
				Error: err.Error(),
			})

			// Insert a failed job for debugging if not expected.
			if se.Response.StatusCode == http.StatusNotFound {
				// No activity? Just drop the job, nothing to do.
				return river.RecordOutput(ctx, fmt.Sprintf("activity not found: https://www.strava.com/activities/%d", args.ActivityID))
			}

			if se.Response.StatusCode == http.StatusBadGateway && strings.Contains(string(se.Body), "Strava is temporarily unavailable") {
				// TODO: Pause the queue and awake it later.
				_ = river.RecordOutput(ctx, "strava is temporarily unavailable, retrying later")
				return w.mgr.StravaSnooze(ctx)
			}

			_, _ = w.mgr.db.InsertFailedJob(ctx, string(jobData))
			return river.RecordOutput(ctx, fmt.Sprintf("failed to fetch: %+v", se))
		}
		if se != nil && se.Response.StatusCode == http.StatusTooManyRequests {
			return w.mgr.StravaSnooze(ctx)
		}
		return err
	}

	logger.Debug().
		Int64("activity_id", activity.ID).
		Int("segment_count", len(activity.SegmentEfforts)).
		Msg("activity fetched")

	err = w.insert(ctx, activity, athlete, args)
	if err != nil {
		return err
	}

	// Potentially re-fetch the activity if it has 0 segment efforts.
	// Strava is slow sometimes. If less than 5 miles though, just ignore it.
	if len(activity.SegmentEfforts) == 0 && database.DistanceToMiles(activity.Distance) > 5 {
		summary, err := w.mgr.db.GetActivitySummary(ctx, activity.ID)
		if err == nil {
			if summary.DownloadCount == 0 {
				// If 0 segment efforts, and has never been redownloaded. Retry in 1hr.
				// This might be correct, but we should check again.
				_, err := w.mgr.EnqueueFetchActivity(ctx, database.ActivityDetailSourceZeroSegmentRefetch, args.AthleteID, args.ActivityID, args.HugelPotential, args.OnHugelDates, job.Priority, func(opt *river.InsertOpts) {
					opt.ScheduledAt = time.Now().Add(time.Hour * 2)
				})
				if err != nil {
					logger.Error().
						Err(err).
						Int("activity_id", int(args.ActivityID)).
						Int("athlete_id", int(args.AthleteID)).
						Msg("error re-enqueuing activity with 0 segments")
				}
				riverlog.Logger(ctx).Info("Activity had 0 segments, re-enqueued for later fetch")
			}
		}
	}

	return river.RecordOutput(ctx, map[string]any{
		"segments": len(activity.SegmentEfforts),
		"link":     fmt.Sprintf("https://www.strava.com/activities/%d", activity.ID),
	})
}

func (w *FetchActivityWorker) insert(ctx context.Context, activity strava.DetailedActivity, athlete database.AthleteLogin, args FetchActivityArgs) error {
	// Parse the activity, save all efforts.
	err := w.mgr.db.InTx(func(store database.Store) error {
		_, err := store.UpsertMapData(ctx, database.UpsertMapDataParams{
			ID:              activity.Map.ID,
			Polyline:        activity.Map.Polyline,
			SummaryPolyline: activity.Map.SummaryPolyline,
		})
		if err != nil {
			return fmt.Errorf("upsert map: %w", err)
		}

		_, err = store.UpsertActivitySummary(ctx, database.UpsertActivitySummaryParams{
			ID:                 activity.ID,
			AthleteID:          activity.Athlete.ID,
			UploadID:           activity.UploadID,
			ExternalID:         activity.ExternalID,
			Name:               activity.Name,
			Distance:           activity.Distance,
			MovingTime:         activity.MovingTime,
			ElapsedTime:        activity.ElapsedTime,
			TotalElevationGain: activity.TotalElevationGain,
			ActivityType:       activity.Type,
			SportType:          activity.SportType,
			WorkoutType:        activity.WorkoutType,
			StartDate:          database.Timestamptz(activity.StartDate),
			StartDateLocal:     database.Timestamptz(activity.StartDateLocal),
			Timezone:           activity.Timezone,
			UtcOffset:          activity.UtcOffset,
			AchievementCount:   activity.AchievementCount,
			KudosCount:         activity.KudosCount,
			CommentCount:       activity.CommentCount,
			AthleteCount:       activity.AthleteCount,
			PhotoCount:         activity.PhotoCount,
			MapID:              activity.Map.ID,
			Trainer:            activity.Trainer,
			Commute:            activity.Commute,
			Manual:             activity.Manual,
			Private:            activity.Private,
			Flagged:            activity.Flagged,
			GearID:             activity.GearID,
			AverageSpeed:       activity.AverageSpeed,
			MaxSpeed:           activity.MaxSpeed,
			DeviceWatts:        activity.DeviceWatts,
			HasHeartrate:       activity.HasHeartrate,
			PrCount:            activity.PrCount,
			TotalPhotoCount:    activity.TotalPhotoCount,
			AverageHeartrate:   activity.AverageHeartrate,
			MaxHeartrate:       activity.MaxHeartrate,
		})
		if err != nil {
			return fmt.Errorf("upsert activity summary: %w", err)
		}

		_, err = store.UpsertActivityDetail(ctx, database.UpsertActivityDetailParams{
			ID:                       activity.ID,
			AthleteID:                activity.Athlete.ID,
			StartLatlng:              activity.StartLatlng,
			EndLatlng:                activity.EndLatlng,
			MapID:                    activity.Map.ID,
			FromAcceptedTag:          activity.FromAcceptedTag,
			AverageCadence:           activity.AverageCadence,
			AverageTemp:              activity.AverageTemp,
			AverageWatts:             activity.AverageWatts,
			WeightedAverageWatts:     activity.WeightedAverageWatts,
			Kilojoules:               activity.Kilojoules,
			MaxWatts:                 activity.MaxWatts,
			ElevHigh:                 activity.ElevHigh,
			ElevLow:                  activity.ElevLow,
			SufferScore:              int32(activity.SufferScore),
			EmbedToken:               activity.EmbedToken,
			SegmentLeaderboardOptOut: activity.SegmentLeaderboardOptOut,
			LeaderboardOptOut:        activity.LeaderboardOptOut,
			Calories:                 activity.Calories,
			Source:                   args.Source,
			//
			PremiumFetch:      athlete.Summit,
			NumSegmentEfforts: int32(len(activity.SegmentEfforts)),
		})
		if err != nil {
			return fmt.Errorf("upsert activity details: %w", err)
		}

		err = store.IncrementActivitySummaryDownload(ctx, activity.ID)
		if err != nil {
			return fmt.Errorf("increment download count: %w", err)
		}

		// Insert efforts.
		starAtheletes := make([]int64, 0, len(activity.SegmentEfforts))
		starSegments := make([]int64, 0, len(activity.SegmentEfforts))
		starStarred := make([]bool, 0, len(activity.SegmentEfforts))
		starSegmentsAdded := make(map[int64]bool)
		for i, effort := range activity.SegmentEfforts {
			if _, ok := starSegmentsAdded[effort.Segment.ID]; !ok {
				starAtheletes = append(starAtheletes, effort.Athlete.ID)
				starSegments = append(starSegments, effort.Segment.ID)
				starStarred = append(starStarred, effort.Segment.Starred)
				starSegmentsAdded[effort.Segment.ID] = true
			}
			_, err := store.UpsertSegmentEffort(ctx, database.UpsertSegmentEffortParams{
				ID:             effort.ID,
				AthleteID:      effort.Athlete.ID,
				SegmentID:      effort.Segment.ID,
				Name:           effort.Name,
				ElapsedTime:    effort.ElapsedTime,
				MovingTime:     effort.MovingTime,
				StartDate:      database.Timestamptz(effort.StartDate),
				StartDateLocal: database.Timestamptz(effort.StartDateLocal),
				Distance:       effort.Distance,
				StartIndex:     effort.StartIndex,
				EndIndex:       effort.EndIndex,
				DeviceWatts:    effort.DeviceWatts,
				AverageWatts:   effort.AverageWatts,
				KomRank: pgtype.Int4{
					Int32: effort.KomRank,
					Valid: effort.KomRank != 0,
				},
				PrRank: pgtype.Int4{
					Int32: effort.PrRank,
					Valid: effort.PrRank != 0,
				},
				ActivitiesID: activity.ID,
			})
			if err != nil {
				return fmt.Errorf("upsert segment effort index=%d, id=%d: %w", i, effort.ID, err)
			}
		}

		err = store.StarSegments(ctx, database.StarSegmentsParams{
			AthleteID: starAtheletes,
			SegmentID: starSegments,
			Starred:   starStarred,
		})
		if err != nil {
			return fmt.Errorf("star segments: %w", err)
		}

		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	return nil
}
