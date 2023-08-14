import { FC } from "react";
import {
  AbsoluteCenter,
  AlertIcon,
  Box,
  Checkbox,
  Container,
  Divider,
  Flex,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Table,
  TableContainer,
  Tabs,
  Tbody,
  Td,
  Text,
  Thead,
  Stack,
  Alert,
  Tr,
  AlertDescription,
  AlertTitle,
  CircularProgress,
  CircularProgressLabel,
  useTheme,
  Tooltip,
} from "@chakra-ui/react";
import { Link, useParams } from "react-router-dom";
import {
  getAthlete,
  getAthleteHugels,
  getAthleteSyncSummary,
  getRoute,
} from "../../api/rest";
import { useQuery } from "@tanstack/react-query";
import { NotFound } from "../404/404";
import { Loading } from "../../components/Loading/Loading";
import { useAuthenticated } from "../../contexts/Authenticated";
import { AthletePageHeader } from "./AthleteHeader";
import {
  AthleteSyncSummary,
  SyncActivitySummary,
} from "../../api/typesGenerated";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCircleXmark,
  faCircleCheck,
} from "@fortawesome/free-solid-svg-icons";
import { FormatDate } from "../HugelBoard/CalcActivity";
import { StravaLink } from "../../components/StravaLink/StravaLink";

export const AthleteMePage: FC<{}> = ({}) => {
  const { athlete_id } = useParams();
  const { authenticatedUser } = useAuthenticated();

  const queryKey = ["sync-summary"];
  const {
    data: athleteData,
    error: athleteError,
    isLoading: athleteLoading,
    isFetched: athleteFetched,
  } = useQuery({
    queryKey,
    enabled: !!athlete_id,
    queryFn: () => getAthleteSyncSummary(),
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
        athlete={athleteData?.athlete_summary}
        hugel_efforts={athleteHugelsData}
      />
      <Box position="relative" padding="10">
        <Divider />
        <AbsoluteCenter bg="chakra-body-bg" px="4">
          <Text fontSize={"1.3em"}>Strava Sync Details</Text>
        </AbsoluteCenter>
      </Box>

      <Container maxW="3xl">
        <AthleteMeTotals summary={athleteData} />
      </Container>
    </>
  );
};

const AthleteMeTotals: FC<{ summary: AthleteSyncSummary }> = ({ summary }) => {
  const load = summary.athlete_load;
  return (
    <Stack spacing={0.5}>
      <Alert status={load.earliest_activity_done ? "success" : "warning"}>
        <AlertIcon />
        {load.earliest_activity_done
          ? "All historical activities have been loaded!"
          : "Not all historical activities have been loaded. Loading starts with the oldest activity, and walks forward in time."}
        <br />- Earliest activity synced:
        {FormatDate(load.earliest_activity)}
        <br />- Latest activity loaded:
        {FormatDate(load.last_backload_activity_start)}
      </Alert>
      <Alert
        status={
          summary.total_summary === summary.total_detail ? "success" : "warning"
        }
      >
        <Flex flexDirection={"row"} alignItems={"center"} gap="15px">
          <Box>
            <CircularProgress
              value={
                Math.ceil(summary.total_detail / summary.total_summary) * 100
              }
              color="green.400"
            >
              <CircularProgressLabel>
                {Math.ceil(
                  (summary.total_detail / summary.total_summary) * 100
                )}
                %
              </CircularProgressLabel>
            </CircularProgress>
          </Box>
          <Text>
            Loaded activities still need to be synced one by one to find all
            segment details. {summary.total_detail} of {summary.total_summary}{" "}
            activities are complete.
          </Text>
        </Flex>
      </Alert>

      {load.last_load_error && (
        <Alert flexDirection={"column"} status="error">
          <Flex>
            <AlertIcon />
            <AlertTitle>
              Error message on the last sync attempt. Please report this error.
            </AlertTitle>
          </Flex>

          <Box>
            <AlertDescription maxWidth="sm">
              <pre>{load.last_load_error}</pre>
            </AlertDescription>
          </Box>
        </Alert>
      )}
      {/* <Alert status="success">
        <AlertIcon />
        Data uploaded to the server. Fire on!
      </Alert>

      <Alert status="warning">
        <AlertIcon />
        Seems your account is about expire, upgrade now
      </Alert>

      <Alert status="info">
        <AlertIcon />
        Chakra is going live on August 30th. Get ready!
      </Alert> */}
      <AthleteMeHugelTable summary={summary} />
    </Stack>
  );
};

const AthleteMeHugelTable: FC<{ summary: AthleteSyncSummary }> = ({
  summary,
}) => {
  const theme = useTheme();
  return (
    <TableContainer
      sx={{
        td: { padding: "2px" },
      }}
    >
      <Table>
        <Thead>
          <Tr>
            <Td>Synced</Td>
            <Td>Activity</Td>
            <Td>Date</Td>
          </Tr>
        </Thead>
        <Tbody>
          {summary.synced_activities.map((act) => {
            return (
              <Tr key={act.activity_summary.activity_id} padding={"0px"}>
                <Td width="70px" textAlign={"center"}>
                  {act.synced ? (
                    <Tooltip label={`Synced on ${FormatDate(act.synced_at)}`}>
                      <FontAwesomeIcon
                        cursor={"pointer"}
                        color={theme.colors.green["500"]}
                        icon={faCircleCheck}
                      />
                    </Tooltip>
                  ) : (
                    <FontAwesomeIcon
                      color={theme.colors.red["500"]}
                      icon={faCircleXmark}
                    />
                  )}
                </Td>
                <Td>
                  <Flex flexDirection={"row"} gap="5px" alignItems={"center"}>
                    <StravaLink
                      href={`https://www.strava.com/activities/${act.activity_summary.activity_id}`}
                      target="_blank"
                      height={"24px"}
                      width={"24px"}
                    />
                    <Text
                      textAlign={"center"}
                      fontWeight={"bold"}
                      maxW="300px"
                      isTruncated
                    >
                      {act.activity_summary.name}
                    </Text>
                  </Flex>
                </Td>
                <Td>{FormatDate(act.activity_summary.start_date_local)}</Td>
              </Tr>
            );
          })}
        </Tbody>
      </Table>
    </TableContainer>
  );
};
