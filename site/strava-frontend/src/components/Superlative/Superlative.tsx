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
  return (
    <Tooltip
      placement="right-start"
      background={"none"}
      p="0px"
      m="0px"
      label={<SuperlativeCard category={category} entry={entry} />}
    >
      <Avatar key={category} src={""} name={category} />
    </Tooltip>
  );
};

export const SuperlativeCard: FC<SuperlativeProps> = ({ category, entry }) => {
  const [title, value] = mutate(category, entry);

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

const mutate = (
  category: string,
  entry: SuperlativeEntry<any>
): [string, string] => {
  switch (category) {
    case "early_bird":
    case "earliest_start":
      const d = new Date(entry.value);
      return ["Early Bird", `Started at ${FormatDateTime(entry.value)}!`];
    case "night_owl":
    case "latest_end":
      return ["Night Owl", `Ended at ${FormatDateTime(entry.value)}!`];
    case "most_stoppage":
      return [
        "Most Relaxed",
        `Total break of ${ElapsedDurationText(entry.value, false)}`,
      ];
    case "least_stoppage":
      // TODO: Rename
      return [
        "Extreme",
        `Total break of ${ElapsedDurationText(entry.value, false)}`,
      ];
    case "most_watts":
      return ["Watt Machine", entry.value as string];
    case "most_cadence":
      return ["Spin to Win", entry.value as string];
    case "least_cadence":
      return ["Grinder", entry.value as string];
    case "most_suffer":
      return ["Most Pain", entry.value as string];
    case "most_achievements":
      return ["Most Decorated", entry.value as string];
    case "longest_ride":
      return ["Has no car", entry.value as string];
    case "shortest_ride":
      return ["Most Efficient", entry.value as string];
  }

  return [category, entry.value as string];
};
