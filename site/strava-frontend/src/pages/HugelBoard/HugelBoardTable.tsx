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
import { HugelLeaderBoard, HugelLeaderBoardActivity } from "../../api/typesGenerated"
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
            <Th>Elapsed Time</Th>
            <Th>Activity</Th>
            {
              hugelSegments && hugelSegments.segments[0] && SortSegments(hugelSegments.segments).map((segment) => {
                return <Th>
                  {segment.name}
                </Th>
              })
            }
          </Tr>
        </Thead>
        <Tbody>
          {data && data.activities?.map((activity) => {
            return <HugelBoardTableRow key={activity.activity_id} {...activity} />
          })
          }
        </Tbody>
      </Table>
    </TableContainer>
  </>
}



export const HugelBoardTableRow: FC<PropsWithChildren<HugelLeaderBoardActivity>> = (activity) => {
  const {
    firstname,
    lastname,
    profile_pic_link,
    username,
  } = activity.athlete

  const {
    dateText,
    elapsedText,
    totalElapsedText,
    elevationText,
    distance,
    showWatts,
    avgWatts
  } = CalculateActivity(activity)
  return <Tr key={`row-${activity.activity_id}`}>
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
    <Td>
      {elapsedText}
    </Td>
    <Td>
      <Box>
        {activity.activity_name}
      </Box>
      <Box>
        {distance} miles | {elevationText} feet
      </Box>
    </Td>
    {
      SortEfforts(activity.efforts).map((effort) => {
        return <Td>
          <Box>
            {ElapsedDurationText(false, effort.elapsed_time)}
          </Box>
          <Box>
            {effort.device_watts ? effort.average_watts + "w" : "--"}
          </Box>
        </Td>
      })
    }
  </Tr>
}