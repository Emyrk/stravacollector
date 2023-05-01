import { FC } from "react"
import { AthleteSummary } from "../../api/typesGenerated"
import { Text, Menu, MenuList, MenuItem, MenuButton, Button, useTheme, Container, Link } from "@chakra-ui/react"
import { AthleteAvatar } from "../AthleteAvatar/AthleteAvatar"
import { ChevronDownIcon } from "@chakra-ui/icons"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowRightFromBracket, faGear } from '@fortawesome/free-solid-svg-icons'
import {
  Link as RouteLink,
} from "react-router-dom";


export const AthleteAvatarDropdown: FC<{
  athlete: AthleteSummary
}> = ({ athlete }) => {

  return <Menu>
    {/* as={Button} rightIcon={<ChevronDownIcon />} */}
    <MenuButton>
      <AthleteAvatar
        firstName={athlete.firstname}
        lastName={athlete.lastname}
        athleteID={athlete.athlete_id}
        username={athlete.username}
        profilePicLink={athlete.profile_pic_link}
        size="lg"
      />
    </MenuButton>
    <MenuList>
      <Link href="/logout" style={{ textDecoration: 'none' }}>
        <MenuItem>
          <FontAwesomeIcon icon={faArrowRightFromBracket} />
          <Text as="span" paddingLeft={"10px"}>Logout</Text>
        </MenuItem>
      </Link>
      <MenuItem>
        <FontAwesomeIcon icon={faGear} />
        <Container paddingLeft={"10px"}>Settings</Container>
      </MenuItem>
    </MenuList>
  </Menu >
}

