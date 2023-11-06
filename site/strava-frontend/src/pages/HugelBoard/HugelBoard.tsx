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

export interface HugelBoardProps {
  data?: TypesGen.HugelLeaderBoard | TypesGen.SuperHugelLeaderBoard;
  error?: Error | unknown;
  isLoading: boolean;
  isFetched: boolean;
}

export const HugelBoard: FC = () => {
  let present = true;
  if (window.location.search.includes("all=true")) {
    present = false;
  }
  const queryKey = ["hugel-leaderboard"];
  const {
    data: hugelLeaderboard,
    error: hugelLeaderboardError,
    isLoading: hugelLoading,
    isFetched: hugelFetched,
  } = useQuery({
    queryKey,
    queryFn: () => getHugelLeaderBoard(present),
  });

  return (
    <>
      <Flex w="100%" justifyContent={"center"}>
        <Heading>Tour Das Hugel</Heading>
      </Flex>
      <Tabs isFitted align="center" p="0 1rem">
        <TabList>
          <Tab>🖼️ Gallery</Tab>
          <Tab>📋 Table</Tab>
        </TabList>
        <TabPanels>
          <TabPanel key="gallery">
            <HugelBoardGallery
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
    </>
  );
};
