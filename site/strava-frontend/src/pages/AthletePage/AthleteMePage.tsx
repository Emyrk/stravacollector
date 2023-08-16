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
  WithCSSVar,
  Th,
  IconButton,
  Select,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
} from "@chakra-ui/react";
import { Dict } from "@chakra-ui/utils";
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
import {
  getCoreRowModel,
  Column,
  ColumnDef,
  useReactTable,
  flexRender,
  getPaginationRowModel,
} from "@tanstack/react-table";
import {
  ArrowRightIcon,
  ArrowLeftIcon,
  ChevronRightIcon,
  ChevronLeftIcon,
} from "@chakra-ui/icons";

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
    queryFn: () => getAthleteSyncSummary(athlete_id || "me"),
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
  const theme = useTheme();
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
      <AthleteMeActivitiesTable
        data={summary.synced_activities}
        columns={activityColumns(theme)}
      />
    </Stack>
  );
};

const activityColumns = (theme: WithCSSVar<Dict>) => {
  return [
    {
      header: "Synced",
      accessor: "synced",
      cell: (row) => {
        const act = row.row.original;
        return act.synced ? (
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
        );
      },
    },
    {
      header: "Activity",
      accessor: "activity_summary.name",
      cell: (row) => {
        const act = row.row.original;
        return (
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
        );
      },
    },
    {
      header: "Date",
      // accessor: "activity_summary.start_date_local",
      cell: (row) => {
        const act = row.row.original;
        return FormatDate(act.activity_summary.start_date_local);
      },
    },
  ] as ColumnDef<SyncActivitySummary>[];
};

interface ReactTableProps<T extends object> {
  data: T[];
  columns: ColumnDef<T>[];
}

const AthleteMeActivitiesTable = <T extends object>({
  data,
  columns,
}: ReactTableProps<T>) => {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });
  const pageIndex = 0;

  return (
    <TableContainer
      sx={{
        td: { padding: "2px" },
      }}
    >
      <Table>
        <Thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <Tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <Th key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody>
          {table.getRowModel().rows.map((row) => (
            <Tr key={row.id}>
              {row.getVisibleCells().map((cell, index) => (
                <Td
                  key={cell.id}
                  width={index === 0 ? "20px" : ""}
                  textAlign={index === 0 ? "center" : "inherit"}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </Td>
              ))}
            </Tr>
          ))}
        </Tbody>
      </Table>

      <Flex justifyContent="space-between" m={4} alignItems="center">
        <Flex>
          <Tooltip label="First Page">
            <IconButton
              // onClick={() => gotoPage(0)}
              // isDisabled={!canPreviousPage}
              icon={<ArrowLeftIcon h={3} w={3} />}
              mr={4}
              aria-label={""}
            />
          </Tooltip>
          <Tooltip label="Previous Page">
            <IconButton
              // onClick={previousPage}
              // isDisabled={!canPreviousPage}
              icon={<ChevronLeftIcon h={6} w={6} />}
              aria-label={""}
            />
          </Tooltip>
        </Flex>

        <Flex alignItems="center">
          <Text flexShrink="0" mr={8}>
            Page{" "}
            <Text fontWeight="bold" as="span">
              {pageIndex + 1}
            </Text>{" "}
            of{" "}
            <Text fontWeight="bold" as="span">
              {/* {pageOptions.length} */}
              100
            </Text>
          </Text>
          <Text flexShrink="0">Go to page:</Text>{" "}
          <NumberInput
            ml={2}
            mr={8}
            w={28}
            min={1}
            // max={pageOptions.length}
            max={100}
            // onChange={(value) => {
            //   const page = value ? value - 1 : 0;
            //   gotoPage(page);
            // }}
            defaultValue={pageIndex + 1}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>
          <Select
            w={32}
            // value={pageSize}
            // onChange={(e) => {
            //   setPageSize(Number(e.target.value));
            // }}
          >
            {[10, 20, 30, 40, 50].map((pageSize) => (
              <option key={pageSize} value={pageSize}>
                Show {pageSize}
              </option>
            ))}
          </Select>
        </Flex>

        <Flex>
          <Tooltip label="Next Page">
            <IconButton
              // onClick={nextPage}
              // isDisabled={!canNextPage}
              icon={<ChevronRightIcon h={6} w={6} />}
              aria-label={""}
            />
          </Tooltip>
          <Tooltip label="Last Page">
            <IconButton
              // onClick={() => gotoPage(pageCount - 1)}
              // isDisabled={!canNextPage}
              icon={<ArrowRightIcon h={3} w={3} />}
              ml={4}
              aria-label={""}
            />
          </Tooltip>
        </Flex>
      </Flex>
    </TableContainer>
  );
};
