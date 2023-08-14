import { FC } from "react";
import {
  AthleteHugelActivities,
  AthleteSummary,
} from "../../api/typesGenerated";
import {
  AvatarBadge,
  Box,
  Container,
  Flex,
  Show,
  Table,
  TableContainer,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import { StravaLink } from "../../components/StravaLink/StravaLink";
import { ElapsedDurationText, FormatDate } from "../HugelBoard/CalcActivity";

export const AthletePageHeader: FC<{
  athlete: AthleteSummary;
  hugel_efforts: AthleteHugelActivities;
}> = ({ athlete, hugel_efforts }) => {
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
          <TableContainer>
            <Table variant="striped">
              <Thead>
                <Tr>
                  <Th>Hugels Completed</Th>
                  <Show above="md">
                    <Th>Segment Time</Th>
                  </Show>
                  <Th>Date</Th>
                </Tr>
              </Thead>
              <Tbody>
                {hugel_efforts.activities.length == 0 && (
                  <Tr>
                    <Td textAlign="center" colSpan={3}>
                      No hugels completed
                    </Td>
                  </Tr>
                )}
                {hugel_efforts.activities.map((effort) => {
                  return (
                    <Tr key={effort.summary.activity_id}>
                      <Td>
                        <Flex
                          flexDirection={"row"}
                          alignItems={"center"}
                          gap="10px"
                        >
                          <StravaLink
                            href={`https://strava.com/activities/${effort.summary.activity_id}`}
                            target="_blank"
                          />
                          {effort.summary.name}
                        </Flex>
                      </Td>
                      <Show above="md">
                        <Td>
                          {ElapsedDurationText(effort.total_time_seconds)}
                        </Td>
                      </Show>
                      <Td>{FormatDate(effort.summary.start_date_local)}</Td>
                    </Tr>
                  );
                })}
              </Tbody>
            </Table>
          </TableContainer>
        </Flex>
      </Container>
    </>
  );
};
