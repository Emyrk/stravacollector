
// Code generated by 'make site/strava-frontend/src/api/typesGenerated.ts'. DO NOT EDIT.

// From codersdk/athlete.go
export interface AthleteLogin {
  readonly athlete_id: number
  readonly summit: boolean
}

// From codersdk/athlete.go
export interface AthleteSummary {
  readonly athlete_id: number
  readonly summit: boolean
  readonly username: string
  readonly firstname: string
  readonly lastname: string
  readonly sex: string
  readonly profile_pic_link: string
  readonly profile_pic_link_medium: string
  readonly updated_at: string
}

// From codersdk/athlete.go
export interface HugelLeaderBoard {
  readonly personal_best?: HugelLeaderBoardActivity
  readonly activities: HugelLeaderBoardActivity[]
}

// From codersdk/athlete.go
export interface HugelLeaderBoardActivity {
  readonly activity_id: number
  readonly athlete_id: number
  readonly elapsed: number
  readonly rank: number
  readonly efforts: SegmentEffort[]
  readonly athlete: MinAthlete
  readonly activity_name: string
  readonly activity_distance: number
  readonly activity_moving_time: number
  readonly activity_elapsed_time: number
  readonly activity_start_date: string
  readonly activity_total_elevation_gain: number
}

// From codersdk/athlete.go
export interface MinAthlete {
  readonly athlete_id: number
  readonly username: string
  readonly firstname: string
  readonly lastname: string
  readonly sex: string
  readonly profile_pic_link: string
}

// From codersdk/response.go
export interface Response {
  readonly message: string
  readonly detail?: string
}

// From codersdk/athlete.go
export interface SegmentEffort {
  readonly effort_id: number
  readonly start_date: string
  readonly segment_id: number
  readonly elapsed_time: number
  readonly moving_time: number
  readonly device_watts: boolean
  readonly average_watts: number
}
