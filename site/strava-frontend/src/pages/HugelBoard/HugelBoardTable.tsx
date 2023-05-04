import { FC, PropsWithChildren, useEffect } from "react"
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
  Tabs, TabList, TabPanels, Tab, TabPanel,
} from '@chakra-ui/react'
import { getErrorDetail, getErrorMessage, getHugelLeaderBoard, getHugelSegments } from "../../api/rest"
import { useMutation, useQuery } from "@tanstack/react-query";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { HugelBoardProps } from "./HugelBoard";
import { HugelLeaderBoard, HugelLeaderBoardActivity, SegmentEffort, SegmentSummary, SuperHugelLeaderBoardActivity } from "../../api/typesGenerated"
import { CalculateActivity, ElapsedDurationText, SortEfforts, SortSegments } from "./CalcActivity";
import React from "react";

export const HugelBoardTable: FC<HugelBoardProps> = ({
  data, error, isLoading, isFetched
}) => {
  const queryKey = ["hugel-segments"]
  const {
    data: hugelSegments,
    error: hugelSegmentsError,
    isLoading: hugelSegmentsLoading,
    isFetched: hugelSegmentsFetched,
  } = useQuery({
    queryKey,
    queryFn: getHugelSegments,
  })

  const segmentMapping = hugelSegments?.segments.reduce((acc, segment) => {
    acc[segment.id] = segment
    return acc
  }, {} as { [key: string]: SegmentSummary })

  return <>
    {
      error && <Alert status='error'>
        <AlertIcon />
        <AlertTitle>Failed to load leaderboard</AlertTitle>
        <AlertDescription>{getErrorMessage(error, "Leaderboard did not load.")}</AlertDescription>
      </Alert>
    }
    <TableContainer>
      <Table size="sm" variant='striped' colorScheme='gray'>
        <TableCaption>Das Hugel Leaderboad</TableCaption>
        <Thead>
          <Tr>
            <Th>Athlete</Th>
            <Th>Elapsed</Th>
            <Th>Activity</Th>
          </Tr>
        </Thead>
        <Tbody>
          {
            data && data.activities?.map((activity) => {
              return <HugelBoardTableRow key={`tbr-${activity.athlete_id}`} activity={activity} segmentSummaries={segmentMapping} />
            })
          }
        </Tbody>
      </Table>
    </TableContainer >
  </>
}



export const HugelBoardTableRow: FC<PropsWithChildren<{
  activity: HugelLeaderBoardActivity | SuperHugelLeaderBoardActivity
  segmentSummaries?: { [key: string]: SegmentSummary }
}>> = ({ activity, segmentSummaries }) => {
  const {
    firstname,
    lastname,
    profile_pic_link,
    username,
  } = activity.athlete

  // Sort by the length of the segment name to group similar length names
  const efforts = activity.efforts.sort((a, b) => {
    if (segmentSummaries && segmentSummaries[a.segment_id] && segmentSummaries[b.segment_id]) {
      return segmentSummaries[a.segment_id].name.toLowerCase() < segmentSummaries[b.segment_id].name.toLowerCase() ? -1 : 1
    }
    return a.segment_id.toLowerCase() < b.segment_id.toLowerCase() ? -1 : 1
  })


  const pairedEfforts: (SegmentEffort | null)[][] = []
  for (let i = 0; i < efforts.length; i += 2) {
    pairedEfforts.push(efforts.slice(i, i + 2))
  }
  if (pairedEfforts[pairedEfforts.length - 1].length === 1) {
    pairedEfforts[pairedEfforts.length - 1].push(null)
  }

  const {
    dateText,
    elapsedText,
    totalElapsedText,
    elevationText,
    distance,
    showWatts,
    avgWatts,
    marginText
  } = CalculateActivity(activity)

  return <Tr key={`row-${activity.athlete_id}`}>
    <Td>
      <Flex p={3} alignItems={'center'}>
        <Text fontWeight='bold' fontSize={30} pr={5}>
          {activity.rank}
        </Text>
        <AthleteAvatar
          firstName={firstname}
          lastName={lastname}
          athleteID={activity.athlete_id}
          profilePicLink={profile_pic_link}
          username={username}
          size="md"
          styleProps={{
            mr: 3,
          }}
        />
        <Box>
          <Text fontWeight='bold' textAlign='left' >
            {firstname} {lastname}
          </Text>
          <Text fontSize='sm' fontFamily={'monospace'} opacity={.6} textAlign='left' >{dateText}</Text>
        </Box>
      </Flex>
    </Td>
    <Td textAlign={"center"}>
      <Text fontSize='lg'>{elapsedText}</Text>
      <Text variant="minor">{marginText}</Text>
    </Td>
    {
      "activity_id" in activity && <Td>
        <Link href={`https://strava.com/activities/${activity.activity_id}`} target="_blank">
          <Text fontWeight={"bold"}>
            {activity.activity_name}
          </Text>
          <Box>
            {distance} miles | {elevationText} feet
          </Box>
        </Link>
      </Td>
    }
    {
      pairedEfforts.map((efforts) => {
        return <Td>
          {efforts.map((effort, index) => {
            if (!effort) {
              // Blank box of height 2 lines
              return <Box pb={index === 0 ? 3 : 0} height="2em"></Box>
            }
            return <>
              <Link target="_blank" href={`https://strava.com/activities/${effort.activity_id.toString()}/segments/${effort.effort_id.toString()}`}>
                <Text maxWidth={"100px"} isTruncated fontWeight={"bold"}>{segmentSummaries && segmentSummaries[effort.segment_id] ? segmentSummaries[effort.segment_id].name : "????"}</Text>
                <Box pb={index === 0 ? 3 : 0}>
                  {ElapsedDurationText(effort.elapsed_time, false)} @ {effort.device_watts ? Math.floor(effort.average_watts) + "w" : "--"}
                </Box>
              </Link >
            </>
          })}
        </Td>
      })
    }
  </Tr >
}