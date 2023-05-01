import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge } from "@chakra-ui/react";

export const AthleteAvatar: FC<{
  firstName: string
  lastName: string
  athleteID: number
  profilePicLink: string
  username: string
  size?: string
}> = ({ firstName, lastName, athleteID, username, profilePicLink, size = "md" }) => {
  let name = firstName + " " + lastName
  if (name === "") {
    name = username
  }
  return <Avatar
    name={name}
    src={profilePicLink}
    size={size}
  >
    {/* <AvatarBadge boxSize='1.25em' bg='green.500' /> */}
  </Avatar>
}