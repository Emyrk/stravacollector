package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/internal/hugeldate"
	"github.com/Emyrk/strava/strava"
)

// verifyRoute is janky
func (api *API) verifyRoute(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx                 = r.Context()
		id, athleteLoggedIn = httpmw.AuthenticatedAthleteIDOptional(r)
	)
	routeName := chi.URLParam(r, "route-name")
	verify := chi.URLParam(r, "route-id")

	if routeName != "das-hugel" {
		// Only support hugel for now
		httpapi.Write(ctx, rw, http.StatusNotFound, modelsdk.Response{
			Message: "Route not found",
		})
	}

	if !athleteLoggedIn || id != 2661162 {
		if id != 2661162 {
			httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
				Message: "Not authorized",
			})
			return
		}
	}

	verifyID, err := strconv.ParseInt(verify, 10, 64)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Invalid route id",
			Detail:  err.Error(),
		})
		return
	}

	login, err := api.Opts.DB.GetAthleteLogin(ctx, id)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to load user login",
			Detail:  err.Error(),
		})
		return
	}

	cli := strava.NewOAuthClient(api.OAuthConfig.Client(ctx, login.OAuthToken()))
	stravaRoute, err := cli.GetRoute(ctx, verifyID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed fetch strava route",
			Detail:  err.Error(),
		})
		return
	}

	route, err := api.HugelRouteCache.Load(ctx)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed fetch internal route",
			Detail:  err.Error(),
		})
		return
	}
	hugel := convertRoute(route)
	required := make(map[int64]modelsdk.SegmentSummary)
	for _, seg := range hugel.Segments {
		required[int64(seg.ID)] = seg
	}

	for _, have := range stravaRoute.Segments {
		delete(required, have.ID)
	}

	missingArr := make([]modelsdk.SegmentSummary, 0, len(required))
	for _, missing := range required {
		missingArr = append(missingArr, missing)
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.VerifyRouteResponse{
		MissingSegments: missingArr,
	})
}

func (api *API) competitiveRoute(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	routeName := chi.URLParam(r, "route-name")
	if routeName == "das-hugel" {
		route, err := api.HugelRouteCache.Load(ctx)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
				Message: "Failed to load route",
				Detail:  err.Error(),
			})
			return
		}
		httpapi.Write(ctx, rw, http.StatusOK, convertRoute(route))
		return
	}

	// Only support hugel for now
	httpapi.Write(ctx, rw, http.StatusNotFound, modelsdk.Response{
		Message: "Route not found",
	})
}

func (api *API) superHugelboard(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx                 = r.Context()
		id, athleteLoggedIn = httpmw.AuthenticatedAthleteIDOptional(r)
	)

	activities, err := api.SuperHugelBoardCache.Load(ctx)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to load leaderboard",
			Detail:  err.Error(),
		})
		return
	}

	board := modelsdk.SuperHugelLeaderBoard{
		PersonalBest: nil,
		Activities:   convertSuperHugelActivities(activities),
	}

	if athleteLoggedIn {
		for _, act := range board.Activities {
			if act.AthleteID == modelsdk.StringInt(id) {
				act := act
				board.PersonalBest = &act
				break
			}
		}

		if board.PersonalBest == nil {
			athleteAct, err := api.Opts.DB.SuperHugelLeaderboard(ctx, id)
			if err == nil {
				if len(athleteAct) == 1 {
					pb := convertSuperHugelActivity(athleteAct[0])
					board.PersonalBest = &pb
				}
			}
		}
	}
	httpapi.Write(ctx, rw, http.StatusOK, board)
}

