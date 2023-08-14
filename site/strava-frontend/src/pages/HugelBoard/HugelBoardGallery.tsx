import { FC, PropsWithChildren } from "react";
import { HugelBoardProps } from "./HugelBoard";
import {
  Flex,
  Grid,
  Box,
  Avatar,
  Badge,
  Text,
  useColorModeValue,
  BoxProps,
  Hide,
  VisuallyHidden,
  GridItem,
  FlexProps,
} from "@chakra-ui/react";
import {
  HugelLeaderBoardActivity,
  SuperHugelLeaderBoardActivity,
} from "../../api/typesGenerated";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import {
  DistanceToLocal,
  DistanceToLocalElevation,
} from "../../lib/Distance/Distance";
import { CalculateActivity } from "./CalcActivity";
import { ResponsiveCard } from "../../components/ResponsiveCard/ResponsiveCard";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { CardStat } from "../../components/CardStat/CardStat";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Loading } from "../../components/Loading/Loading";

export const HugelBoardGallery: FC<HugelBoardProps> = ({
  data,
  error,
  isLoading,
  isFetched,
}) => {
  if (isLoading) {
    return <Loading />;
  }
  return (
    <>
      <Grid
        gridTemplateColumns={{
          base: "repeat(1, 1fr)",
          // md: "repeat(3, 1fr)",
        }}
        //   templateAreas={`"first first first first first"
        // "" "second "third" "" ""`}
        alignItems={"center"}
        maxW={"1100px"}
      >
        <GridItem colSpan={3}>
          <GalleryCard activity={data?.activities[0]} position={1} />
        </GridItem>
        <GridItem colSpan={3}>
          <Flex
            flexDirection={{
              base: "column",
              lg: "row",
            }}
            width={"100%"}
            // justifyItems={"center"}
            alignItems={"center"}
            justifyContent={"center"}
          >
            <GalleryCard activity={data?.activities[1]} position={2} />
            <GalleryCard activity={data?.activities[2]} position={3} />
          </Flex>
        </GridItem>
        {data?.activities.slice(3).map((activity, index) => (
          <GridItem
            key={`activity-${index}`}
            colSpan={{
              base: 3,
              lg: 1,
            }}
          >
            <GalleryCard activity={activity} position={index + 4} />
          </GridItem>
        ))}
      </Grid>
    </>
  );
};

const SkeletonGalleryCard: React.FC<{
  position: number;
}> = ({}) => {
  return <Box w="100%" maxW="350px" h="300px" bg="transparent" />;
};

const GalleryCard: React.FC<{
  activity?: HugelLeaderBoardActivity | SuperHugelLeaderBoardActivity;
  position: number;
}> = ({ activity, position }) => {
  if (!activity) {
    return <GalleryCardBox display="none" />;
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
    numActivities,
  } = CalculateActivity(activity);

  const { firstname, lastname, profile_pic_link, username, hugel_count } =
    activity.athlete;

  const isSuper = !("activity_name" in activity);

  return (
    <GalleryCardBox>
      <Flex justifyContent={"space-between"}>
        <Flex p={3}>
          <AthleteAvatar
            firstname={firstname}
            lastname={lastname}
            athlete_id={activity.athlete_id}
            profile_pic_link={profile_pic_link}
            username={username}
            hugel_count={hugel_count}
            size="lg"
            styleProps={{
              mr: 3,
            }}
          />
          <Box>
            <Text fontWeight="bold" textAlign="left">
              {firstname} {lastname}
            </Text>
            <Text
              fontSize="sm"
              fontFamily={"monospace"}
              opacity={0.6}
              textAlign="left"
            >
              {dateText || "All time"}
            </Text>
            <Text
              isTruncated
              pl={"20px"}
              maxW="180px"
              textAlign={"right"}
              // noOfLines={1}
              fontWeight={"bold"}
              pt="10px"
            >
              {"activity_name" in activity
                ? activity.activity_name
                : `Across all rides`}
            </Text>
          </Box>
        </Flex>
        <Flex flexDirection={"column"} rowGap={2}>
          <Flex
            // Rank
            bg={"rgba(0,0,0,0.25)"}
            p={"1.5rem"}
            maxHeight={"2.5rem"}
            borderRadius={"0 10px 0 10px"}
            alignItems={"center"}
            justifyContent={"center"}
          >
            <Text fontWeight="bold" fontSize={"1.8rem"}>
              {position}
            </Text>
          </Flex>
          <StravaLink
            height={"45px"}
            width={"45px"}
            target="_blank"
            tooltip={
              "activity_id" in activity
                ? `View activity on strava`
                : `View athlete on strava`
            }
            href={
              "activity_id" in activity
                ? `"https://www.strava.com/activities/${activity.activity_id}"`
                : `https://www.strava.com/athletes/${activity.athlete_id}}`
            }
          />
        </Flex>
      </Flex>
      <Grid
        gridTemplateColumns={"repeat(3, 1fr)"}
        gap={3}
        p={4}
        padding={"20px"}
      >
        <GridItem colSpan={2} fontFamily="monospace">
          <Flex justifyContent={"center"} alignItems={"center"}>
            <FontAwesomeIcon color="#fc4c02" size="2x" icon={faClock} />
            <StatBox pl={"15px"} value={marginText} title={elapsedText} />
          </Flex>
        </GridItem>

        <StatBox title={distance.toFixed(1)} value={"mi"} />
        {isSuper ? (
          <StatBox title={numActivities.toString()} value={"rides"} />
        ) : (
          <StatBox title={totalElapsedText} value={"total"} />
        )}
        <StatBox
          title={showWatts ? avgWatts.toString() : "--"}
          value={"watts"}
        />
        {isSuper ? <></> : <StatBox title={elevationText} value={"ft"} />}
      </Grid>
    </GalleryCardBox>
  );
};

const StatBox: React.FC<
  FlexProps & {
    title: string;
    value: string;
  }
> = ({ title, value, ...props }) => {
  return (
    <CardStat
      flexDirection={"column-reverse"}
      value={title}
      title={value}
      fontSize={"1.3rem"}
      fontFamily="monospace"
      {...props}
    />
  );
};

const GalleryCardBox: React.FC<PropsWithChildren<BoxProps>> = ({
  children,
  ...props
}) => {
  return (
    <ResponsiveCard m={"10px"} w="100%" maxW="350px" h="270px" {...props}>
      {children}
    </ResponsiveCard>
  );
};

// {Array.from({ length: 20 }).map(e =>
//   <GalleryCard />
// )}
