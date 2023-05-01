package api

import (
	"encoding/json"
	"net/http"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
)

func (api *API) hugelboard(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx                 = r.Context()
		id, athleteLoggedIn = httpmw.AuthenticatedAthleteIDOptional(r)
	)
	var _ = athleteLoggedIn

	activities, err := api.HugelBoardCache.Load(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	board := modelsdk.HugelLeaderBoard{
		PersonalBest: nil,
		Activities:   convertHugelActivities(activities),
	}

	if athleteLoggedIn {
		for _, act := range board.Activities {
			if act.AthleteID == id {
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
		ActivityID: activity.ActivityID,
		AthleteID:  activity.AthleteID,
		Elapsed:    activity.TotalTimeSeconds,
		Rank:       activity.Rank,
		Efforts:    efforts,
		Athlete: modelsdk.MinAthlete{
			AthleteID:      activity.AthleteID,
			Username:       activity.Username,
			Firstname:      activity.Firstname,
			Lastname:       activity.Lastname,
			Sex:            activity.Sex,
			ProfilePicLink: activity.ProfilePicLink,
		},
		ActivityName:               activity.Name,
		ActivityDistance:           activity.Distance,
		ActivityMovingTime:         int64(activity.MovingTime),
		ActivityElapsedTime:        int64(activity.ElapsedTime),
		ActivityStartDate:          activity.StartDate,
		ActivityTotalElevationGain: activity.TotalElevationGain,
	}
}
