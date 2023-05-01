import { FC } from "react"
import { HugelBoardProps } from "./HugelBoard"
import {
  Flex,
  Grid,
  Box,
  Avatar,
  Badge,
  Text,
  useColorModeValue,

} from '@chakra-ui/react'
import { HugelLeaderBoardActivity } from "../../api/typesGenerated"
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar"
import { DistanceToLocal, DistanceToLocalElevation } from "../../lib/Distance/Distance"

export const HugelBoardGallery: FC<HugelBoardProps> = ({
  data, error, isLoading, isFetched
}) => {

  return <>
    <FirstPlaceContainer>
      <GalleryCard activity={data?.activities[0]} position={1} />
    </FirstPlaceContainer>

    <SecondPlaceContainer>
      <GalleryCard activity={data?.activities[1]} position={2} />
      <GalleryCard activity={data?.activities[2]} position={3} />
    </SecondPlaceContainer>

    <RemainingPlacesContainer>
      {data?.activities.slice(3).map((activity, index) =>
        <GalleryCard activity={activity} position={index + 4} />
      )}
    </RemainingPlacesContainer>
  </>
}

const FirstPlaceContainer: React.FC<React.PropsWithChildren> = ({ children }) => {
  return <Flex m="1rem auto" maxW='xl' justifyContent='center'>
    {children}
  </Flex>
}
const SecondPlaceContainer: React.FC<React.PropsWithChildren> = ({ children }) => {
  return <Flex gap={3} maxW='2xl' m='0 auto 1rem' alignItems='center' flexDir={{ base: 'column', md: 'row' }}>
    {children}
  </Flex>
}
const RemainingPlacesContainer: React.FC<React.PropsWithChildren> = ({ children }) => {
  return <Grid gridTemplateColumns='repeat(auto-fit, minmax(300px, 1fr))' gap={3} maxWidth={'6xl'}>
    {children}
  </Grid>
}

const SkeletonGalleryCard: React.FC<{
  position: number
}> = ({}) => {
  return <Box w='100%' maxW='350px' h='300px' bg='transparent' />
}

// https://www.color-hex.com/color-palette/50061
const gradientsLight = [
  'rgba(212,175,55,1) 0%, rgba(227,227,227,1) 250%',
  'rgba(192,192,192,1) 0%, rgba(227,227,227,1) 250%',
  'rgba(205,127,50,1) 0%, rgba(227,227,227,1) 250%'
]
// Original
// 'linear-gradient(90deg, rgba(255,217,61,1) 0%, rgba(255,132,0,1) 100%)',
// 'linear-gradient(90deg, rgba(231,246,242,1) 0%, rgba(165,201,202,1) 100%)'

const gradientsDark = [
  'rgba(212,175,55,1) 0%, rgba(0,0,0,1) 200%',
  'rgba(192,192,192,1) 0%, rgba(0,0,0,1) 200%',
  'rgba(205,127,50,1) 0%, rgba(0,0,0,1) 200%'
]


const GalleryCard: React.FC<{
  activity?: HugelLeaderBoardActivity
  position: number
}> = ({ activity, position }) => {
  const gradients = useColorModeValue(gradientsLight, gradientsDark)
  const defaultColor = useColorModeValue(
    "rgba(252,76,2,1) 0%, rgba(0,0,0,1) 350%",
    "rgba(252,76,2,1) 0%, rgba(9,1,1,1) 200%"
  )
  const bgColor = `radial-gradient(circle, ${gradients[position - 1] || defaultColor})`
  const shadowColorRGB = useColorModeValue('0,0,0', '210,210,210')


  if (!activity) {
    // Return empty
    return <Box w='100%' maxW='350px' h='300px' bg={"transparent"} borderRadius={'1rem'} />
  }

  console.log({ activity })


  const {
    firstname,
    lastname,
    profile_pic_link,
    username,
  } = activity.athlete

  // 2022-11-27T15:42:54Z
  // Dates come over in UTC
  const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric' };
  const dateText = new Date(activity.activity_start_date).toLocaleDateString(undefined, options)
  const elapsed = `${Math.floor(activity.elapsed / 3600)}:${Math.floor(activity.elapsed / 60) % 60}:${activity.elapsed % 60}`
  const totalElapsed = `${Math.floor(activity.activity_elapsed_time / 3600)}:${Math.floor(activity.activity_elapsed_time / 60) % 60}`
  const elevation = `${Math.floor(DistanceToLocalElevation(activity.activity_total_elevation_gain) / 100) / 10}k`
  const showWatts = activity.efforts.every(effort => effort.average_watts > 0 && effort.device_watts)
  const avgWatts = Math.floor(activity.efforts.reduce((acc, effort) => acc + effort.average_watts * effort.elapsed_time, 0) / activity.elapsed).toString()

  return <Box w='100%' maxW='350px' h='300px' bg={bgColor} borderRadius={'1rem'}
    filter={`drop-shadow(2px 2px 2px rgba(${shadowColorRGB}, 0.25))`}
    transition={'all 0.25s ease-in-out'}
    _hover={{ filter: `drop-shadow(3px 3px 3px rgba(${shadowColorRGB}, 0.45))`, transform: 'translate(-5px, -5px)' }
    }>
    <Flex justifyContent={'space-between'}>
      <Flex p={3}>
        <AthleteAvatar
          firstName={firstname}
          lastName={lastname}
          athleteID={activity.athlete_id}
          profilePicLink={profile_pic_link}
          username={username}
          size="lg"
          styleProps={{
            mr: 3,
          }}
        />
        <Box>
          <Text fontWeight='bold' textAlign='left' >
            {firstname} {lastname}
          </Text>
          <Text fontSize='sm' fontFamily={'monospace'} opacity={.6} textAlign='left' >{dateText}</Text>
        </Box>
      </Flex>
      <Flex bg={'rgba(0,0,0,0.25)'} color={bgColor} p={'1.25rem'} maxHeight={'2.5rem'} borderRadius={'0 1rem 0 1rem'} alignItems={'center'} justifyContent={'center'}>
        <Text>
          {position}
        </Text>
      </Flex>

    </Flex>
    <Text fontWeight='bold'>{activity.activity_name}</Text>
    <Grid gridTemplateColumns='2fr 1fr' gap={3} p={4}>
      {/* Margins are not yet computed */}
      <StatBox stat={elapsed} label={"+00:00:05"} />
      <StatBox stat={Math.floor(DistanceToLocal(activity.activity_distance)).toString()} label={"miles"} />
    </Grid>
    <Grid gridTemplateColumns='1fr 1fr 1fr' gap={3} px={4}>
      <StatBox stat={totalElapsed} label={"total hours"} />
      <StatBox stat={showWatts ? avgWatts : "--"} label={"watts"} />
      <StatBox stat={elevation} label={"feet"} />
    </Grid>

  </Box>
}

const StatBox: React.FC<{
  stat?: string
  label?: string
}> = ({
  stat, label
}) => {
    return <Flex flexDir={'column'} justifyContent={'center'} alignItems={'center'} textAlign={'center'} bg={'antiquewhite'} h={'4rem'} borderRadius={3} color={'black'}>
      <Text fontWeight={700} fontFamily='monospace' fontSize='1rem'>{stat || "123"}</Text>
      <Text opacity={.5}>{label || "miles"}</Text>
    </Flex >
  }

// {Array.from({ length: 20 }).map(e =>
//   <GalleryCard />
// )}
