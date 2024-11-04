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
  Stack,
  AvatarGroup,
} from "@chakra-ui/react";
import {
  HugelLeaderBoardActivity,
  SuperHugelLeaderBoardActivity,
  SuperlativeEntry,
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
import { Superlative } from "../../components/Superlative/Superlative";
import { ErrorBox } from "../../components/ErrorBox/ErrorBox";

export const HugelBoardGallery: FC<HugelBoardProps> = ({
  data,
  error,
  isLoading,
  isFetched,
  disableSuperlatives,
}) => {
  if (isLoading) {
    return <Loading />;
  }
  if (error) {
    return <ErrorBox error={error.toString()} />;
  }
  if (data && data.activities.length === 0) {
    return <>No rides completed for this year.</>;
  }

  var superlatives: Record<string, Record<string, SuperlativeEntry<any>>> = {};
  if (!disableSuperlatives && data && "superlatives" in data) {
    for (const [key, value] of Object.entries(data.superlatives)) {
      const entry = value as SuperlativeEntry<any>;
      if (!superlatives[entry.activity_id]) {
        superlatives[entry.activity_id] = {};
      }

      superlatives[entry.activity_id][key] = entry;
    }
  }

  var athSuperlatives = (
    activity?: SuperHugelLeaderBoardActivity | HugelLeaderBoardActivity
  ): Record<string, SuperlativeEntry<any>> | undefined => {
    // let sups: Record<string, SuperlativeEntry<any>> | undefined = undefined;
    if (superlatives && activity && "activity_id" in activity) {
      return superlatives[activity.activity_id as string];
    }
    return undefined;
  };

  return (
    <>
      <Grid
        gridTemplateColumns={
          {
            // base: "repeat(1, 1fr)",
            // md: "repeat(3, 1fr)",
          }
        }
        //   templateAreas={`"first first first first first"
        // "" "second "third" "" ""`}
        alignItems={"center"}
        maxW={"1200px"}
      >
        <GridItem colSpan={3}>
          <GalleryCard activity={data?.activities[0]} position={1} />
        </GridItem>
        <GridItem colSpan={3}>
          <Flex
            gap={{
              // 6
              lg: "6",
              md: "1",
            }}
            flexDirection={{
              base: "column",
              lg: "row",
            }}
            width={"100%"}
            // justifyItems={"center"}
            alignItems={"center"}
            justifyContent={"center"}
          >
            <GalleryCard
              activity={data?.activities[1]}
              position={2}
              superlatives={athSuperlatives(data?.activities[1])}
            />
            <GalleryCard
              activity={data?.activities[2]}
              position={3}
              superlatives={athSuperlatives(data?.activities[2])}
            />
          </Flex>
        </GridItem>
        {data?.activities.slice(3).map((activity, index) => {
          return (
            <GridItem
              key={`activity-${index}`}
              colSpan={{
                base: 3,
                lg: 1,
              }}
            >
              <GalleryCard
                activity={activity}
                position={index + 4}
                superlatives={athSuperlatives(activity)}
              />
            </GridItem>
          );
        })}
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
  superlatives?: Record<string, SuperlativeEntry<any>>;
}> = ({ activity, position, superlatives }) => {
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
      {/* Superlatives */}
      <Box position="relative" top="0px" left="0px">
        <Stack dir="column" position={"absolute"} top="70px" left="-20px">
          {
            superlatives &&
              Object.entries(superlatives).map(([key, value]) => {
                return <Superlative category={key} entry={value} />;
              })

            // Object.entries(superlatives).map((item) => (
            //   <Superlative category="earliest_start" />
            //   // <Avatar key={item} src={""} name={item} />
            // ))
          }
        </Stack>
      </Box>
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
              maxW="170px"
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
                ? `https://www.strava.com/activities/${activity.activity_id}`
                : `https://www.strava.com/athletes/${activity.athlete_id}`
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
    <ResponsiveCard m={"10px"} w="350px" maxW="350px" h="270px" {...props}>
      {children}
    </ResponsiveCard>
  );
};

// {Array.from({ length: 20 }).map(e =>
//   <GalleryCard />
// )}
