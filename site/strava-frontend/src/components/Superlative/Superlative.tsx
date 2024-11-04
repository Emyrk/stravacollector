import { Avatar, AvatarProps, Box, Stack, Text } from "@chakra-ui/react";
import { FC } from "react";
import { SuperlativeEntry } from "../../api/typesGenerated";
import { Tooltip, TooltipProps } from "@chakra-ui/react";
import { ResponsiveCard } from "../ResponsiveCard/ResponsiveCard";
import {
  ElapsedDurationText,
  FormatDate,
  FormatDateTime,
} from "../../pages/HugelBoard/CalcActivity";

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
      width={"200px"}
      height={"100px"}
      border={"white"}
      borderStyle={"solid"}
      opacity={"99%"}
      color={"white"}
      p={"10px"}
      // boxShadow={"#fc4c02 0px 3px 6px"}
    >
      <Stack>
        <Text fontSize={"2em"}>{title}</Text>
        <Text>{value}</Text>
      </Stack>
    </ResponsiveCard>
  );
};

const SuperlativeLookup = (
  category: string,
  entry: SuperlativeEntry<any>
): [string, string, string] => {
  switch (category) {
    case "early_bird":
    case "earliest_start":
      const d = new Date(entry.value);
      return [
        "EarlyBird.png",
        "Early Bird",
        `Gets the worm with their ${FormatDateTime(entry.value)} start time.`,
      ];
    case "night_owl":
    case "latest_end":
      return [
        "NightOwn.png",
        "Night Owl",
        `Aren’t you glad you didn’t wait up with their ${FormatDateTime(
          entry.value
        )} end time?`,
      ];
    case "most_stoppage":
      return [
        "CoffeeBreak.png",
        "Coffee Break",
        `Stopped and smelled the roses with ${(entry.value / 3600).toFixed(
          0
        )} minutes of stoppage.`,
      ];
    case "least_stoppage":
      // TODO: Rename
      return [
        "Dory.png",
        "Dory",
        `Just keep swimming. Only ${(entry.value / 3600).toFixed(
          0
        )} minutes of stoppage.`,
      ];
    case "most_watts":
      return [
        "TheEdison.png",
        "The Edison",
        `Powering Austin with ${entry.value} average watts.`,
      ];
    case "most_cadence":
      return [
        "Roadrunner.png",
        "Roadrunner",
        `Legs a’blur with average cadence of ${entry.value} rpm.`,
      ];
    case "least_cadence":
      return [
        "Mortar&Pestle.png",
        "Mortar & Pestle",
        `Grinding so hard with average cadence of ${entry.value} rpm.`,
      ];
    case "most_suffer":
      return [
        "Masochist.png",
        "Masochist",
        `Definitely type 2 fun with this ${entry.value} suffer score.`,
      ];
    case "most_achievements":
      return [
        "Overachiever.png",
        "Overachiever",
        `Thinking they’re so cool with ${entry.value} achievements.`,
      ];
    case "longest_ride":
      return [
        "Wanderer.png",
        "Wanderer",
        `Must’ve gotten lost taking ${entry.value} miles to finish.`,
      ];
    case "shortest_ride":
      return [
        "MVP.png",
        "MVP",
        `Most Vigilant Path-Follower took no detours with only ${entry.value} miles to finish.`,
      ];
  }

  return ["", category, entry.value as string];
};
