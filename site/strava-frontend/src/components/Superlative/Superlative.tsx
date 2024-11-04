import { Avatar, AvatarProps, Box, Stack, Text } from "@chakra-ui/react";
import { FC, ReactElement } from "react";
import { SuperlativeEntry } from "../../api/typesGenerated";
import { Tooltip, TooltipProps } from "@chakra-ui/react";
import { ResponsiveCard } from "../ResponsiveCard/ResponsiveCard";
import {
  ElapsedDurationText,
  FormatDate,
  FormatDateTime,
} from "../../pages/HugelBoard/CalcActivity";
import { DistanceToMiles } from "../../lib/Distance/Distance";

export type SuperlativeProps = AvatarProps & {
  category: string;
  entry: SuperlativeEntry<any>;
};

export const Superlative: FC<SuperlativeProps> = ({
  category,
  entry,
  ...props
}) => {
  const [src, title, value] = SuperlativeLookup(category, entry);

  return (
    <Tooltip
      placement="right-start"
      background={"none"}
      p="0px"
      m="0px"
      label={<SuperlativeCard title={title} value={value} />}
    >
      <Avatar key={category} src={`/img/superlatives/${src}`} name={category} />
    </Tooltip>
  );
};

export const SuperlativeCard: FC<{ title: string; value: any }> = ({
  title,
  value,
}) => {
  return (
    <ResponsiveCard
      width={"270px"}
      // height={"100px"}
      opacity={"93%"}
      color={"white"}
      p={"10px"}
      // boxShadow={"#fc4c02 0px 3px 6px"}
    >
      <Stack>
        <Text fontSize={"1.2em"} fontWeight={800}>
          {title}
        </Text>
        {value}
      </Stack>
    </ResponsiveCard>
  );
};

const SuperlativeLookup = (
  category: string,
  entry: SuperlativeEntry<any>
): [string, string, ReactElement] => {
  switch (category) {
    case "early_bird":
    case "earliest_start":
      const d = new Date(entry.value);
      return [
        "EarlyBird.png",
        "Early Bird",
        <Text>
          Gets the worm with their {FormatDateTime(entry.value)} start time.
        </Text>,
      ];
    case "night_owl":
    case "latest_end":
      return [
        "NightOwl.png",
        "Night Owl",
        <Text>
          Aren’t you glad you didn’t wait up with their{" "}
          {FormatDateTime(entry.value)} end time?
        </Text>,
      ];
    case "most_stoppage":
      return [
        "CoffeeBreak.png",
        "Coffee Break",
        <Text>
          Stopped and smelled the roses with {Math.floor(entry.value / 3600)}{" "}
          hrs and {((entry.value / 60) % 60).toFixed(0)} minutes of stoppage.
        </Text>,
      ];
    case "least_stoppage":
      return [
        "Dory.png",
        "Dory",
        <Text>
          Just keep swimming. Only {(entry.value / 60).toFixed(0)} minutes of
          stoppage.
        </Text>,
      ];
    case "most_avg_watts":
    case "most_watts":
      return [
        "TheEdison.png",
        "The Edison",
        <Text>Powering Austin with {entry.value} average watts.</Text>,
      ];
    case "most_avg_cadence":
      return [
        "Roadrunner.png",
        "Roadrunner",
        <Text>Legs a'blur with average cadence of {entry.value} rpm.</Text>,
      ];
    case "least_avg_cadence":
      return [
        "Mortar&Pestle.png",
        "Mortar & Pestle",
        <Text>
          Grinding so hard with average cadence of {entry.value.toFixed(0)} rpm.
        </Text>,
      ];
    case "most_suffer":
      return [
        "Masochist.png",
        "Masochist",
        <Text>
          Definitely type 2 fun with this {entry.value} suffer score.
        </Text>,
      ];
    case "most_achievements":
      return [
        "Overachiever.png",
        "Overachiever",
        <Text>Thinking they're so cool with {entry.value} achievements.</Text>,
      ];
    case "longest_ride":
      return [
        "Wanderer.png",
        "Wanderer",
        <Text>
          Must've gotten lost taking {DistanceToMiles(entry.value).toFixed(1)}{" "}
          miles to finish.
        </Text>,
      ];
    case "shortest_ride":
      return [
        "MVP.png",
        "MVP",
        <Text>
          Most Vigilant Path-Follower took no detours with only{" "}
          {DistanceToMiles(entry.value).toFixed(1)} miles to finish.
        </Text>,
      ];
    case "least_avg_hr":
      return [
        "Yawner.png",
        "Yawner",
        <Text>
          Not working hard with a {entry.value.toFixed(0)} bpm average heart
          rate.
        </Text>,
      ];
    case "most_avg_hr":
      // TODO
      return [
        "CardiacArrest.png",
        "Cardiac Arrest",
        <Text>
          Anyone have a defib? Gonna need it with this {entry.value} bpm average
          heart rate.
        </Text>,
      ];
    case "least_avg_speed":
      return [
        "Turtle.png",
        "Turtle",
        <Text>
          Taking their sweet time with an average of{" "}
          {DistanceToMiles(entry.value * 3600).toFixed(2)} mph.
        </Text>,
      ];
    case "most_avg_speed":
      return [
        "Hare.png",
        "Hare",
        <Text>
          Must’ve wanted the ride to be over with an average of{" "}
          {DistanceToMiles(entry.value * 3600).toFixed(2)} mph.
        </Text>,
      ];
  }

  return ["", category, <></>];
};
