import {
  FC,
  useRef,
  PropsWithChildren,
  Children,
  ReactNode,
  useState,
  useEffect,
  useMemo,
} from "react";
import {
  AthleteSummary,
  CompetitiveRoute,
  DetailedSegment,
  PersonalSegment,
  SegmentSummary,
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
  Card,
  chakra,
  position,
} from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getDetailedSegments, getErrorMessage, getRoute } from "../../api/rest";
import { NotFound } from "../404/404";
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
import L, { Control, LatLngBoundsExpression, Layer, map } from "leaflet";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import {
  faStar as fasStar,
  faCircleInfo,
  faLocationDot,
} from "@fortawesome/free-solid-svg-icons";
import { faStar as farStar } from "@fortawesome/free-regular-svg-icons";
import { useAuthenticated } from "../../contexts/Authenticated";
import { ElapsedDurationText, FormatDate } from "../HugelBoard/CalcActivity";
import { ResponsiveCard } from "../../components/ResponsiveCard/ResponsiveCard";
import { ConditionalLink } from "../../components/ConditionalLink/ConditionalLink";
import { CardStat } from "../../components/CardStat/CardStat";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { Loading } from "../../components/Loading/Loading";

const LongRouteColor = "#fc4c02";
const LiteRouteColor = "#FCFB03";

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

  const mainRoute = routeData?.routes[name || ""];
  const liteRoute = routeData?.routes["lite-" + name];

  const {
    data: segmentsData,
    error: segmentsError,
    isLoading: segmentsLoading,
    isFetched: segmentsFetched,
  } = useQuery({
    queryKey: ["hugel-segments", name],
    queryFn: () =>
      getDetailedSegments(mainRoute?.segments.map((e) => e.id) || []),
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
          <Flex flexDirection={"column"} pb="0.5em">
            <Heading fontSize={"4em"}>{mainRoute?.display_name}</Heading>
            <Text maxWidth={"1050px"} pt="1em">
              The Tour Das Hügel is an unsanctioned ride in Austin, Texas that
              showcases the challenging terrain around the city. To qualify as a
              Hügel, you must complete all of the following segments.
            </Text>
            <Flex
              direction={"row"}
              alignItems={"center"}
              justifyContent={"center"}
              gap={"2em"}
              pt="15px"
              pb="15px"
            >
              <Link href="https://www.strava.com/routes/3287834903479380952">
                <Flex
                  direction={"row"}
                  alignItems={"center"}
                  justifyContent={"center"}
                >
                  <chakra.img
                    height="1.5em"
                    width="1.5em"
                    src={"/logos/stravalogo.png"}
                    mr="5px"
                  />
                  <Text color="brand.stravaOrange">2024 Official Route</Text>
                </Flex>
              </Link>
              <Link href="https://www.strava.com/routes/3288659325860289186">
                <Flex
                  direction={"row"}
                  alignItems={"center"}
                  justifyContent={"center"}
                >
                  <chakra.img
                    height="1.5em"
                    width="1.5em"
                    src={"/logos/stravalogo.png"}
                    mr="5px"
                  />
                  <Text color={LiteRouteColor}>2024 Official Lite Route</Text>
                </Flex>
              </Link>
            </Flex>

            <Text>
              When: Saturday, November 9, 2024 7 a.m. Meetup and 7:15 a.m.
              Rollout
            </Text>
            <Text>
              Where:{" "}
              <Link
                target="_blank"
                href="https://www.google.com/maps/place/30%C2%B016'18.3%22N+97%C2%B046'22.3%22W/@30.271743,-97.772862,17z/data=!3m1!4b1!4m4!3m3!8m2!3d30.271743!4d-97.772862?entry=tts"
              >
                Under MoPac Expressway Off Stratford Dr. Next To Zilker Park
              </Link>
            </Text>
            <Link target="_blank" href="https://www.facebook.com/tourdashugel">
              <Flex
                direction={"row"}
                justifyContent={"center"}
                alignItems={"center"}
              >
                <chakra.img
                  src="/logos/facebook.svg"
                  width="2em"
                  height="2em"
                />

                <Text>@tourdashugel</Text>
              </Flex>
            </Link>
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
              liteRoute={liteRoute}
            />
            <TileLayer
              attribution='Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery &copy; <a href="https://www.mapbox.com/">Mapbox</a>'
              url={`https://api.mapbox.com/styles/v1/${mapboxUsername}/${mapboxStyleID}/tiles/256/{z}/{x}/{y}@2x?access_token=${mapboxAccessToken}`}
            />
          </MapContainer>
        </Flex>

        <SegmentCardContainer>
          {/* <Flex w="100%" flexDirection="column" p="2em"> */}
          {segmentsData
            .sort((a, b) =>
              computeSegmentName(
                a.detailed_segment.friendly_name,
                a.detailed_segment.name
              ).localeCompare(
                computeSegmentName(
                  b.detailed_segment.friendly_name,
                  b.detailed_segment.name
                )
              )
            )
            .map((segment) => (
              <Box key={segment.detailed_segment.id}>
                <SegmentCard
                  segment={segment}
                  setSelectedSegment={setSelectedSegment}
                  liteSegment={Boolean(
                    liteRoute?.segments.find(
                      (e) => e.id === segment.detailed_segment.id
                    )
                  )}
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
  liteRoute?: CompetitiveRoute;
}> = ({ selectedSegment, segments, outerBounds, liteRoute }) => {
  const mapRef = useMap();

  useEffect(() => {
    const legend = new L.Control({
      position: "topleft",
    });

    legend.onAdd = function () {
      const existing = document.getElementsByClassName("legend");
      if (existing && existing.length > 0) {
        existing[0].remove();
      }
      const container = L.DomUtil.create("div", "legend"); // Create the legend container

      // Style the container (optional)
      container.style.backgroundColor = "white";
      container.style.padding = "10px";
      container.style.border = "2px solid #ccc";
      container.style.color = "black";
      container.style.fontSize = "1.5em";

      // Define the HTML content for the legend
      container.innerHTML = `
        <h4><b>Color Legend</b></h4>
        <div><span style="background-color: ${LongRouteColor}; width: 10px; height: 10px; display: inline-block;"></span> Das Hügel</div>
        <div><span style="background-color: ${LiteRouteColor}; width: 10px; height: 10px; display: inline-block;"></span> Hügel Lite</div>
    `;
      legend.onRemove = function () {};

      return container;
    };
    legend.addTo(mapRef);
  }, [mapRef]);
  // const legend = new LegendControl();

  const segmentLineColor = (segmentID: string): string => {
    if (
      liteRoute &&
      liteRoute.segments.find((e) => e.id === segmentID) !== undefined
    ) {
      return LiteRouteColor;
    }
    return LongRouteColor;
  };

  const [polys, setPolys] = useState([] as [string, L.Polyline][]);
  useEffect(() => {
    const polys = segments.map((segment) => {
      const points = decode(segment.detailed_segment.map.polyline);
      const circleRadius = 5;
      // const popUp = <Popup>{segment.detailed_segment.name}</Popup>;
      const poly = L.polyline(points, {
        weight: 3,
        color: segmentLineColor(segment.detailed_segment.id),
      });
      const start = L.circleMarker(points[0], {
        radius: circleRadius,
        color: "green",
      });
      const end = L.circleMarker(points[points.length - 1], {
        radius: circleRadius,
        color: "red",
      });

      const group = L.featureGroup([poly, start, end]);
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
        poly.setStyle({ color: segmentLineColor(id), weight: 3 });
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

const computeSegmentName = (friendlyName: string, name: string): string => {
  return friendlyName !== "" ? friendlyName : name;
};

const SegmentCard: FC<{
  segment: PersonalSegment;
  setSelectedSegment: (id: string) => void;
  liteSegment?: boolean;
}> = ({ segment, setSelectedSegment, liteSegment }) => {
  const { authenticatedUser } = useAuthenticated();

  const bestActHref = segment.personal_best
    ? `https://www.strava.com/activities/${segment.personal_best.best_effort_activities_id}/segments/${segment.personal_best.best_effort_id}`
    : "";

  const segmentName = computeSegmentName(
    segment.detailed_segment.friendly_name,
    segment.detailed_segment.name
  );

  return (
    <ResponsiveCard height={"170px"} width={"350px"}>
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
            <Grid templateColumns="repeat(13, 1fr)" columnGap={2}>
              <GridItem
                fontSize={segmentName.length < 20 ? "1.2rem" : "1rem"}
                // fontSize={segment.detailed_segment.name.length < 20 ? "2em" : "1em"}
                textAlign={"center"}
                colSpan={10}
                color={liteSegment ? LiteRouteColor : LongRouteColor}
              >
                {segment.detailed_segment.friendly_name !== "" ? (
                  <Tooltip label={segment.detailed_segment.name}>
                    {segment.detailed_segment.friendly_name}
                  </Tooltip>
                ) : (
                  segmentName
                )}
              </GridItem>
              <GridItem colSpan={3}>
                <Flex columnGap={3}>
                  <Link
                    as="span"
                    onClick={() => {
                      setSelectedSegment(segment.detailed_segment.id);
                    }}
                    transition={"all .1s ease"}
                    _hover={{
                      transform: "scale(1.1)",
                    }}
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
                  <StravaLink
                    href={`https://strava.com/segments/${segment.detailed_segment.id}`}
                    target="_blank"
                    height={"34px"}
                    width={"34px"}
                    tooltip="View segment on Strava"
                  />
                </Flex>
              </GridItem>
            </Grid>
          </GridItem>

          {/* Now the stats */}
          <GridItem>
            <CardStat
              title={"Distance"}
              value={
                DistanceToLocal(segment.detailed_segment.distance).toFixed(0) +
                " mi"
              }
            />
          </GridItem>
          <GridItem>
            <CardStat
              title={"Avg Grade"}
              value={segment.detailed_segment.average_grade + "%"}
            />
          </GridItem>
          <GridItem>
            <CardStat
              title={"Elevation"}
              value={
                DistanceToLocalElevation(
                  segment.detailed_segment.total_elevation_gain
                ).toFixed(0) + " ft"
              }
            />
          </GridItem>

          <GridItem>
            <ConditionalLink href={bestActHref}>
              <CardStat
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
              <CardStat
                title="PR"
                value={
                  segment.personal_best
                    ? ElapsedDurationText(
                        segment.personal_best.best_effort_elapsed_time,
                        true,
                        false
                      )
                    : "--:--"
                }
              />
            </ConditionalLink>
          </GridItem>
          <GridItem textAlign={"center"}>
            {/* <StarredIcon starred={segment.starred} /> */}
            <Flex
              justifyContent="center"
              alignItems={"center"}
              height={"100%"}
            ></Flex>
          </GridItem>
        </Grid>
      </Box>
    </ResponsiveCard>
  );
};

const StarredIcon: FC<{ starred?: boolean }> = ({ starred }) => {
  const { authenticatedUser } = useAuthenticated();
  const starIcon =
    authenticatedUser === undefined
      ? faCircleInfo
      : starred
      ? fasStar
      : farStar;

  const starColor = authenticatedUser === undefined ? "#709df8" : "#fcaf02";

  const starTooltip =
    authenticatedUser === undefined
      ? "Connect with strava to see if you have this segment starred"
      : starred
      ? "You have this segment starred"
      : "Segment is not starred";

  return (
    <Box
      cursor="help"
      _hover={{
        opacity: 1,
      }}
      opacity={0.7}
    >
      <Tooltip label={starTooltip} aria-label="Starred Segment Tooltip">
        <FontAwesomeIcon
          icon={starIcon}
          size="2x"
          style={{
            color: starColor,
          }}
        />
      </Tooltip>
    </Box>
  );
};
