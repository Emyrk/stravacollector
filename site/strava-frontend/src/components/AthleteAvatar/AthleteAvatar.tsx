import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps, Text } from "@chakra-ui/react";

export const AthleteAvatar: FC<{
  firstname: string;
  lastname: string;
  athlete_id: string;
  profile_pic_link: string;
  username: string;
  size?: string;
  styleProps?: AvatarProps;
  hugel_count?: number;
}> = ({
  firstname: firstName,
  lastname: lastName,
  athlete_id: athleteID,
  username,
  profile_pic_link: profilePicLink,
  size = "md",
  hugel_count: hugelCount,
  styleProps,
}) => {
  let name = firstName + " " + lastName;
  if (name === "") {
    name = username;
  }
  let boxSize = "1em";
  let fontSize = "sm";
  if (size == "xxl") {
    boxSize = "2.5em";
    fontSize = "md";
  }

  return (
    <Avatar name={name} src={profilePicLink} size={size} {...styleProps}>
      {/* Disable Hugel count for now */}
      {/* {hugelCount && hugelCount > 0 ? (
        <AvatarBadge
          boxSize={boxSize}
          bg="green.500"
          fontWeight={"bold"}
          borderColor="transparent"
        >
          <Text fontSize={fontSize}>{hugelCount}</Text>
        </AvatarBadge>
      ) : (
        ""
      )} */}
    </Avatar>
  );
};
