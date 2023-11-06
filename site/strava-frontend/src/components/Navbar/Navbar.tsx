import {
  Box,
  Flex,
  Text,
  Stack,
  Collapse,
  Icon,
  Link,
  Popover,
  PopoverTrigger,
  PopoverContent,
  useColorModeValue,
  useDisclosure,
  Image,
  Container,
  Tag,
  useTheme,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";
import { Link as RouteLink } from "react-router-dom";
import {
  AddIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  EditIcon,
  ExternalLinkIcon,
  HamburgerIcon,
  RepeatIcon,
} from "@chakra-ui/icons";
import { StravaConnectOrUser } from "./StravaConnect";
import { useAuthenticated } from "../../contexts/Authenticated";
import { getErrorMessage, getErrorDetail } from "../../api/rest";
import React, { FC, useEffect } from "react";
import { AthleteAvatar } from "../AthleteAvatar/AthleteAvatar";
import { AthleteAvatarDropdown } from "./AthleteAvatarDropdown";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrophy, faBars } from "@fortawesome/free-solid-svg-icons";
import { ColorModeSwitcher } from "../ColorModeSwitcher/ColorModeSwitcher";

const Navbar: React.FC = () => {
  const { isOpen, onToggle } = useDisclosure();

  return (
    <>
      <Flex
        w="100%"
        maxW={"7xl"}
        m={"1rem auto 0"}
        justifyContent="space-between"
        alignItems={"center"}
        p={3}
        pb={0}
      >
        <Box>
          <RouteLink to="/">
            {/* https://chakra-ui.com/docs/components/image/usage */}
            <Image
              maxHeight={"80px"}
              src="/logos/Logomark.png"
              alt="Hugel Ranker"
              display={{ base: "block", md: "none" }}
            />
            <Image
              maxHeight={"80px"}
              src="/logos/LongDasHugelWhite.png"
              alt="Hugel Ranker"
              display={{ base: "none", md: "block" }}
            />
          </RouteLink>
        </Box>

        <Flex alignItems={"center"} gap={2} marginLeft={"auto"}>
          <DesktopNav display={{ base: "none", md: "block" }} />
          <StravaConnectOrUser />
          <MobileNav2 display={{ base: "block", md: "none" }} />
        </Flex>
      </Flex>
      {/* <MobileNav display={{ base: "block", md: "none" }} /> */}
    </>
  );
};

export default Navbar;

const DesktopNav: React.FC<{ display: { base: string; md: string } }> = ({
  display,
}) => {
  return (
    <Stack direction={"row"} spacing={4} display={display}>
      {NAV_ITEMS.map((navItem, index) => (
        <Box key={navItem.label}>
          <Popover trigger={"hover"} placement={"bottom-end"}>
            <PopoverTrigger>
              <Container p={2} fontSize={"md"} fontWeight={500}>
                <RouteLink to={navItem.href ?? "#"}>
                  <Tag p={3} display={"flex"} gap={2} borderRadius={"2px"}>
                    <FontAwesomeIcon icon={faTrophy} />
                    <Text>{navItem.label}</Text>
                  </Tag>
                </RouteLink>
              </Container>
            </PopoverTrigger>

            {navItem.children && (
              <PopoverContent
                border={0}
                boxShadow={"xl"}
                p={4}
                rounded={"xl"}
                minW={"sm"}
              >
                <Stack>
                  {navItem.children.map((child) => (
                    <DesktopSubNav key={child.label} {...child} />
                  ))}
                </Stack>
              </PopoverContent>
            )}
          </Popover>
        </Box>
      ))}
    </Stack>
  );
};

