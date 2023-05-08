import { FC } from "react";
import { AthleteSummary, DetailedSegment } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps, Box, Flex, Heading, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getDetailedSegments, getErrorMessage, getHugelLeaderBoard, getRoute } from "../../api/rest";
import { NotFound } from "../../pages/404/404";
import { useQuery } from "@tanstack/react-query";

export const ChallengeRoute: FC<{

}> = ({  }) => {
  const {name} = useParams()
  
  const queryKey = ["hugel-leaderboard", name]
  const {
    data: routeData,
    error: routeError,
    isLoading: routeLoading,
    isFetched: routeFetched,
  } = useQuery({
    queryKey,
    queryFn: () => getRoute(name || ""),
    enabled: !!name,
  })
  
  const {
    data: segmentsData,
    error: segmentsError,
    isLoading: segmentsLoading,
    isFetched: segmentsFetched,
  } = useQuery({
    queryKey: ["hugel-segments", name],
    queryFn: () => getDetailedSegments(routeData?.segments.map(e => e.id) || []),
    enabled: !!name && !!routeData,
  })

  if(!name) {
    return <NotFound />
  }


  if(routeLoading || segmentsLoading) {
    return <Loading />
  }

  if(routeError || segmentsError) {
    return <Text>{getErrorMessage(routeError, "route failed to load")}</Text>
  }

  if(!routeData || !segmentsData) {
    return <NotFound />
  }

  return <>
  <Flex w='100%' maxW={'7xl'} m={'1rem auto 0'} flexDirection="column">
    <Flex w="100%" justifyContent={"center"} alignItems={"center"} textAlign="center">
      <Flex flexDirection={"column"}>
        <Heading fontSize={"8em"}>{routeData.display_name}</Heading>
        <Text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum rhoncus in arcu eu luctus. Donec vel risus suscipit, dignissim quam sit amet, blandit nulla. Donec molestie tincidunt odio, ut sodales mi porttitor ut. Cras fermentum vestibulum viverra. Mauris blandit gravida nibh, sit amet viverra risus convallis eget. Donec efficitur commodo luctus. Vestibulum elementum odio arcu, vitae fermentum urna feugiat eu.</Text>
      </Flex>
    </Flex>

    <Flex w="100%" flexDirection="column" >
      Map here
    </Flex>
    <Flex w="100%" flexDirection="column" >
        {segmentsData.map(segment =>
          <SegmentCard key={segment.id} segment={segment}/>
        )} 
    </Flex>
  </Flex>
    
  </>
}

const SegmentCard: FC<{
  segment: DetailedSegment
}> = ({segment}) => {
  return <Box m={1} bg="lightblue" width="100%" height="2em">
    {segment.name}
    </Box>
}

export const Loading: FC = () => {
  return <Text>Loading...</Text>
}



{/* <iframe height='405' width='590' frameborder='0' allowtransparency='true' scrolling='no' src='https://www.strava.com/segments/7041089/embed'></iframe> */}