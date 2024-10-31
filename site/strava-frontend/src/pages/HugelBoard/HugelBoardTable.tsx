import { FC, PropsWithChildren, useEffect } from "react";
import {
  Grid,
  Table,
  Thead,
  Tbody,
  Tfoot,
  Tr,
  Th,
  Td,
  TableCaption,
  TableContainer,
  Box,
  Spinner,
  Text,
  Link,
  LinkBox,
  LinkOverlay,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  useToast,
  Alert,
  AlertDescription,
  Flex,
  AlertIcon,
  AlertTitle,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Tooltip,
  GridItem,
  Heading,
  Divider,
  AbsoluteCenter,
} from "@chakra-ui/react";
import {
  getErrorDetail,
  getErrorMessage,
  getHugelSegments,
} from "../../api/rest";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { HugelBoardProps } from "./HugelBoard";
import {
  HugelLeaderBoard,
  HugelLeaderBoardActivity,
  SegmentEffort,
  SegmentSummary,
  SuperHugelLeaderBoardActivity,
} from "../../api/typesGenerated";
import {
  CalculateActivity,
  ElapsedDurationText,
  SortEfforts,
  SortSegments,
} from "./CalcActivity";
import React from "react";
import { CardStat } from "../../components/CardStat/CardStat";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { Loading } from "../../components/Loading/Loading";

export const HugelBoardTable: FC<HugelBoardProps> = ({
  data,
  error,
  isLoading,
  isFetched,
}) => {
  const queryKey = ["hugel-segments"];
  const {
    data: hugelSegments,
    error: hugelSegmentsError,
    isLoading: hugelSegmentsLoading,
    isFetched: hugelSegmentsFetched,
  } = useQuery({
    queryKey,
    queryFn: getHugelSegments,
  });

  if (isLoading || hugelSegmentsLoading) {
    return <Loading />;
  }

  const segmentMapping = hugelSegments?.segments.reduce((acc, segment) => {
    acc[segment.id] = segment;
    return acc;
  }, {} as { [key: string]: SegmentSummary });

  return (
    <>
      {error && (
        <Alert status="error">
          <AlertIcon />
          <AlertTitle>Failed to load leaderboard</AlertTitle>
          <AlertDescription>
            {getErrorMessage(error, "Leaderboard did not load.")}
          </AlertDescription>
        </Alert>
      )}
      <TableContainer>
        <Table size="sm" variant="striped" colorScheme="gray">
          <TableCaption>Das Hügel Results</TableCaption>
          <Thead>
            <Tr>
              <Th>Athlete</Th>
              <Th>Elapsed</Th>
              <Th>Activity</Th>
            </Tr>
          </Thead>
          <Tbody>
            {data?.personal_best && (
              <HugelBoardTableRow
                key={`tbr-${data?.personal_best.athlete_id}`}
                activity={data?.personal_best}
                segmentSummaries={segmentMapping}
                personal={true}
              />
            )}
            {data?.personal_best && (
              <Tr>
                <Td colSpan={6}>
                  <Box position="relative" padding="10">
                    <Divider />
                    <AbsoluteCenter bg="gray.800" px="4">
                      <Heading size="md"> Results </Heading>
                    </AbsoluteCenter>
                  </Box>
                </Td>
              </Tr>
            )}
            {data &&
              data.activities?.map((activity) => {
                return (
                  <HugelBoardTableRow
                    key={`tbr-${activity.athlete_id}`}
                    activity={activity}
                    segmentSummaries={segmentMapping}
                  />
                );
              })}
          </Tbody>
        </Table>
      </TableContainer>
    </>
  );
};

export const HugelBoardTableRow: FC<
  PropsWithChildren<{
    personal?: boolean;
    activity: HugelLeaderBoardActivity | SuperHugelLeaderBoardActivity;
    segmentSummaries?: { [key: string]: SegmentSummary };
  }>
