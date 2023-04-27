package modelsdk

type AthleteLogin struct {
	AthleteID int64 `db:"athlete_id" json:"athlete_id"`
	Summit    bool  `db:"summit" json:"summit"`
}