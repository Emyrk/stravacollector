package dbmetrics

import (
	"context"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

var (
	// Force these imports, for some reason the autogen does not include them.
	_ uuid.UUID
	_ context.Context
	_ pgx.Conn
)

// NewQueryMetrics returns a database.Store that registers metrics for all queries to reg.
func NewQueryMetrics(s database.Store, logger zerolog.Logger, reg prometheus.Registerer) database.Store {
	queryLatencies := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "coderd",
		Subsystem: "db",
		Name:      "query_latencies_seconds",
		Help:      "Latency distribution of queries in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"query"})
	reg.MustRegister(queryLatencies)
	return &queryMetricsStore{
		s:              s,
		queryLatencies: queryLatencies,
		dbMetrics:      NewDBMetrics(s, logger, reg).(*metricsStore),
	}
}

var _ database.Store = (*queryMetricsStore)(nil)

type queryMetricsStore struct {
	s              database.Store
	queryLatencies *prometheus.HistogramVec
	dbMetrics      *metricsStore
}

func (m queryMetricsStore) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	duration, err := m.s.Ping(ctx)
	m.queryLatencies.WithLabelValues("Ping").Observe(time.Since(start).Seconds())
	return duration, err
}

func (m queryMetricsStore) InTx(f func(database.Store) error, options *pgx.TxOptions) error {
	return m.dbMetrics.InTx(f, options)
}

func (m queryMetricsStore) Close() error {
	return m.dbMetrics.Close()
}

