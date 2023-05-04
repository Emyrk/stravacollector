import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps } from "@chakra-ui/react";

export const AthleteAvatar: FC<{
  firstName: string
  lastName: string
  athleteID: string
  profilePicLink: string
  username: string
  size?: string
  styleProps?: AvatarProps
}> = ({ firstName, lastName, athleteID, username, profilePicLink, size = "md", styleProps }) => {
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
    {/* <AvatarBadge boxSize='1.25em' bg='green.500' /> */}
  </Avatar>
}