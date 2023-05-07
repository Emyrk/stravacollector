import { FC } from "react"
import { AthleteSummary } from "../../api/typesGenerated"
import { Text, Menu, MenuDivider, MenuList, MenuItem, MenuButton, Flex, Button, useTheme, Center, Container, Link } from "@chakra-ui/react"
import { AthleteAvatar } from "../AthleteAvatar/AthleteAvatar"
import { ChevronDownIcon, ChevronUpIcon } from "@chakra-ui/icons"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowRightFromBracket, faGear } from '@fortawesome/free-solid-svg-icons'
import {
  Link as RouteLink,
} from "react-router-dom";


export const AthleteAvatarDropdown: FC<{
  athlete: AthleteSummary
}> = ({ athlete }) => {

  return <Menu  placement={'bottom-end'}>
    {({ isOpen }) => (<>
      <MenuButton
        rounded={'full'}
        variant={'link'}
        as={Button}
        // rightIcon={isOpen ? <ChevronUpIcon /> : <ChevronDownIcon />}
      >
        <AthleteAvatar
          firstName={athlete.firstname}
          lastName={athlete.lastname}
          athleteID={athlete.athlete_id}
          username={athlete.username}
          profilePicLink={athlete.profile_pic_link}
          hugelCount={athlete.hugel_count}
          size="lg"
        />
      </MenuButton>
      <MenuList alignItems={'center'}>
        <Center>
          <AthleteAvatar
            firstName={athlete.firstname}
            lastName={athlete.lastname}
            athleteID={athlete.athlete_id}
            username={athlete.username}
            profilePicLink={athlete.profile_pic_link}
            hugelCount={athlete.hugel_count}
            size="2xl"
          />
        </Center>
        <br />
        <Center>
          <p>{athlete.firstname} {athlete.lastname}</p>
        </Center>
        <MenuDivider />
        <MenuItem>
          <FontAwesomeIcon icon={faGear} />
          <Container paddingLeft={"10px"}>Settings</Container>
        </MenuItem>
        <Link href="/logout" style={{ textDecoration: 'none' }}>
          <MenuItem>
            <FontAwesomeIcon icon={faArrowRightFromBracket} />
            <Text as="span" paddingLeft={"10px"}>Logout</Text>
          </MenuItem>
        </Link>
      </MenuList>
    </>
    )}
  </Menu >
}

