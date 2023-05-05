import { Link, Flex, Text, Box, Grid, VStack, Code, Heading, useColorModeValue } from "@chakra-ui/react"
import { FC } from "react"
import { Logo } from "../../Logo"
import {
  Link as RouteLink
} from "react-router-dom";

const bgColors = {
  light: 'antiquewhite',
  dark: 'saddlebrown',
}

export const Landing: FC = () => {
  const bg = useColorModeValue(bgColors.light, bgColors.dark)

  return <Box textAlign="center" fontSize="xl" maxW={'7xl'} m={'1rem auto'}>
    <Heading>Welcome to Tour Das Hügel!</Heading>
    <Box bg={bg} p={5} borderRadius={3} m={'0 auto'} maxW='2xl' w='100%'>
      <Text>
        Austin, Texas is home to the notorious Tour das Hügel — a challenging bike route boasting over 100 miles of treacherous
        hills and climbs totaling over 13,000 feet. Unfortunately, riders attempting this ride often encounter dangerous stop signs
        that can prompt them to take unnecessary risks in order to get better times. To help cyclists stay safe and get the most out
        of their Tour das Hügel experience, our site provides an alternative solution by calculating the best times for riders on the
        famous climbs, allowing them to avoid being penalized for stop signs or traffic jams.
      </Text>
      <Text>In Austin Texas exists this notorious biking route! This route is filled with dangerous stop-signs which push racers to take risks to get the best times. This site instead calculates riders best times on the climbs to avoid penalizations for stop-signs/traffic.</Text>
    </Box>

    <Flex w='100%' maxW='4xl' m='0 auto' justifyContent='space-between' p={3} gap={'1rem'} flexDir={{ base: 'column', md: 'row' }}>
      <Card header={'Connect With Strava'} text='Link your Strava account to our site to get started and unlock exclusive benefits.' />
      <Card header={'Wait to Sync'} text={`Our API may be slow to load due to our current alpha stage. We are actively working on improving its speed and performance. Please come back after signing in and your data should be available shortly. Thank you for your patience!`} />
      <Card header={'Climb the Leaderboard'} text='Track your scores against others and challenge yourself to be the best on the Hügel!' />
    </Flex>
  </Box >
}

// Connect with Strava
// Wait for your activities to sync
// See how you stack up!

const Card: React.FC<{ header: string, text: string }> = ({ header, text }) => {

  const bg = useColorModeValue(bgColors.light, bgColors.dark)
  return <Box bg={bg} p={3} w='100%' borderRadius={3} textAlign={'left'}>
    <Text fontWeight={700}>{header}</Text>
    <Text fontSize='0.9rem' opacity={0.6}>{text}</Text>
  </Box>
}