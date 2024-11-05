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
      {NAV_ITEMS.map((navItem, index) => {
        return (
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
                      <DesktopSubNav key={child.label} {...child} depth={0} />
                    ))}
                  </Stack>
                </PopoverContent>
              )}
            </Popover>
          </Box>
        );
      })}
    </Stack>
  );
};

const DesktopSubNav = ({
  label,
  href,
  subLabel,
  children,
  depth,
}: NavItem & { depth: number }) => {
  return (
    <>
      <Container
        p={2}
        pl={depth > 0 ? depth * 30 + "px" : 2}
        role={"group"}
        display={"block"}
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
      {children &&
        children.map((child) => (
          <DesktopSubNav key={child.label} depth={depth + 1} {...child} />
        ))}
    </>
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

interface NavItem {
  label: string;
  subLabel?: string;
  children?: Array<NavItem>;
  href?: string;
}

const NAV_ITEMS: Array<NavItem> = [
  {
    label: "Das Hügel",
    children: [
      // {
      //   label: "Das Hügel Super Scores",
      //   subLabel: "Your best segments across all your rides.",
      //   href: "/superhugelboard",
      // },
      {
        label: "Das Hügel Route",
        subLabel: "Plan your ride with the required segments",
        href: "/route/das-hugel",
      },
      {
        label: "Das Hügel Results",
        subLabel: "See how you did on the Tour Das Hügel",
        href: "/results",
        children: [
          {
            label: "2024 Das Hügel Results",
            href: "/hugelboard/2024",
          },
          {
            label: "2024 Das Hügel Lite Results",
            // subLabel: "See how you did on the Tour Das Hügel",
            href: "/hugelboard/2024?lite=true",
          },
          {
            label: "Other Results",
            // subLabel: "See how you did on the Tour Das Hügel",
            href: "/results",
          },
        ],
      },
    ],
  },
];
