import { Box, Container, Flex, Text } from "@chakra-ui/react";
import { FC, useState } from "react";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { useAuthenticated } from "../../contexts/Authenticated";
import { useParams } from "react-router-dom";
import { getAthlete, getAthleteSyncSummary } from "../../api/rest";
import { AthleteSummary } from "../../api/typesGenerated";
import { useQuery } from "@tanstack/react-query";
import { Loading } from "../../components/Loading/Loading";
import { EddingtonChart } from "./EddingtonChart";




export const Eddington: FC<{}> = ({}) => {
  const { athlete_id } = useParams();
  const { authenticatedUser } = useAuthenticated();
  const [ athlete, setAthlete] = useState<AthleteSummary>();

  const queryKey = ["athlete", athlete_id];
    const {
      data: athleteData,
      error: athleteError,
      isLoading: athleteLoading,
      isFetched: athleteFetched,
      // refetch: athleteRefetch,
    } = useQuery({
      queryKey,
      enabled: !!athlete_id,
      queryFn: () =>
        getAthlete(athlete_id || "me"),
      onSuccess: (data) => {
        if(data) {
          setAthlete(data);
        }
      },
      onError: (error) => {
        console.error("Error fetching athlete data:", error);
      }
    });

  if (
    (!athlete || athleteLoading)
  ) {
    return <Loading />;
  }
    

  return (
    <>
     <Container maxW="3xl">
      <Flex flexDirection={"column"} gap="60px">
          <Box textAlign={"center"}>
            <AthleteAvatar
              styleProps={{ marginBottom: "20px" }}
              size="xxl"
              {...athlete}
            ></AthleteAvatar>
            <Flex
              flexDirection="row"
              gap="10px"
              width={"100%"}
              alignItems={"center"}
              justifyContent={"center"}
            >
              <StravaLink
                href={`https://strava.com/athletes/${athlete.athlete_id}`}
                target="_blank"
              />
              <Text fontSize="2xl" fontWeight="bold">
                {athlete.firstname} {athlete.lastname}
              </Text>
            </Flex>
          </Box>
        </Flex>
     </Container>
     <EddingtonChart />
    </>
  );
};

