import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge } from "@chakra-ui/react";

export const AthleteAvatar: FC<{
  athlete: AthleteSummary
  size?: string
}> = ({ athlete, size = "md" }) => {
  const name = athlete.firstname + " " + athlete.lastname
  return <Avatar
    name={name}
    src={athlete.profile_pic_link}
    size={size}
  >
    {/* <AvatarBadge boxSize='1.25em' bg='green.500' /> */}
  </Avatar>
}