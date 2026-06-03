import React, { useState } from 'react';
import {
  Box,
  Container,
  Heading,
  Text,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  VStack,
  HStack,
  Badge,
  useColorModeValue,
  Icon,
  Flex,
  SimpleGrid,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Avatar,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Card,
  CardHeader,
  CardBody,
  Divider,
} from '@chakra-ui/react';
import { FaTrophy, FaCrown, FaMedal, FaFire, FaMountain } from 'react-icons/fa';
import { mockKingLeaderboard, mockQueenLeaderboard, mockLeaderboardStats } from './mockData';

export const KingQueenAustin: React.FC = () => {
  const [selectedTab, setSelectedTab] = useState(0);
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const accentColor = useColorModeValue('orange.500', 'orange.300');

  return (
    <Container maxW="container.xl" py={8}>
      {/* Header */}
      <VStack spacing={6} mb={8}>
        <HStack spacing={4}>
          <Icon as={FaTrophy} boxSize={12} color={accentColor} />
          <Heading size="2xl">King & Queen of Austin</Heading>
        </HStack>
        <Text fontSize="lg" color="gray.500" textAlign="center">
          Top 50 Cycling Segments • Points-Based Rankings • Separate Male & Female Leaderboards
        </Text>

        {/* Stats Cards */}
        <SimpleGrid columns={{ base: 1, md: 4 }} spacing={6} w="full">
          <Stat
            bg={bgColor}
            p={4}
            borderRadius="lg"
            borderWidth="1px"
            borderColor={borderColor}
          >
            <StatLabel>Featured Segments</StatLabel>
            <StatNumber>{mockLeaderboardStats.totalSegments}</StatNumber>
            <StatHelpText>Across Austin, TX</StatHelpText>
          </Stat>

          <Stat
            bg={bgColor}
            p={4}
            borderRadius="lg"
            borderWidth="1px"
            borderColor={borderColor}
          >
            <StatLabel>Scoring Method</StatLabel>
            <StatNumber>10-1 pts</StatNumber>
            <StatHelpText>Top 10 finishes</StatHelpText>
          </Stat>

          <Stat
            bg={bgColor}
            p={4}
            borderRadius="lg"
            borderWidth="1px"
            borderColor={borderColor}
          >
            <StatLabel>King Leader</StatLabel>
            <StatNumber>{mockKingLeaderboard[0].totalPoints}</StatNumber>
            <StatHelpText>{mockKingLeaderboard[0].firstName} {mockKingLeaderboard[0].lastName}</StatHelpText>
          </Stat>

          <Stat
            bg={bgColor}
            p={4}
            borderRadius="lg"
            borderWidth="1px"
            borderColor={borderColor}
          >
            <StatLabel>Queen Leader</StatLabel>
            <StatNumber>{mockQueenLeaderboard[0].totalPoints}</StatNumber>
            <StatHelpText>{mockQueenLeaderboard[0].firstName} {mockQueenLeaderboard[0].lastName}</StatHelpText>
          </Stat>
        </SimpleGrid>
      </VStack>

      {/* Leaderboard Tabs */}
      <Tabs
        index={selectedTab}
        onChange={setSelectedTab}
        variant="soft-rounded"
        colorScheme="orange"
      >
        <TabList mb={4} justifyContent="center">
          <Tab>
            <Icon as={FaCrown} mr={2} />
            King Leaderboard
          </Tab>
          <Tab>
            <Icon as={FaCrown} mr={2} />
            Queen Leaderboard
          </Tab>
        </TabList>

        <TabPanels>
          {/* King Leaderboard */}
          <TabPanel>
            <LeaderboardTable data={mockKingLeaderboard} type="King" />
          </TabPanel>

          {/* Queen Leaderboard */}
          <TabPanel>
            <LeaderboardTable data={mockQueenLeaderboard} type="Queen" />
          </TabPanel>
        </TabPanels>
      </Tabs>

      {/* How It Works */}
      <Card mt={8}>
        <CardHeader>
          <Heading size="md">How It Works</Heading>
        </CardHeader>
        <CardBody>
          <VStack align="stretch" spacing={4}>
            <Text>
              <strong>Scoring:</strong> Athletes earn points based on their ranking on each of the 50 featured segments.
              1st place = 10 points, 2nd = 9 points, down to 10th = 1 point.
            </Text>
            <Text>
              <strong>Rankings:</strong> King (male) and Queen (female) compete separately for the top position.
            </Text>
            <Text>
              <strong>Segments:</strong> 50 carefully selected segments across Austin including iconic climbs like Mount Bonnell,
              Pennybacker Bridge Vista, and popular sprints on 360, Lamar Blvd, and South Congress.
            </Text>
            <Text>
              <strong>Updates:</strong> Leaderboard is recalculated daily based on best efforts.
            </Text>
          </VStack>
        </CardBody>
      </Card>
    </Container>
  );
};

interface LeaderboardTableProps {
  data: any[];
  type: 'King' | 'Queen';
}

const LeaderboardTable: React.FC<LeaderboardTableProps> = ({ data, type }) => {
  const bgColor = useColorModeValue('white', 'gray.800');
  const hoverBg = useColorModeValue('gray.50', 'gray.700');

  const getRankBadge = (rank: number) => {
    if (rank === 1) return <Icon as={FaTrophy} color="gold" boxSize={6} />;
    if (rank === 2) return <Icon as={FaMedal} color="silver" boxSize={6} />;
    if (rank === 3) return <Icon as={FaMedal} color="bronze" boxSize={6} />;
    return <Badge colorScheme="gray">{rank}</Badge>;
  };

  return (
    <Box overflowX="auto" bg={bgColor} borderRadius="lg" borderWidth="1px">
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Rank</Th>
            <Th>Athlete</Th>
            <Th isNumeric>Total Points</Th>
            <Th isNumeric>{type === 'King' ? 'KOMs' : 'QOMs'}</Th>
            <Th isNumeric>Top 5</Th>
            <Th isNumeric>Top 10</Th>
            <Th isNumeric>Segments</Th>
          </Tr>
        </Thead>
        <Tbody>
          {data.map((athlete) => (
            <Tr
              key={athlete.athleteId}
              _hover={{ bg: hoverBg }}
              cursor="pointer"
            >
              <Td>{getRankBadge(athlete.rank)}</Td>
              <Td>
                <HStack spacing={3}>
                  <Avatar
                    size="sm"
                    name={`${athlete.firstName} ${athlete.lastName}`}
                    src={athlete.profilePic}
                  />
                  <VStack align="start" spacing={0}>
                    <Text fontWeight="bold">
                      {athlete.firstName} {athlete.lastName}
                    </Text>
                    {athlete.rank <= 3 && (
                      <Badge colorScheme="orange" size="sm">
                        Top {athlete.rank}
                      </Badge>
                    )}
                  </VStack>
                </HStack>
              </Td>
              <Td isNumeric>
                <Text fontWeight="bold" fontSize="lg">
                  {athlete.totalPoints}
                </Text>
              </Td>
              <Td isNumeric>
                <HStack justifyContent="flex-end" spacing={1}>
                  <Icon as={FaFire} color="orange.500" />
                  <Text>{athlete.komsQoms}</Text>
                </HStack>
              </Td>
              <Td isNumeric>{athlete.top5Finishes}</Td>
              <Td isNumeric>{athlete.top10Finishes}</Td>
              <Td isNumeric>
                <Text>
                  {athlete.segmentsCompleted} / {mockLeaderboardStats.totalSegments}
                </Text>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Box>
  );
};