> = ({ personal, activity, segmentSummaries }) => {
  const { firstname, lastname, profile_pic_link, username, hugel_count } =
    activity.athlete;

  // Sort by the length of the segment name to group similar length names
  const efforts = activity.efforts.sort((a, b) => {
    if (
      segmentSummaries &&
      segmentSummaries[a.segment_id] &&
      segmentSummaries[b.segment_id]
    ) {
      return segmentSummaries[a.segment_id].name.toLowerCase() <
        segmentSummaries[b.segment_id].name.toLowerCase()
        ? -1
        : 1;
    }
    return a.segment_id.toLowerCase() < b.segment_id.toLowerCase() ? -1 : 1;
  });

  const pairedEfforts: (SegmentEffort | null)[][] = [];
  for (let i = 0; i < efforts.length; i += 2) {
    pairedEfforts.push(efforts.slice(i, i + 2));
  }
  if (pairedEfforts[pairedEfforts.length - 1].length === 1) {
    pairedEfforts[pairedEfforts.length - 1].push(null);
  }

  const {
    dateText,
    elapsedText,
    totalElapsedText,
    elevationText,
    distance,
    showWatts,
    avgWatts,
    marginText,
  } = CalculateActivity(activity);

  return (
    <Tr key={`row-${activity.athlete_id}`}>
      <Td>
        {personal && (
          <Flex justifyContent={"center"} width={"100%"} paddingBottom={"10px"}>
            <Heading size="md" color={"#fc4c02"}>
              Personal Best
            </Heading>
          </Flex>
        )}
        <Flex pl={3} alignItems={"center"}>
          <Text fontWeight="bold" fontSize={30} pr={5}>
            {activity.rank}
          </Text>
          <AthleteAvatar
            firstname={firstname}
            lastname={lastname}
            athlete_id={activity.athlete_id}
            profile_pic_link={profile_pic_link}
            username={username}
            hugel_count={hugel_count}
            size="md"
            styleProps={{
              mr: 3,
            }}
          />
          <Box>
            <Link
              href={`https://www.strava.com/athletes/${activity.athlete_id}`}
            >
              <Text fontWeight="bold" textAlign="left">
                {firstname} {lastname}
              </Text>
            </Link>
            <Text
              fontSize="sm"
              fontFamily={"monospace"}
              opacity={0.6}
              textAlign="left"
            >
              {dateText}
            </Text>
          </Box>
        </Flex>
      </Td>
      <Td textAlign={"center"}>
        <Text fontSize="lg">{elapsedText}</Text>
        <Text variant="minor">{marginText}</Text>
      </Td>
      {"activity_id" in activity && (
        <Td>
          <Flex direction={"column"}>
            <Text
              textAlign={"center"}
              fontWeight={"bold"}
              maxW="150px"
              isTruncated
            >
              {activity.activity_name}
            </Text>
            <Flex
              alignContent={"center"}
              alignItems={"center"}
              direction={"row"}
            >
              <StravaLink
                m={1}
                href={`https://strava.com/activities/${activity.activity_id}`}
                target="_blank"
                height={"32px"}
                width={"32px"}
                tooltip="View activity on Strava"
              />
              <CardStat
                m={1}
                title="Distance"
                value={distance.toFixed(1) + "mi"}
              />
              <CardStat m={1} title="Elevation" value={elevationText + "ft"} />
            </Flex>
          </Flex>
        </Td>
      )}
      {pairedEfforts.map((efforts, index) => {
        return (
          <EffortPair
            key={`pair-${index}`}
            pair={efforts as [SegmentEffort, SegmentEffort | null]}
            segmentSummaries={segmentSummaries}
          />
        );
      })}
    </Tr>
  );
};

const EffortPair: FC<{
  pair: [SegmentEffort, SegmentEffort | null];
  segmentSummaries?: { [key: string]: SegmentSummary };
}> = ({ pair, segmentSummaries }) => {
  return (
    <Td p={"10px"} m={0}>
      {pair.map((effort, index) => {
        if (!effort) {
          // Blank box of height 2 lines
          return (
            <Box
              key={`effort-${index}`}
              pb={index === 0 ? 3 : 0}
              height="2em"
            ></Box>
          );
        }

        const name =
          segmentSummaries && segmentSummaries[effort.segment_id]
            ? segmentSummaries[effort.segment_id].name
            : "????";

        return (
          <Tooltip
            key={`effort-${index}`}
            mt={index === 0 ? "-10px" : "0px"}
            label={"Link to effort on strava"}
            aria-label="Strava logo tooltip"
          >
            <Link
              target="_blank"
              href={`https://strava.com/activities/${effort.activity_id.toString()}/segments/${effort.effort_id.toString()}`}
            >
              <CardStat
                title={name}
                titleProps={{
                  maxW: "110px",
                  isTruncated: true,
                }}
                pb={index !== 0 ? "0px" : "10px"}
                value={`${ElapsedDurationText(effort.elapsed_time, false)} @ ${
                  effort.device_watts
                    ? effort.average_watts.toFixed(0) + "w"
                    : "--"
                }`}
              />
            </Link>
          </Tooltip>
        );
      })}
    </Td>
  );
};
