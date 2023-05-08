import { FC } from "react";
import { AthleteSummary } from "../../api/typesGenerated";
import { Avatar, AvatarBadge, AvatarProps, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";

export const Activity: FC<{

}> = ({  }) => {

  const {activity_id} = useParams()
  return <>
    Hi
  </>
}


{/* <iframe height='405' width='590' frameborder='0' allowtransparency='true' scrolling='no' src='https://www.strava.com/segments/7041089/embed'></iframe> */}