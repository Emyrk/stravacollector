import {
  HugelLeaderBoardActivity,
  SegmentEffort,
  SegmentSummary,
  SuperHugelLeaderBoard,
  SuperHugelLeaderBoardActivity,
} from "../../api/typesGenerated";
import {
  DistanceToLocal,
  DistanceToLocalElevation,
} from "../../lib/Distance/Distance";

export interface ActivityCalResults {
  dateText: string;
  elapsedText: string;
  totalElapsedText: string;
  elevationText: string;
  showWatts: boolean;
  avgWatts: number;
  distance: number;
  marginText: string;
  numActivities: number;
}

export const FormatDate = (
  data: string,
  shortMonth: boolean = false
): string => {
  const options: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: shortMonth ? "short" : "long",
    day: "numeric",
  };
  // 2022-11-27T15:42:54Z
  // Dates come over in UTC
  return new Date(data).toLocaleDateString(undefined, options);
};

export const FormatDateTime = (data: string): string => {
  const options: Intl.DateTimeFormatOptions = {
    year: undefined,
    month: undefined,
    day: undefined,
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: true,
  };
  // 2022-11-27T15:42:54Z
  // Dates come over in UTC
  return new Date(data).toLocaleTimeString(undefined, options);
};

export const CalculateActivity = (
  activity: HugelLeaderBoardActivity | SuperHugelLeaderBoardActivity
): ActivityCalResults => {
  const elapsedText = ElapsedDurationText(activity.elapsed);
  const showWatts = activity.efforts.every(
    (effort) => effort.average_watts > 0 && effort.device_watts
  );
  const avgWatts = Math.floor(
    activity.efforts.reduce(
      (acc, effort) => acc + effort.average_watts * effort.elapsed_time,
      0
    ) / activity.elapsed
  );
  let marginText =
    "+" + ElapsedDurationText(activity.elapsed - activity.rank_one_elapsed);
  if (!activity.rank_one_elapsed) {
    marginText = "--:--:--";
  }
  const uniqueActs = activity.efforts.reduce((acc, effort) => {
    acc[effort.activity_id] += 1;
    return acc;
  }, {} as Record<string, number>);
  const numActivities = Object.keys(uniqueActs).length;

  // Abort early on super hugel activities
  if (!("activity_elapsed_time" in activity)) {
    return {
      elapsedText,
      showWatts,
      avgWatts,
      marginText,
      numActivities,
      dateText: "",
      totalElapsedText: "",
      elevationText: "",
      distance: 0,
    };
  }

  const totalElapsedText = ElapsedDurationText(
    activity.activity_elapsed_time,
    true,
    true,
    true,
    false
  );
  const elevationText = `${
    Math.floor(
      DistanceToLocalElevation(activity.activity_total_elevation_gain) / 100
    ) / 10
  }k`;
  const distance = DistanceToLocal(activity.activity_distance);
  const dateText = FormatDate(activity.activity_start_date, true);

  return {
    dateText,
    elapsedText,
    totalElapsedText,
    elevationText,
    distance,
    showWatts,
    avgWatts,
    marginText,
    numActivities,
  };
};

export const ElapsedDurationText = (
  seconds: number,
  paddedHour: boolean = true,
  h: boolean = true,
  m: boolean = true,
  s: boolean = true
): string => {
  let msg = [];
  if (h) {
    const pad = paddedHour ? 2 : 1;
    msg.push(`${PaddedNumber(Math.floor(seconds / 3600), pad)}`);
  }
  if (m) {
    msg.push(`${PaddedNumber(Math.floor(seconds / 60) % 60)}`);
  }
  if (s) {
    msg.push(`${PaddedNumber(seconds % 60)}`);
  }
  return msg.join(":");
};

const PaddedNumber = (num: number, padLength: number = 2): string => {
  return num.toString().padStart(padLength, "0");
};

export const SortSegments = (efforts: SegmentSummary[]): SegmentSummary[] => {
  return efforts.sort((a, b) => {
    return a.id.toLowerCase() < b.id.toLowerCase() ? -1 : 1;
  });
};

export const SortEfforts = (efforts: SegmentEffort[]): SegmentEffort[] => {
  return efforts.sort((a, b) => {
    return a.segment_id.toLowerCase() < b.segment_id.toLowerCase() ? -1 : 1;
  });
};
