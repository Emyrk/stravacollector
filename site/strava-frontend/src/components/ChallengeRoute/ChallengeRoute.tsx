import {
  FC,
  useRef,
  PropsWithChildren,
  Children,
  ReactNode,
  useState,
  useEffect,
} from "react";
import {
  AthleteSummary,
  DetailedSegment,
  PersonalSegment,
} from "../../api/typesGenerated";
import {
  Image,
  Avatar,
  AvatarBadge,
  AvatarProps,
  Box,
  Flex,
  Heading,
  Text,
  Link,
  Grid,
  LinkBox,
  Table,
  SimpleGrid,
  GridItem,
  Tooltip,
  background,
  Button,
} from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import {
  getDetailedSegments,
  getErrorMessage,
  getHugelLeaderBoard,
  getRoute,
} from "../../api/rest";
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
  useMap,
  FeatureGroup,
  Rectangle,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import "leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css";
import "leaflet-defaulticon-compatibility";
import { decode } from "@mapbox/polyline";
import {
  DistanceToLocal,
  DistanceToLocalElevation,
} from "../../lib/Distance/Distance";
import L, { LatLngBoundsExpression, Layer } from "leaflet";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import {
  faStar as fasStar,
  faCircleInfo,
  faLocationDot,
} from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";
import { useAuthenticated } from "../../contexts/Authenticated";
import {
  ElapsedDurationText,
  FormatDate,
} from "../../pages/HugelBoard/CalcActivity";

