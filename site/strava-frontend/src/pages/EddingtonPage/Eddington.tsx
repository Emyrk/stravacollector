import { Box, Container, Flex, Heading, Text, Link } from "@chakra-ui/react";
import { FC, useState } from "react";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { useAuthenticated } from "../../contexts/Authenticated";
import { useParams } from "react-router-dom";
import { getAthlete, getAthleteSyncSummary } from "../../api/rest";
import { AthleteSummary } from "../../api/typesGenerated";
import { useQuery } from "@tanstack/react-query";
import { Loading } from "../../components/Loading/Loading";
import { EddingtonChart } from "../../components/Eddington/EddingtonChart";
import { ErrorBox } from "../../components/ErrorBox/ErrorBox";
import { EddingtonAllChart } from "../../components/Eddington/EddingtonAllChart";




export const Eddington: FC<{}> = ({ }) => {
  const { athlete_id } = useParams();
  const { authenticatedUser } = useAuthenticated();
  const [athlete, setAthlete] = useState<AthleteSummary>();

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
      if (data) {
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

  if (athleteError) {
    return <ErrorBox error="Error fetching athlete data." detail={athleteError} />;
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

      <Box padding="10" textAlign="center"> </Box>

      <Flex
        w="100%"
        justifyContent={"center"}
        alignItems={"center"}
        textAlign="center"
      >
        <Flex flexDirection={"column"} pb="0.5em">
          <Text maxWidth={"1050px"} pt="1em">
              The chart below visualizes your Eddington Number — a metric representing the largest number <strong>n</strong> such that you’ve completed <strong>n</strong> rides of at least <strong>n</strong> miles. Each bar shows how many rides you’ve completed at a given mileage, while the diagonal red line indicates the <em>y = x</em> threshold. The point where your ride count drops below the line defines your current Eddington Number.
          </Text>
          <Link color={"#36c"}  href="https://en.wikipedia.org/wiki/Arthur_Eddington#Eddington_number_for_cycling"> Wikipedia</Link>
        </Flex>
      </Flex>

      <Container maxW="7xl">
        <EddingtonChart />

        { authenticatedUser?.athlete_id?.toString() === "2661162" ? <EddingtonAllChart /> : <></>}

      </Container>
    </>
  );
};

