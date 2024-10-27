import { Avatar, AvatarProps, Box, Stack, Text } from "@chakra-ui/react";
import { FC } from "react";
import { SuperlativeEntry } from "../../api/typesGenerated";
import { Tooltip, TooltipProps } from "@chakra-ui/react";

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
      label={<SuperlativeCard category={category} entry={entry} />}
    >
      <Avatar key={category} src={""} name={category} />
    </Tooltip>
  );
};

export const SuperlativeCard: FC<SuperlativeProps> = ({ category, entry }) => {
  const [title, value] = mutate(category, entry);

  return (
    <Box width={"200px"} height={"100px"}>
      <Stack>
        <Text>{title}</Text>
        <Text>{value}</Text>
      </Stack>
    </Box>
  );
};

const mutate = (
  category: string,
  entry: SuperlativeEntry<any>
): [string, string] => {
  switch (category) {
    case "earliest_start":
      return ["Early Bird", entry.value as string];
    case "latest_end":
      return ["Night Owl", entry.value as string];
    case "most_stoppage":
      return ["Most Relaxed", entry.value as string];
    case "least_stoppage":
      // TODO: Rename
      return ["Extreme", entry.value as string];
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
