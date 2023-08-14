import { FC } from "react";
import { Flex, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getAthlete, getAthleteHugels, getRoute } from "../../api/rest";
import { useQuery } from "@tanstack/react-query";
import { NotFound } from "../404/404";
import { Loading } from "../../components/Loading/Loading";
import { AthletePageHeader } from "./AthleteHeader";

export const AthleteOtherPage: FC<{}> = ({}) => {
  const { athlete_id } = useParams();

  const queryKey = ["athlete-summary"];
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
    queryKey: queryHugelsKey,
    enabled: !!athlete_id,
    queryFn: () => getAthleteHugels(athlete_id || ""),
  });

  if (athleteLoading || athleteHugelsLoading) {
    return <Loading />;
  }

  if (!athlete_id || !athleteData || !athleteHugelsData) {
    return <NotFound />;
  }

  return (
    <>
      <AthletePageHeader
        athlete={athleteData}
        hugel_efforts={athleteHugelsData}
      />
    </>
  );
};
