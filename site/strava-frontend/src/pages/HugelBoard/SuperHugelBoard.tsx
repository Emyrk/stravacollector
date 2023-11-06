import { FC } from "react";
import {
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useTheme,
  Heading,
  Flex,
} from "@chakra-ui/react";
import { getSuperHugelLeaderBoard } from "../../api/rest";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar";
import * as TypesGen from "./../../api/typesGenerated";
import { HugelBoardGallery } from "./HugelBoardGallery";
import { HugelBoardTable } from "./HugelBoardTable";

export const SuperHugelBoard: FC = () => {
  const queryKey = ["hugel-leaderboard"];
  const {
    data: hugelLeaderboard,
    error: hugelLeaderboardError,
    isLoading: hugelLoading,
    isFetched: hugelFetched,
  } = useQuery({
    queryKey,
    queryFn: getSuperHugelLeaderBoard,
  })

  return <>
    <Flex w="100%" justifyContent={"center"}>
      <Heading>Super Scored Hugel</Heading>
    </Flex>
    <Tabs isFitted align="center" p='0 1rem'>
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
}

