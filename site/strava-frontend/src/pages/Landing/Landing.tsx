import { Link, Flex, Text, Box, Grid, VStack, Code, Heading, useColorModeValue } from "@chakra-ui/react"
import { FC } from "react"
import { ColorModeSwitcher } from "../../ColorModeSwitcher"
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
    <ColorModeSwitcher justifySelf="flex-end" />
    <Heading>Welcome to Tour Das Hügel!</Heading>
    <Box bg={bg} p={5} borderRadius={3} m={'0 auto'} maxW='2xl' w='100%'>
      <Text>In Austin Texas exists this notorious biking route! This route is filled with dangerous stop-signs which push racers to take risks to get the best times. This site instead calculates riders best times on the climbs to avoid penalizations for stop-signs/traffic.</Text>
    </Box>

    <Flex w='100%' maxW='4xl' m='0 auto' justifyContent='space-between' p={3} gap={'1rem'} flexDir={{ base: 'column', md: 'row' }}>
      <Card header={'Connect with Strava'} text='Connect your strava account to this site' />
      <Card header={'Wait for Sync'} text={`We're currently in alpha so our api takes awhile to load your data. Give us some time and come back after signing in.`} />
      <Card header={'Climb the leaderboard'} text='Track your scores against others against the Hügel!' />
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