func (api *API) hugelboard(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx                 = r.Context()
		id, athleteLoggedIn = httpmw.AuthenticatedAthleteIDOptional(r)
	)

	before, _ := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	after, _ := strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
	present, _ := strconv.ParseBool(r.URL.Query().Get("present"))
	var beforeTime time.Time
	var afterTime time.Time
	var activities []database.HugelLeaderboardRow
	var err error

	if before > 0 && after > 0 {
		if id != 2661162 {
			httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
				Message: "Not authorized",
			})
			return
		}
		beforeTime = time.Unix(before, 0)
		afterTime = time.Unix(after, 0)
		activities, err = api.Opts.DB.HugelLeaderboard(ctx, database.HugelLeaderboardParams{
			AthleteID: -1,
			Before:    beforeTime,
			After:     afterTime,
		})
	} else {
		if present {
			activities, err = api.HugelBoard2023Cache.Load(ctx)
			beforeTime = hugeldate.StartHugel
			afterTime = hugeldate.EndHugel
		} else {
			activities, err = api.HugelBoardCache.Load(ctx)
		}
	}

	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to load leaderboard",
			Detail:  err.Error(),
		})
		return
	}

	board := modelsdk.HugelLeaderBoard{
		PersonalBest: nil,
		Activities:   convertHugelActivities(activities),
	}

	if athleteLoggedIn {
		for _, act := range board.Activities {
			if act.AthleteID == modelsdk.StringInt(id) {
				act := act
				board.PersonalBest = &act
				break
			}
		}

		if board.PersonalBest == nil {
			athleteAct, err := api.Opts.DB.HugelLeaderboard(ctx, database.HugelLeaderboardParams{
				AthleteID: id,
				Before:    beforeTime,
				After:     afterTime,
			})
			if err == nil {
				if len(athleteAct) == 1 {
					pb := convertHugelActivity(athleteAct[0])
					board.PersonalBest = &pb
				}
			}
		}
	}
	httpapi.Write(ctx, rw, http.StatusOK, board)
}

func convertRoute(route database.GetCompetitiveRouteRow) modelsdk.CompetitiveRoute {
	sdkRoute := modelsdk.CompetitiveRoute{
		Name:        route.Name,
		DisplayName: route.DisplayName,
		Description: route.Description,
		Segments:    []modelsdk.SegmentSummary{},
	}

	_ = json.Unmarshal(route.SegmentSummaries, &sdkRoute.Segments)
	return sdkRoute
}

func convertSuperHugelActivities(activites []database.SuperHugelLeaderboardRow) []modelsdk.SuperHugelLeaderBoardActivity {
	sdk := make([]modelsdk.SuperHugelLeaderBoardActivity, 0, len(activites))
	for _, act := range activites {
		sdk = append(sdk, convertSuperHugelActivity(act))
	}
	return sdk
}

func convertSuperHugelActivity(activity database.SuperHugelLeaderboardRow) modelsdk.SuperHugelLeaderBoardActivity {
	var efforts []modelsdk.SegmentEffort
	_ = json.Unmarshal(activity.Efforts, &efforts)
	return modelsdk.SuperHugelLeaderBoardActivity{
		RankOneElapsed: activity.BestTime,
		AthleteID:      modelsdk.StringInt(activity.AthleteID),
		Elapsed:        activity.TotalTimeSeconds,
		Rank:           activity.Rank,
		Efforts:        efforts,
		Athlete: modelsdk.MinAthlete{
			AthleteID:      modelsdk.StringInt(activity.AthleteID),
			Username:       activity.Username,
			Firstname:      activity.Firstname,
			Lastname:       activity.Lastname,
			Sex:            activity.Sex,
			ProfilePicLink: activity.ProfilePicLink,
			HugelCount:     int(activity.HugelCount),
		},
	}
}

func convertHugelActivities(activites []database.HugelLeaderboardRow) []modelsdk.HugelLeaderBoardActivity {
	sdk := make([]modelsdk.HugelLeaderBoardActivity, 0, len(activites))
	for _, act := range activites {
		sdk = append(sdk, convertHugelActivity(act))
	}
	return sdk
}

func convertHugelActivity(activity database.HugelLeaderboardRow) modelsdk.HugelLeaderBoardActivity {
	var efforts []modelsdk.SegmentEffort
	_ = json.Unmarshal(activity.Efforts, &efforts)
	return modelsdk.HugelLeaderBoardActivity{
		RankOneElapsed: activity.BestTime,
		ActivityID:     modelsdk.StringInt(activity.ActivityID),
		AthleteID:      modelsdk.StringInt(activity.AthleteID),
		Elapsed:        activity.TotalTimeSeconds,
		Rank:           activity.Rank,
		Efforts:        efforts,
		Athlete: modelsdk.MinAthlete{
			AthleteID:      modelsdk.StringInt(activity.AthleteID),
			Username:       activity.Username,
			Firstname:      activity.Firstname,
			Lastname:       activity.Lastname,
			Sex:            activity.Sex,
			ProfilePicLink: activity.ProfilePicLink,
			HugelCount:     int(activity.HugelCount),
		},
		ActivityName:               activity.Name,
		ActivityDistance:           activity.Distance,
		ActivityMovingTime:         int64(activity.MovingTime),
		ActivityElapsedTime:        int64(activity.ElapsedTime),
		ActivityStartDate:          activity.StartDate,
		ActivityTotalElevationGain: activity.TotalElevationGain,
	}
}
