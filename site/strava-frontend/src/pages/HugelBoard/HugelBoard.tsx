import { FC, useEffect } from "react";
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
import { useParams } from "react-router-dom";

export interface HugelBoardProps {
  disableSuperlatives?: boolean;
  data?: TypesGen.HugelLeaderBoard | TypesGen.SuperHugelLeaderBoard;
  error?: Error | unknown;
  isLoading: boolean;
  isFetched: boolean;
}

export const HugelBoard: FC = () => {
  const params = new URLSearchParams(window.location.search);

  const { year } = useParams();
  // Default to this year
  const yearNumber = parseInt(year || "2024");
  // TODO: Remove this disable when 2024 ride is complete.
  const disableSuperlatives = year !== "2024";
  const lite = params.get("lite") === "true";

  const queryKey = ["hugel-leaderboard", year, lite];
  const {
    data: hugelLeaderboard,
    error: hugelLeaderboardError,
    isLoading: hugelLoading,
    isFetched: hugelFetched,
  } = useQuery({
    queryKey,
    enabled: yearNumber != 2025,
    queryFn: () => {
      return getHugelLeaderBoard(yearNumber, lite);
    },
  });

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
      { yearNumber === 2025 ? (
        <Container maxW="7xl">
          <Alert status="info" borderRadius={"md"}>
            <AlertIcon />
            <Box>
              <AlertTitle>2025 Das Hügel is not live yet!</AlertTitle>
              <AlertDescription>
                The 2025 Das Hügel event will be live on November 8th, 2025.
                Ride will meet at 7:00 AM with a 7:15 AM rollout.
                <br />
                Please check back then to see the results and leaderboards.
              </AlertDescription>
            </Box>
          </Alert>
        </Container>
      ) : (
        <Tabs isFitted align="center" p="0 1rem">
          <TabList>
            <Tab>🖼️ Gallery</Tab>
            <Tab>📋 Table</Tab>
          </TabList>
          <TabPanels>
            <TabPanel key="gallery">
              <HugelBoardGallery
                disableSuperlatives={disableSuperlatives}
                data={hugelLeaderboard}
                error={hugelLeaderboardError}
                isLoading={hugelLoading}
                isFetched={hugelFetched}
              />
            </TabPanel>
            <TabPanel key="table">
              <HugelBoardTable
                data={hugelLeaderboard}
                error={hugelLeaderboardError}
                isLoading={hugelLoading}
                isFetched={hugelFetched}
              />
            </TabPanel>
          </TabPanels>
        </Tabs>
        )
      }
    </>
  );
};
