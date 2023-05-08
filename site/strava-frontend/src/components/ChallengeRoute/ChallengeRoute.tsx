import { FC, useRef } from "react";
import { AthleteSummary, DetailedSegment } from "../../api/typesGenerated";
import { Image, Avatar, AvatarBadge, AvatarProps, Box, Flex, Heading, Text, Link } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getDetailedSegments, getErrorMessage, getHugelLeaderBoard, getRoute } from "../../api/rest";
import { NotFound } from "../../pages/404/404";
import { useQuery } from "@tanstack/react-query";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  CircleMarker,
  Circle,
  Polyline,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import "leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css";
import "leaflet-defaulticon-compatibility";
import { decode } from "@mapbox/polyline";
import { faCirclePlay } from '@fortawesome/free-solid-svg-icons'
import { DistanceToLocal, DistanceToLocalElevation } from "../../lib/Distance/Distance";

export const ChallengeRoute: FC<{

}> = ({  }) => {
  const {name} = useParams()
  const mapRef = useRef(null);
  
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


  const mapboxAccessToken = "pk.eyJ1IjoiZW15cmsiLCJhIjoiY2wweW93ZnYzMGp0OTNvbzN5a2VvNWVldyJ9.QyM0MUn75YqHqMUvMlMaag"
  const mapboxStyleID = "clhebqdem028w01p85pnzcsch"
  const mapboxUsername = "emyrk"

  //mapbox://styles/emyrk/clhe4rd8l027g01pa3bdh5u4v

  // https://www.paigeniedringhaus.com/blog/render-multiple-colored-lines-on-a-react-map-with-polylines
  // pk.eyJ1IjoiZW15cmsiLCJhIjoiY2xoZTR6YjAxMWh0ODNqbzc5NjRxdzBxbCJ9._SlRHXQG5-DqZTucbZUagA
  // https://api.mapbox.com/styles/v1/emyrk/clhe4rd8l027g01pa3bdh5u4v.html?title=view&access_token=pk.eyJ1IjoiZW15cmsiLCJhIjoiY2xoZTR6YjAxMWh0ODNqbzc5NjRxdzBxbCJ9._SlRHXQG5-DqZTucbZUagA&zoomwheel=true&fresh=true#7.5/42.2/9.1
  return <>
  <Flex w='100%' maxW={'7xl'} m={'1rem auto 0'} flexDirection="column">
    <Flex w="100%" justifyContent={"center"} alignItems={"center"} textAlign="center">
      <Flex flexDirection={"column"}>
        <Heading fontSize={"6em"}>{routeData.display_name}</Heading>
        <Text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum rhoncus in arcu eu luctus. Donec vel risus suscipit, dignissim quam sit amet, blandit nulla. Donec molestie tincidunt odio, ut sodales mi porttitor ut. Cras fermentum vestibulum viverra. Mauris blandit gravida nibh, sit amet viverra risus convallis eget. Donec efficitur commodo luctus. Vestibulum elementum odio arcu, vitae fermentum urna feugiat eu.</Text>
      </Flex>
    </Flex>

    <Flex w="100%" flexDirection="column" alignItems={"center"} textAlign={"center"} pt={"2em"}>
      <MapContainer ref={mapRef} style={{ zIndex: 0, borderRadius:"10px", height: "650px", width: "80%" }} center={[30.349426, -97.774007]} zoom={12}>
        <TileLayer 
        attribution="Map data &copy; <a href=&quot;https://www.openstreetmap.org/&quot;>OpenStreetMap</a> contributors, <a href=&quot;https://creativecommons.org/licenses/by-sa/2.0/&quot;>CC-BY-SA</a>, Imagery &copy; <a href=&quot;https://www.mapbox.com/&quot;>Mapbox</a>"
        url={`https://api.mapbox.com/styles/v1/${mapboxUsername}/${mapboxStyleID}/tiles/256/{z}/{x}/{y}@2x?access_token=${mapboxAccessToken}`}
        // url={`https://api.mapbox.com/styles/v1/mapbox/streets-v11/tiles/256/{z}/{x}/{y}@2x?access_token=pk.eyJ1IjoiZW15cmsiLCJhIjoiY2wweW93ZnYzMGp0OTNvbzN5a2VvNWVldyJ9.QyM0MUn75YqHqMUvMlMaag`}
        />
        {segmentsData.map(segment =>{
          const points = decode(segment.map.polyline)
          const circleRadius = 5
          const popUp = <Popup>
              {segment.name}
          </Popup>
          return <Box key={segment.id}>
            <Polyline weight={3} key={segment.id} pathOptions={{ color: "#fc4c02" }} positions={points}>
              {popUp}
            </Polyline>
            <CircleMarker center={points[0]} radius={circleRadius} color="green">
              {popUp}
            </CircleMarker>
            <CircleMarker center={points[points.length-1]} radius={circleRadius} color="red">
              {popUp}
            </CircleMarker>
          </Box>
        })} 
      </MapContainer>
    </Flex>
    <Flex w="100%" flexDirection="column" p="2em">
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
  //  bg="rgba(248,248,248,0.3)"
  return <Flex 
    _hover={{ bg: "rgba(248,248,248,0.5)" }}
    alignItems={"center"} textAlign={"center"}
    m={"0.5em"} bg="rgba(248,248,248,0.3)" width="100%" borderRadius={"4px"} p={0}
  >
    {/* TODO: Add a checkmark if the athlete has done the segment or not */}
    <Link href={`https://strava.com/segments/${segment.id}`} target="_blank" height={"4em"} p={"10px"}>
      <Image src={"/logos/stravalogo.png"} height="100%"/>
    </Link>
    <Image src={segment.elevation_profile} p={"10px"}/>
    <Text fontSize={"1em"} fontWeight={"bold"} p={"10px"}>{segment.name}</Text>
    <Text>{` Distance: ${Math.floor(DistanceToLocal(segment.distance)*10)/10} mi | Avg Grade: ${segment.average_grade}% | Elevation Gain: ${Math.floor(DistanceToLocalElevation(segment.total_elevation_gain))} ft`}</Text>

    </Flex>
}

export const Loading: FC = () => {
  return <Text>Loading...</Text>
}



{/* <iframe height='405' width='590' frameborder='0' allowtransparency='true' scrolling='no' src='https://www.strava.com/segments/7041089/embed'></iframe> */}