export const ChallengeRoute: FC<{}> = ({}) => {
  const { name } = useParams();
  const mapRef = useRef(null);
  const [selectedSegment, setSelectedSegment] = useState<string>("");

  const queryKey = ["hugel-leaderboard", name];
  const {
    data: routeData,
    error: routeError,
    isLoading: routeLoading,
    isFetched: routeFetched,
  } = useQuery({
    queryKey,
    queryFn: () => getRoute(name || ""),
    enabled: !!name,
  });

  const {
    data: segmentsData,
    error: segmentsError,
    isLoading: segmentsLoading,
    isFetched: segmentsFetched,
  } = useQuery({
    queryKey: ["hugel-segments", name],
    queryFn: () =>
      getDetailedSegments(routeData?.segments.map((e) => e.id) || []),
    enabled: !!name && !!routeData,
  });

  if (!name) {
    return <NotFound />;
  }

  if (routeLoading || segmentsLoading) {
    return <Loading />;
  }

  if (routeError || segmentsError) {
    return <Text>{getErrorMessage(routeError, "route failed to load")}</Text>;
  }

  if (!routeData || !segmentsData) {
    return <NotFound />;
  }

  const mapboxAccessToken =
    "pk.eyJ1IjoiZW15cmsiLCJhIjoiY2wweW93ZnYzMGp0OTNvbzN5a2VvNWVldyJ9.QyM0MUn75YqHqMUvMlMaag";
  const mapboxStyleID = "clhebqdem028w01p85pnzcsch";
  const mapboxUsername = "emyrk";

  // Calculate the center of the route
  const sum = segmentsData.reduce(
    (acc, segment) => {
      const points = decode(segment.detailed_segment.map.polyline);
      const sum = points.reduce(
        (acc, point) => {
          return [acc[0] + point[0], acc[1] + point[1]];
        },
        [0, 0]
      );
      const center = [sum[0] / points.length, sum[1] / points.length];

      return [acc[0] + center[0], acc[1] + center[1]];
    },
    [0, 0]
  );
  const center: [number, number] = [
    sum[0] / segmentsData.length,
    sum[1] / segmentsData.length,
  ];
  // Calculate the bounds of the route
  const strictBounds: [[number, number], [number, number]] =
    segmentsData.reduce(
      (acc, segment) => {
        const points = decode(segment.detailed_segment.map.polyline);

        const min = points.reduce(
          (acc, point) => {
            return [Math.min(acc[0], point[0]), Math.min(acc[1], point[1])];
          },
          [360, 360]
        );
        const max = points.reduce(
          (acc, point) => {
            return [Math.max(acc[0], point[0]), Math.max(acc[1], point[1])];
          },
          [-360, -360]
        );

        return [
          [Math.min(acc[0][0], min[0]), Math.min(acc[0][1], min[1])],
          [Math.max(acc[1][0], max[0]), Math.max(acc[1][1], max[1])],
        ];
      },
      [
        [360, 360],
        [-360, -360],
      ]
    );

  // Add some tolerance to the bounds for easier viewing
  const bounds: [[number, number], [number, number]] = [
    [strictBounds[0][0] - 0.1, strictBounds[0][1] - 0.1],
    [strictBounds[1][0] + 0.1, strictBounds[1][1] + 0.1],
  ];

  //mapbox://styles/emyrk/clhe4rd8l027g01pa3bdh5u4v
  // https://www.paigeniedringhaus.com/blog/render-multiple-colored-lines-on-a-react-map-with-polylines
  // pk.eyJ1IjoiZW15cmsiLCJhIjoiY2xoZTR6YjAxMWh0ODNqbzc5NjRxdzBxbCJ9._SlRHXQG5-DqZTucbZUagA
  // https://api.mapbox.com/styles/v1/emyrk/clhe4rd8l027g01pa3bdh5u4v.html?title=view&access_token=pk.eyJ1IjoiZW15cmsiLCJhIjoiY2xoZTR6YjAxMWh0ODNqbzc5NjRxdzBxbCJ9._SlRHXQG5-DqZTucbZUagA&zoomwheel=true&fresh=true#7.5/42.2/9.1
  return (
    <>
      <Flex
        w="100%"
        maxW={"7xl"}
        m={"1rem auto 0"}
        flexDirection="column"
        alignItems={"center"}
      >
        <Flex
          w="100%"
          justifyContent={"center"}
          alignItems={"center"}
          textAlign="center"
        >
          <Flex flexDirection={"column"}>
            <Heading fontSize={"6em"}>{routeData.display_name}</Heading>
            <Text>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit.
              Vestibulum rhoncus in arcu eu luctus. Donec vel risus suscipit,
              dignissim quam sit amet, blandit nulla. Donec molestie tincidunt
              odio, ut sodales mi porttitor ut. Cras fermentum vestibulum
              viverra. Mauris blandit gravida nibh, sit amet viverra risus
              convallis eget. Donec efficitur commodo luctus. Vestibulum
              elementum odio arcu, vitae fermentum urna feugiat eu.
            </Text>
          </Flex>
        </Flex>

        <Flex
          w="100%"
          flexDirection="column"
          alignItems={"center"}
          textAlign={"center"}
          pt={"2em"}
        >
          <MapContainer
            ref={mapRef}
            style={{
              zIndex: 0,
              borderRadius: "10px",
              height: "650px",
              width: "80%",
            }}
            center={center}
            zoom={12}
            maxBounds={bounds}
          >
            {/* <Rectangle bounds={bounds}></Rectangle> */}
            <MapController
              segments={segmentsData}
              selectedSegment={selectedSegment}
              outerBounds={strictBounds}
            />
            <TileLayer
              attribution='Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery &copy; <a href="https://www.mapbox.com/">Mapbox</a>'
              url={`https://api.mapbox.com/styles/v1/${mapboxUsername}/${mapboxStyleID}/tiles/256/{z}/{x}/{y}@2x?access_token=${mapboxAccessToken}`}
            />
          </MapContainer>
        </Flex>

        <SegmentCardContainer>
          {/* <Flex w="100%" flexDirection="column" p="2em"> */}
          {segmentsData.map((segment) => (
            <Box key={segment.detailed_segment.id}>
              <SegmentCard
                segment={segment}
                setSelectedSegment={setSelectedSegment}
              />
            </Box>
          ))}
          {/* </Flex> */}
        </SegmentCardContainer>
      </Flex>
    </>
  );
};

