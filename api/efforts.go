package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
)

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

	activities, err := api.HugelBoardCache.Load(ctx)
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
			athleteAct, err := api.Opts.DB.HugelLeaderboard(ctx, id)
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