func (m queryMetricsStore) AllCompetitiveRoutes(ctx context.Context) ([]database.CompetitiveRoute, error) {
	start := time.Now()
	r0, r1 := m.s.AllCompetitiveRoutes(ctx)
	m.queryLatencies.WithLabelValues("AllCompetitiveRoutes").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) AthleteHugelActivites(ctx context.Context, athleteID int64) ([]database.AthleteHugelActivitesRow, error) {
	start := time.Now()
	r0, r1 := m.s.AthleteHugelActivites(ctx, athleteID)
	m.queryLatencies.WithLabelValues("AthleteHugelActivites").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) AthleteSyncedActivities(ctx context.Context, arg database.AthleteSyncedActivitiesParams) ([]database.AthleteSyncedActivitiesRow, error) {
	start := time.Now()
	r0, r1 := m.s.AthleteSyncedActivities(ctx, arg)
	m.queryLatencies.WithLabelValues("AthleteSyncedActivities").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) BestRouteEfforts(ctx context.Context, expectedSegments []int64) ([]database.BestRouteEffortsRow, error) {
	start := time.Now()
	r0, r1 := m.s.BestRouteEfforts(ctx, expectedSegments)
	m.queryLatencies.WithLabelValues("BestRouteEfforts").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) DeleteActivity(ctx context.Context, id int64) (database.ActivitySummary, error) {
	start := time.Now()
	r0, r1 := m.s.DeleteActivity(ctx, id)
	m.queryLatencies.WithLabelValues("DeleteActivity").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) DeleteAthleteLogin(ctx context.Context, athleteID int64) error {
	start := time.Now()
	r0 := m.s.DeleteAthleteLogin(ctx, athleteID)
	m.queryLatencies.WithLabelValues("DeleteAthleteLogin").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) EddingtonActivities(ctx context.Context, athleteID int64) ([]database.EddingtonActivitiesRow, error) {
	start := time.Now()
	r0, r1 := m.s.EddingtonActivities(ctx, athleteID)
	m.queryLatencies.WithLabelValues("EddingtonActivities").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetActivityDetail(ctx context.Context, id int64) (database.ActivityDetail, error) {
	start := time.Now()
	r0, r1 := m.s.GetActivityDetail(ctx, id)
	m.queryLatencies.WithLabelValues("GetActivityDetail").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetActivitySummary(ctx context.Context, id int64) (database.ActivitySummary, error) {
	start := time.Now()
	r0, r1 := m.s.GetActivitySummary(ctx, id)
	m.queryLatencies.WithLabelValues("GetActivitySummary").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthlete(ctx context.Context, athleteID int64) (database.Athlete, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthlete(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthlete").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteFull(ctx context.Context, athleteID int64) (database.GetAthleteFullRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteFull(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthleteFull").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteLoad(ctx context.Context, athleteID int64) (database.AthleteForwardLoad, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteLoad(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthleteLoad").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteLoadDetailed(ctx context.Context, athleteID int64) (database.GetAthleteLoadDetailedRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteLoadDetailed(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthleteLoadDetailed").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteLogin(ctx context.Context, athleteID int64) (database.AthleteLogin, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteLogin(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthleteLogin").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteLoginFull(ctx context.Context, athleteID int64) (database.GetAthleteLoginFullRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteLoginFull(ctx, athleteID)
	m.queryLatencies.WithLabelValues("GetAthleteLoginFull").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetAthleteNeedsForwardLoad(ctx context.Context) ([]database.GetAthleteNeedsForwardLoadRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetAthleteNeedsForwardLoad(ctx)
	m.queryLatencies.WithLabelValues("GetAthleteNeedsForwardLoad").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetBestPersonalSegmentEffort(ctx context.Context, arg database.GetBestPersonalSegmentEffortParams) ([]database.SegmentEffort, error) {
	start := time.Now()
	r0, r1 := m.s.GetBestPersonalSegmentEffort(ctx, arg)
	m.queryLatencies.WithLabelValues("GetBestPersonalSegmentEffort").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetCompetitiveRoute(ctx context.Context, routeName string) (database.GetCompetitiveRouteRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetCompetitiveRoute(ctx, routeName)
	m.queryLatencies.WithLabelValues("GetCompetitiveRoute").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) GetSegments(ctx context.Context, segmentIds []int64) ([]database.GetSegmentsRow, error) {
	start := time.Now()
	r0, r1 := m.s.GetSegments(ctx, segmentIds)
	m.queryLatencies.WithLabelValues("GetSegments").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) HugelLeaderboard(ctx context.Context, arg database.HugelLeaderboardParams) ([]database.HugelLeaderboardRow, error) {
	start := time.Now()
	r0, r1 := m.s.HugelLeaderboard(ctx, arg)
	m.queryLatencies.WithLabelValues("HugelLeaderboard").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) IncrementActivitySummaryDownload(ctx context.Context, id int64) error {
	start := time.Now()
	r0 := m.s.IncrementActivitySummaryDownload(ctx, id)
	m.queryLatencies.WithLabelValues("IncrementActivitySummaryDownload").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) InsertFailedJob(ctx context.Context, rawJson string) (database.FailedJob, error) {
	start := time.Now()
	r0, r1 := m.s.InsertFailedJob(ctx, rawJson)
	m.queryLatencies.WithLabelValues("InsertFailedJob").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) InsertWebhookDump(ctx context.Context, rawJson string) (database.WebhookDump, error) {
	start := time.Now()
	r0, r1 := m.s.InsertWebhookDump(ctx, rawJson)
	m.queryLatencies.WithLabelValues("InsertWebhookDump").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) LoadedSegments(ctx context.Context) ([]database.LoadedSegmentsRow, error) {
	start := time.Now()
	r0, r1 := m.s.LoadedSegments(ctx)
	m.queryLatencies.WithLabelValues("LoadedSegments").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) MissingHugelSegments(ctx context.Context, activityID int64) ([]database.Segment, error) {
	start := time.Now()
	r0, r1 := m.s.MissingHugelSegments(ctx, activityID)
	m.queryLatencies.WithLabelValues("MissingHugelSegments").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) MissingSegments(ctx context.Context, activitiesID int64) ([]string, error) {
	start := time.Now()
	r0, r1 := m.s.MissingSegments(ctx, activitiesID)
	m.queryLatencies.WithLabelValues("MissingSegments").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) NeedsARefresh(ctx context.Context) ([]database.NeedsARefreshRow, error) {
	start := time.Now()
	r0, r1 := m.s.NeedsARefresh(ctx)
	m.queryLatencies.WithLabelValues("NeedsARefresh").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) RefreshHugel2023Activities(ctx context.Context) error {
	start := time.Now()
	r0 := m.s.RefreshHugel2023Activities(ctx)
	m.queryLatencies.WithLabelValues("RefreshHugel2023Activities").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) RefreshHugelActivities(ctx context.Context) error {
	start := time.Now()
	r0 := m.s.RefreshHugelActivities(ctx)
	m.queryLatencies.WithLabelValues("RefreshHugelActivities").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) RefreshHugelLiteActivities(ctx context.Context) error {
	start := time.Now()
	r0 := m.s.RefreshHugelLiteActivities(ctx)
	m.queryLatencies.WithLabelValues("RefreshHugelLiteActivities").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) RefreshSuperHugelActivities(ctx context.Context) error {
	start := time.Now()
	r0 := m.s.RefreshSuperHugelActivities(ctx)
	m.queryLatencies.WithLabelValues("RefreshSuperHugelActivities").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) StarSegments(ctx context.Context, arg database.StarSegmentsParams) error {
	start := time.Now()
	r0 := m.s.StarSegments(ctx, arg)
	m.queryLatencies.WithLabelValues("StarSegments").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) SuperHugelLeaderboard(ctx context.Context, athleteID interface{}) ([]database.SuperHugelLeaderboardRow, error) {
	start := time.Now()
	r0, r1 := m.s.SuperHugelLeaderboard(ctx, athleteID)
	m.queryLatencies.WithLabelValues("SuperHugelLeaderboard").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) TotalActivityDetailsCount(ctx context.Context) (int64, error) {
	start := time.Now()
	r0, r1 := m.s.TotalActivityDetailsCount(ctx)
	m.queryLatencies.WithLabelValues("TotalActivityDetailsCount").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) TotalJobCount(ctx context.Context) (int64, error) {
	start := time.Now()
	r0, r1 := m.s.TotalJobCount(ctx)
	m.queryLatencies.WithLabelValues("TotalJobCount").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) TotalRideActivitySummariesCount(ctx context.Context) (int64, error) {
	start := time.Now()
	r0, r1 := m.s.TotalRideActivitySummariesCount(ctx)
	m.queryLatencies.WithLabelValues("TotalRideActivitySummariesCount").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpdateActivityName(ctx context.Context, arg database.UpdateActivityNameParams) error {
	start := time.Now()
	r0 := m.s.UpdateActivityName(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateActivityName").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) UpdateActivityType(ctx context.Context, arg database.UpdateActivityTypeParams) error {
	start := time.Now()
	r0 := m.s.UpdateActivityType(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateActivityType").Observe(time.Since(start).Seconds())
	return r0
}

func (m queryMetricsStore) UpsertActivityDetail(ctx context.Context, arg database.UpsertActivityDetailParams) (database.ActivityDetail, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertActivityDetail(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertActivityDetail").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertActivitySummary(ctx context.Context, arg database.UpsertActivitySummaryParams) (database.ActivitySummary, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertActivitySummary(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertActivitySummary").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertAthlete(ctx context.Context, arg database.UpsertAthleteParams) (database.Athlete, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertAthlete(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertAthlete").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertAthleteEddington(ctx context.Context, arg database.UpsertAthleteEddingtonParams) (database.AthleteEddington, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertAthleteEddington(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertAthleteEddington").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertAthleteForwardLoad(ctx context.Context, arg database.UpsertAthleteForwardLoadParams) (database.AthleteForwardLoad, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertAthleteForwardLoad(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertAthleteForwardLoad").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertAthleteLogin(ctx context.Context, arg database.UpsertAthleteLoginParams) (database.AthleteLogin, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertAthleteLogin(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertAthleteLogin").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertMapData(ctx context.Context, arg database.UpsertMapDataParams) (database.Map, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertMapData(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertMapData").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertSegment(ctx context.Context, arg database.UpsertSegmentParams) (database.Segment, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertSegment(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertSegment").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) UpsertSegmentEffort(ctx context.Context, arg database.UpsertSegmentEffortParams) (database.SegmentEffort, error) {
	start := time.Now()
	r0, r1 := m.s.UpsertSegmentEffort(ctx, arg)
	m.queryLatencies.WithLabelValues("UpsertSegmentEffort").Observe(time.Since(start).Seconds())
	return r0, r1
}

func (m queryMetricsStore) YearlyHugelLeaderboard(ctx context.Context, arg database.YearlyHugelLeaderboardParams) ([]database.HugelLeaderboardRow, error) {
	start := time.Now()
	r0, r1 := m.s.YearlyHugelLeaderboard(ctx, arg)
	m.queryLatencies.WithLabelValues("YearlyHugelLeaderboard").Observe(time.Since(start).Seconds())
	return r0, r1
}
