import { FC, useEffect, useState } from "react";
import {
  Flex,
  Grid,
  Table,
  Thead,
  Tbody,
  Tfoot,
  Tr,
  Th,
  Td,
  TableCaption,
  TableContainer,
  Box,
  Spinner,
  Text,
  Link,
  LinkBox,
  LinkOverlay,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  useToast,
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useTheme,
  Heading,
  Container,
  ButtonGroup,
  Button,
} from "@chakra-ui/react";
import {
  getErrorDetail,
  getErrorMessage,
  getHugelLeaderBoard,
} from "../../api/rest";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import * as TypesGen from "./../../api/typesGenerated";
import { HugelBoardGallery } from "./HugelBoardGallery";
import { HugelBoardTable } from "./HugelBoardTable";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { SexFilterButtons } from "../../components/SexFilter/SexFilter";

export interface HugelBoardProps {
  disableSuperlatives?: boolean;
  data?: TypesGen.HugelLeaderBoard | TypesGen.SuperHugelLeaderBoard;
  error?: Error | unknown;
  isLoading: boolean;
  isFetched: boolean;
}

export type SexFilter = "all" | "M" | "F";

export const HugelBoard: FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const theme = useTheme();

  
  // const [sexFilter, setSexFilter] = useState<SexFilter>("all");

  const { year } = useParams();
  // Default to this year
  const yearNumber = parseInt(year || "2024");
  // TODO: Remove this disable when 2024 ride is complete.
  const disableSuperlatives = year !== "2024";
  const lite = searchParams.get("lite") === "true";
  const sexFilter =  searchParams.get("sex") as SexFilter || "all";
  const allowSexFilter = searchParams.get("sexboards") === "true" || yearNumber > 2024;


  const queryKey = ["hugel-leaderboard", year, lite];
  const {
    data: hugelLeaderboard,
    error: hugelLeaderboardError,
    isLoading: hugelLoading,
    isFetched: hugelFetched,
  } = useQuery({
    queryKey,
    queryFn: () => {
      return getHugelLeaderBoard(yearNumber, lite);
    },
  });

  // Filter activities by gender
const filteredData = hugelLeaderboard
  ? {
      ...hugelLeaderboard,
      activities: hugelLeaderboard.activities?.filter((activity) => {
        if (sexFilter === "all") return true;
        return activity.athlete.sex === sexFilter;
      }).map((activity, index) => ({
        ...activity,
        rank: index + 1,
      })),
    }
  : undefined;

  const handleSexFilterChange = (newSex: SexFilter) => {
    const params = new URLSearchParams(searchParams);
    if (newSex === "all") {
      params.delete("sex");
    } else {
      params.set("sex", newSex);
    }
    navigate({ search: params.toString() }, { replace: true });
  };

  return (
    <>
      <Flex
        w="100%"
        textAlign={"center"}
        justifyContent={"center"}
        direction={"column"}
      >
        <Heading>
          {year} Das Hügel {lite && "Lite"} Results
        </Heading>
        <Text color="gray.400" pt="5px">
          If your ride is not showing after 24hrs, please email{" "}
          <a href="mailto: help@dashugel.bike">help@dashugel.bike</a> with the
          link to your Hügel activity.
        </Text>
      </Flex>
      {(allowSexFilter) && (
        <SexFilterButtons 
            value={sexFilter} 
            onChange={handleSexFilterChange} 
          />
      )} 
        <Tabs isFitted align="center" p="0 1rem">
          <TabList>
            <Tab>🖼️ Gallery</Tab>
            <Tab>📋 Table</Tab>
          </TabList>
          <TabPanels>
            <TabPanel key="gallery">
              <HugelBoardGallery
                disableSuperlatives={disableSuperlatives}
                data={filteredData}
                error={hugelLeaderboardError}
                isLoading={hugelLoading}
                isFetched={hugelFetched}
              />
            </TabPanel>
            <TabPanel key="table">
              <HugelBoardTable
                data={filteredData}
                error={hugelLeaderboardError}
                isLoading={hugelLoading}
                isFetched={hugelFetched}
              />
            </TabPanel>
          </TabPanels>
        </Tabs>

    </>
  );
};
