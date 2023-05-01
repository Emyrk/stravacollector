import { FC } from "react"
import { HugelBoardProps } from "./HugelBoard"
import {
  Flex,
  Grid,
  Box,
  Avatar,
  Badge,
  Text,

} from '@chakra-ui/react'
import { HugelLeaderBoardActivity } from "../../api/typesGenerated"
import { AthleteAvatar } from "../../components/AthleteAvatar/AthleteAvatar"

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
      {data?.activities.map((activity, index) =>
        <GalleryCard activity={activity} position={index + 3} />
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
  return <Flex gap={3} maxW='xl' m='0 auto 1rem' alignItems='center' flexDir={{ base: 'column', md: 'row' }}>
    {children}
  </Flex>
}
const RemainingPlacesContainer: React.FC<React.PropsWithChildren> = ({ children }) => {
  return <Grid gridTemplateColumns='repeat(auto-fit, minmax(300px, 1fr))' gap={3}>
    {children}
  </Grid>
}

const SkeletonGalleryCard: React.FC<{
  position: number
}> = ({}) => {
  return <Box w='100%' maxW='350px' h='300px' bg='brown' />
}

const GalleryCard: React.FC<{
  activity?: HugelLeaderBoardActivity
  position: number
}> = ({ activity, position }) => {
  if (!activity) {
    // Return empty
    return <Box w='100%' maxW='350px' h='300px' bg='brown' />
  }

  console.log({ activity })

  const {
    firstname,
    lastname,
    profile_pic_link,
    username,
  } = activity.athlete

  // 2022-11-27T15:42:54Z
  //   const dateText = new Date(activity.activity_start_date).format('YYYY-MM-DDTHH:MM:SS')
  const dateText = new Date(activity.activity_start_date).toISOString()


  return <Box w='100%' maxW='350px' h='300px' bg='brown' p={3} borderRadius={3}>
    <Flex>
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
        <Text fontWeight='bold' textAlign='left'>
          {firstname} {lastname}
        </Text>
        <Text fontSize='sm' fontFamily={'monospace'}>{activity.activity_start_date.split('T')[0]}</Text>
      </Box>
    </Flex>

  </Box>
}

// {Array.from({ length: 20 }).map(e =>
//   <GalleryCard />
// )}