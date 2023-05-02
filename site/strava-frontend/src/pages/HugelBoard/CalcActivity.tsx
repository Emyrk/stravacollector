import { HugelLeaderBoardActivity, SegmentEffort, SegmentSummary } from "../../api/typesGenerated"
import { DistanceToLocal, DistanceToLocalElevation } from "../../lib/Distance/Distance"

export interface ActivityCalResults {
    dateText: string
    elapsedText: string
    totalElapsedText: string
    elevationText: string
    showWatts: boolean
    avgWatts: number
    distance: number
}

export const CalculateActivity = (activity: HugelLeaderBoardActivity): ActivityCalResults => {
    // 2022-11-27T15:42:54Z
    // Dates come over in UTC
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric' };
    const dateText = new Date(activity.activity_start_date).toLocaleDateString(undefined, options)
    const elapsedText = ElapsedDurationText(true, activity.elapsed)
    const totalElapsedText = `${Math.floor(activity.activity_elapsed_time / 3600)}:${Math.floor(activity.activity_elapsed_time / 60) % 60}`
    const elevationText = `${Math.floor(DistanceToLocalElevation(activity.activity_total_elevation_gain) / 100) / 10}k`
    const showWatts = activity.efforts.every(effort => effort.average_watts > 0 && effort.device_watts)
    const avgWatts = Math.floor(activity.efforts.reduce((acc, effort) => acc + effort.average_watts * effort.elapsed_time, 0) / activity.elapsed)
    const distance = Math.floor(DistanceToLocal(activity.activity_distance))


    return {
        dateText,
        elapsedText,
        totalElapsedText,
        elevationText,
        distance,
        showWatts,
        avgWatts
    }
}

export const ElapsedDurationText = (includeHours: boolean, seconds: number): string => {
    const msg = `${Math.floor(seconds / 60) % 60}:${seconds % 60}`
    if (includeHours) {
        return `${Math.floor(seconds / 3600)}:` + msg
    }
    return msg
}

export const SortSegments = (efforts: SegmentSummary[]): SegmentSummary[] => {
    return efforts.sort((a, b) => {
        return a.id - b.id
    })
}

export const SortEfforts = (efforts: SegmentEffort[]): SegmentEffort[] => {
    return efforts.sort((a, b) => {
        return a.segment_id - b.segment_id
    })
}