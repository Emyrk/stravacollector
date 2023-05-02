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
    marginText: string
}

export const CalculateActivity = (activity: HugelLeaderBoardActivity): ActivityCalResults => {
    // 2022-11-27T15:42:54Z
    // Dates come over in UTC
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric' };
    const dateText = new Date(activity.activity_start_date).toLocaleDateString(undefined, options)
    const elapsedText = ElapsedDurationText(activity.elapsed)
    const totalElapsedText = ElapsedDurationText(activity.activity_elapsed_time, true, true, false)
    const elevationText = `${Math.floor(DistanceToLocalElevation(activity.activity_total_elevation_gain) / 100) / 10}k`
    const showWatts = activity.efforts.every(effort => effort.average_watts > 0 && effort.device_watts)
    const avgWatts = Math.floor(activity.efforts.reduce((acc, effort) => acc + effort.average_watts * effort.elapsed_time, 0) / activity.elapsed)
    const distance = Math.floor(DistanceToLocal(activity.activity_distance))
    let marginText = "+" + ElapsedDurationText(activity.elapsed - activity.rank_one_elapsed)
    if (!activity.rank_one_elapsed) {
        marginText = "--:--:--"
    }


    return {
        dateText,
        elapsedText,
        totalElapsedText,
        elevationText,
        distance,
        showWatts,
        avgWatts,
        marginText
    }
}



export const ElapsedDurationText = (seconds: number, h: boolean = true, m: boolean = true, s: boolean = true): string => {
    let msg = []
    if (h) {
        msg.push(`${PaddedNumber(Math.floor(seconds / 3600))}`)
    }
    if (m) {
        msg.push(`${PaddedNumber(Math.floor(seconds / 60) % 60)}`)
    }
    if (s) {
        msg.push(`${PaddedNumber(seconds % 60)}`)
    }
    return msg.join(':')
}

const PaddedNumber = (num: number): string => {
    return num.toString().padStart(2, '0')
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