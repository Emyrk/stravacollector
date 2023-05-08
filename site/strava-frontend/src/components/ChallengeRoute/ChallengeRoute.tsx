import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps, Box, Flex, Heading, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getErrorMessage, getHugelLeaderBoard, getRoute } from "../../api/rest";
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
  
  if(!name) {
    return <NotFound />
  }


  if(routeLoading) {
    return <Loading />
  }

  if(routeError) {
    return <Text>{getErrorMessage(routeError, "route failed to load")}</Text>
  }

  if(!routeData) {
    return <NotFound />
  }

  console.log(routeData)  
  return <>
  <Flex w='100%' maxW={'7xl'} m={'1rem auto 0'} flexDirection="column">
    <Flex w="100%" justifyContent={"center"} alignItems={"center"} textAlign="center">
      <Flex flexDirection={"column"}>
        <Heading fontSize={"8em"}>{routeData.display_name}</Heading>
        <Text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum rhoncus in arcu eu luctus. Donec vel risus suscipit, dignissim quam sit amet, blandit nulla. Donec molestie tincidunt odio, ut sodales mi porttitor ut. Cras fermentum vestibulum viverra. Mauris blandit gravida nibh, sit amet viverra risus convallis eget. Donec efficitur commodo luctus. Vestibulum elementum odio arcu, vitae fermentum urna feugiat eu.</Text>
      </Flex>
    </Flex>

    <Flex w="50%" flexDirection="column" >
        {Array.from({ length: 20 }).map(e =>
          <SegmentCard />
        )} 
    </Flex>
  </Flex>
    
  </>
}

const SegmentCard: FC<{}> = ({}) => {
  return <Box m={1} bg="lightblue" width="100%" height="2em">
    asd
    </Box>
}

export const Loading: FC = () => {
  return <Text>Loading...</Text>
}



{/* <iframe height='405' width='590' frameborder='0' allowtransparency='true' scrolling='no' src='https://www.strava.com/segments/7041089/embed'></iframe> */}