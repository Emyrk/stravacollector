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
import { getErrorDetail, getErrorMessage, getHugelLeaderBoard } from "../../api/rest"
import { useMutation, useQuery } from "@tanstack/react-query";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { HugelBoardProps } from "./HugelBoard";
import { HugelLeaderBoard, HugelLeaderBoardActivity } from "../../api/typesGenerated"
import { CalculateActivity, ElapsedDurationText } from "./CalcActivity";
import React from "react";

export const HugelBoardTable: FC<HugelBoardProps> = ({
  data, error, isLoading, isFetched
}) => {

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
            {/* {Segment IDs
              data && data.activities[0] && data.activities[0].efforts.map((effort) => {
                return <Th>
                    {effort.segment_id}
                </Th>
              })
            } */}
          </Tr>
        </Thead>
        <Tbody>
        {data && data.activities?.map((activity) => {
           return  <HugelBoardTableRow {...activity}/>
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
  return <Tr>
    <Td>
    <Flex p={3} alignItems={'center'}>
      <Text fontSize={30} pr={5}>
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
      activity.efforts.map((effort) => {
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