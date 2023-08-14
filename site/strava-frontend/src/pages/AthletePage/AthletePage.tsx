import { FC } from "react";
import { Flex, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getAthlete, getRoute } from "../../api/rest";
import { useQuery } from "@tanstack/react-query";
import { NotFound } from "../404/404";
import { Loading } from "../../components/Loading/Loading";

export const AthletePage: FC<{}> = ({}) => {
  const { athlete_id } = useParams();

  const queryKey = ["athlete", athlete_id];
  const {
    data: athleteData,
    error: athleteError,
    isLoading: athleteLoading,
    isFetched: athleteFetched,
  } = useQuery({
    queryKey,
    enabled: !!athlete_id,
    queryFn: () => getAthlete(athlete_id || ""),
  });

  const queryHugelsKey = ["hugels", athlete_id];
  const {
    data: athleteHugelsData,
    error: athleteHugelsError,
    isLoading: athleteHugelsLoading,
    isFetched: athleteHugelsFetched,
  } = useQuery({
    queryKey,
    enabled: !!athlete_id,
    queryFn: () => getAthlete(athlete_id || ""),
  });

  if (!athlete_id) {
    return <NotFound />;
  }

  if (athleteLoading || athleteHugelsLoading) {
    return <Loading />;
  }

  return (
    <>
      <Flex></Flex>
      <Text>{athlete_id} Page</Text>
    </>
  );
};