const MapController: FC<{
  selectedSegment: string;
  segments: PersonalSegment[];
  outerBounds: LatLngBoundsExpression;
}> = ({ selectedSegment, segments, outerBounds }) => {
  const mapRef = useMap();

  const [polys, setPolys] = useState([] as [string, L.Polyline][]);
  useEffect(() => {
    const polys = segments.map((segment) => {
      const points = decode(segment.detailed_segment.map.polyline);
      const circleRadius = 5;
      // const popUp = <Popup>{segment.detailed_segment.name}</Popup>;
      const poly = L.polyline(points, {
        weight: 3,
        color: "#fc4c02",
      });
      const start = L.circleMarker(points[0], {
        radius: circleRadius,
        color: "green",
      });
      const end = L.circleMarker(points[points.length - 1], {
        radius: circleRadius,
        color: "red",
      });

      const group = L.featureGroup([poly, start, end]).on("click", () => {
        console.log("Clicked", segment.detailed_segment.id);
      });
      // https://www.wrld3d.com/wrld.js/latest/docs/leaflet/L.Popup/
      group.bindTooltip(segment.detailed_segment.name);

      mapRef.addLayer(group);
      return [segment.detailed_segment.id, poly] as [string, L.Polyline];
    });
    setPolys(polys);
  }, [mapRef, segments]);

  useEffect(() => {
    if (selectedSegment === "") {
      mapRef.fitBounds(outerBounds);
    }
    polys.forEach(([id, poly]) => {
      if (id === selectedSegment) {
        poly.setStyle({ color: "#0215fc", weight: 5 });
        mapRef.panTo(poly.getBounds().getCenter());
        mapRef.fitBounds(poly.getBounds());
      } else {
        poly.setStyle({ color: "#fc4c02", weight: 3 });
      }
    });
  }, [mapRef, polys, selectedSegment, outerBounds]);

  return <></>;
};

const SegmentCardContainer: FC<PropsWithChildren<{}>> = ({ children }) => {
  return (
    <Grid
      pt={"2em"}
      gridTemplateColumns={{
        base: "repeat(1, 1fr)",
        md: "repeat(2, 1fr)",
        lg: "repeat(auto-fit, minmax(300px, 1fr))",
      }}
      rowGap={4}
      columnGap={6}
      maxWidth={"1050px"}
    >
      {children}
    </Grid>
  );
};

