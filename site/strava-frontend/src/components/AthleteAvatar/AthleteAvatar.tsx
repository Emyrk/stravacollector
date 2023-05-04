import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps, Text } from "@chakra-ui/react";

export const AthleteAvatar: FC<{
  firstName: string
  lastName: string
  athleteID: string
  profilePicLink: string
  username: string
  size?: string
  styleProps?: AvatarProps
  hugelCount?: number
}> = ({ firstName, lastName, athleteID, username, profilePicLink, size = "md", hugelCount, styleProps }) => {
  let name = firstName + " " + lastName
  if (name === "") {
    name = username
  }
  return <Avatar
    name={name}
    src={profilePicLink}
    size={size}
    {...styleProps}
  >
    {hugelCount && hugelCount > 0 &&
      <AvatarBadge
        boxSize='1em'
        bg='green.500'
        fontWeight={"bold"}
        borderColor="transparent"
      >
        <Text fontSize={"sm"} >
          {hugelCount}
        </Text>
      </AvatarBadge>
    }
  </Avatar>
}