const DesktopSubNav = ({ label, href, subLabel }: NavItem) => {
  return (
    <Container
      role={"group"}
      display={"block"}
      p={2}
      rounded={"md"}
      _hover={{ bg: useColorModeValue("pink.50", "gray.900") }}
    >
      <RouteLink to={href || "#"}>
        <Stack direction={"row"} align={"center"}>
          <Box>
            <Text
              transition={"all .3s ease"}
              _groupHover={{ color: "brand.stravaOrange" }}
              fontWeight={500}
            >
              {label}
            </Text>
            <Text fontSize={"sm"}>{subLabel}</Text>
          </Box>
          <Flex
            transition={"all .3s ease"}
            transform={"translateX(-10px)"}
            opacity={0}
            _groupHover={{ opacity: "100%", transform: "translateX(0)" }}
            justify={"flex-end"}
            align={"center"}
            flex={1}
          >
            <Icon
              color={"brand.stravaOrange"}
              w={5}
              h={5}
              as={ChevronRightIcon}
            />
          </Flex>
        </Stack>
      </RouteLink>
    </Container>
  );
};

export const MobileNav2: React.FC<{
  display: { base: string; md: string };
}> = ({ display }) => {
  const bugerColor = useColorModeValue(
    "brand.stravaOrange",
    "colors.alphaWhite.800"
  );

  // Hugel links + Light/Dark toggle
  return (
    <Box display={display}>
      <Menu>
        <MenuButton
          color={bugerColor}
          variant="outline"
          as={IconButton}
          aria-label="Options"
          icon={<HamburgerIcon />}
        />
        <MenuList>
          {NAV_ITEMS.map((navItem) => {
            return <MobileNav2Item key={navItem.label} item={navItem} />;
          })}
        </MenuList>
      </Menu>
    </Box>
  );
};

const MobileNav2Item: React.FC<{ item: NavItem }> = ({ item }) => {
  if (item.children) {
    return (
      <>
        {item.children.map((child) => {
          return <MobileNav2Item key={child.label} item={child} />;
        })}
      </>
    );
  }
  return (
    <RouteLink to={item.href || ""}>
      <MenuItem>{item.label}</MenuItem>
    </RouteLink>
  );
};

const MobileNav: React.FC<{ display: { base: string; md: string } }> = ({
  display,
}) => {
  return (
    <Stack bg={useColorModeValue("white", "gray.800")} p={4} display={display}>
      {NAV_ITEMS.map((navItem) => (
        <MobileNavItem key={navItem.label} {...navItem} />
      ))}
    </Stack>
  );
};

const MobileNavItem = ({ label, children, href }: NavItem) => {
  const { isOpen, onToggle } = useDisclosure();

  return (
    <Stack spacing={4} onClick={children && onToggle}>
      <Flex
        py={2}
        as={Link}
        href={href ?? "#"}
        justify={"space-between"}
        alignItems={"center"}
        _hover={{
          textDecoration: "none",
        }}
      >
        <Text
          fontWeight={600}
          color={useColorModeValue("gray.600", "gray.200")}
        >
          {label}
        </Text>
        {children && (
          <Icon
            as={ChevronDownIcon}
            transition={"all .25s ease-in-out"}
            transform={isOpen ? "rotate(180deg)" : ""}
            w={6}
            h={6}
          />
        )}
      </Flex>

      <Collapse in={isOpen} animateOpacity style={{ marginTop: "0!important" }}>
        <Stack
          mt={2}
          pl={4}
          borderLeft={1}
          borderStyle={"solid"}
          borderColor={useColorModeValue("gray.200", "gray.700")}
          align={"start"}
        >
          {children &&
            children.map((child) => (
              <Container py={2} key={child.label}>
                <RouteLink to={child.href || "#"}>{child.label}</RouteLink>
              </Container>
            ))}
        </Stack>
      </Collapse>
    </Stack>
  );
};

interface NavItem {
  label: string;
  subLabel?: string;
  children?: Array<NavItem>;
  href?: string;
}

const NAV_ITEMS: Array<NavItem> = [
  {
    label: "Das Hugel",
    children: [
      {
        label: "Das Hugel Results",
        subLabel: "See how you did on the 2023 Tour Das Hügel",
        href: "/hugelboard",
      },
      // {
      //   label: "Das Hugel Super Scores",
      //   subLabel: "Your best segments across all your rides.",
      //   href: "/superhugelboard",
      // },
      {
        label: "Das Hugel Route",
        subLabel: "Plan your ride with the required segments",
        href: "/route/das-hugel",
      },
    ],
  },
];