const SegmentCard: FC<{
  segment: PersonalSegment;
  setSelectedSegment: (id: string) => void;
}> = ({ segment, setSelectedSegment }) => {
  const { authenticatedUser } = useAuthenticated();
  const starIcon =
    authenticatedUser === undefined
      ? faCircleInfo
      : segment.starred
      ? fasStar
      : farStar;

  const starColor = authenticatedUser === undefined ? "#709df8" : "#fcaf02";

  const starTooltip =
    authenticatedUser === undefined
      ? "Connect with strava to see if you have this segment starred"
      : segment.starred
      ? "You have this segment starred"
      : "Segment is not starred";

  const bestActHref = segment.personal_best
    ? `https://www.strava.com/activities/${segment.personal_best?.best_effort_activities_id}`
    : "";

  return (
    <Box
      // bgColor={"#272c35"}
      bgColor={"#3b3f48"}
      height={"170px"}
      width={"350px"}
      borderRadius={"10px"}
      boxShadow={"rgb(20, 20, 20) 0px 3px 6px"}
      _hover={{
        bg: "rgb(43, 49, 59)",
        boxShadow: "rgb(20, 20, 20) 0px 5px 10px",
        marginTop: "-3px",
        marginBottom: "-3px",
      }}
    >
      <Box p="10px">
        <Grid
          // templateAreas={`"header header"
          //         "nav main"
          //         "nav footer"`}
          gridTemplateRows={"repeat(3, 1fr)"}
          gridTemplateColumns={"repeat(3, 1fr)"}
          gap="1"
          rowGap={2}
          maxH={"4rem"}
        >
          {/*  The card header */}
          <GridItem colSpan={3}>
            <Grid templateColumns="repeat(10, 1fr)">
              <GridItem
                fontSize={
                  segment.detailed_segment.name.length < 20 ? "1.2rem" : "1rem"
                }
                // fontSize={segment.detailed_segment.name.length < 20 ? "2em" : "1em"}
                textAlign={"center"}
                colSpan={8}
              >
                {segment.detailed_segment.name}
              </GridItem>
              <GridItem colSpan={2}>
                <Flex alignItems={"center"} columnGap={3}>
                  <Link
                    as="span"
                    onClick={() =>
                      setSelectedSegment(segment.detailed_segment.id)
                    }
                  >
                    <Tooltip
                      label="Zoom to segment"
                      aria-label="Segment view tooltip"
                    >
                      <FontAwesomeIcon
                        style={{ color: "#fc4c02" }}
                        icon={faLocationDot}
                        size="2x"
                      />
                    </Tooltip>
                  </Link>
                  <Tooltip
                    label="Link to segment"
                    aria-label="Strava logo tooltip"
                  >
                    <Link
                      href={`https://strava.com/segments/${segment.detailed_segment.id}`}
                      target="_blank"
                    >
                      <Image
                        src={"/logos/stravalogo.png"}
                        height={"30px"}
                        width={"auto"}
                      />
                    </Link>
                  </Tooltip>
                </Flex>
              </GridItem>
            </Grid>
          </GridItem>

          {/* Now the stats */}
          <GridItem>
            <SegmentStat
              title={"Distance"}
              value={
                Math.floor(
                  DistanceToLocal(segment.detailed_segment.distance) * 10
                ) /
                  10 +
                "mi"
              }
            />
          </GridItem>
          <GridItem>
            <SegmentStat
              title={"Avg Grade"}
              value={segment.detailed_segment.average_grade + "%"}
            />
          </GridItem>
          <GridItem>
            <SegmentStat
              title={"Avg Grade"}
              value={
                Math.floor(
                  DistanceToLocalElevation(
                    segment.detailed_segment.total_elevation_gain
                  )
                ) + "ft"
              }
            />
          </GridItem>

          <GridItem>
            <ConditionalLink href={bestActHref}>
              <SegmentStat
                title="PR Activity"
                value={
                  segment.personal_best
                    ? FormatDate(
                        segment.personal_best.best_effort_start_date,
                        true
                      )
                    : "--/--/----"
                }
              />
            </ConditionalLink>
          </GridItem>
          <GridItem>
            <ConditionalLink href={bestActHref}>
              <SegmentStat
                title="PR"
                value={
                  segment.personal_best
                    ? ElapsedDurationText(
                        segment.personal_best.best_effort_elapsed_time,
                        false
                      )
                    : "--:--"
                }
              />
            </ConditionalLink>
          </GridItem>
          <GridItem textAlign={"center"}>
            <Flex justifyContent="center" alignItems={"center"} height={"100%"}>
              <Tooltip label={starTooltip} aria-label="Starred Segment Tooltip">
                <FontAwesomeIcon
                  icon={starIcon}
                  size="2x"
                  style={{ color: starColor }}
                />
              </Tooltip>
            </Flex>
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

const ConditionalLink: FC<{
  href: string;
  children: ReactNode;
}> = ({ href, children }) => {
  if (href !== "") {
    return (
      <Link target="_blank" href={href}>
        {children}
      </Link>
    );
  }
  return <>{children}</>;
};

const SegmentStat: FC<{
  title: string;
  value: string;
}> = ({ title, value }) => {
  return (
    <Flex
      flexDirection={"column"}
      alignItems={"center"}
      justifyContent={"center"}
    >
      <Text color="#a7afbe" fontSize={"13px"}>
        {title}
      </Text>
      <Text fontSize={"16px"}>{value}</Text>
    </Flex>
  );
};

export const Loading: FC = () => {
  return <Text>Loading...</Text>;
};

{
  /* <iframe height='405' width='590' frameborder='0' allowtransparency='true' scrolling='no' src='https://www.strava.com/segments/7041089/embed'></iframe> */
}
