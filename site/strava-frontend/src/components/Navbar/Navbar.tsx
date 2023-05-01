import {
  Box,
  Flex,
  Text,
  IconButton,
  Button,
  Stack,
  Collapse,
  Icon,
  Link,
  Popover,
  PopoverTrigger,
  PopoverContent,
  useColorModeValue,
  useBreakpointValue,
  useDisclosure,
  useToast,
  Avatar,
  Image,
  Container,
} from '@chakra-ui/react';
import {
  Link as RouteLink,
} from "react-router-dom";
import {
  HamburgerIcon,
  CloseIcon,
  ChevronDownIcon,
  ChevronRightIcon,
} from '@chakra-ui/icons';
import { StravaConnect } from './StravaConnect';
import { useAuthenticated } from '../../contexts/Authenticated';
import { getErrorMessage, getErrorDetail } from '../../api/rest';
import { useEffect } from 'react';
import { AthleteAvatar } from '../AthleteAvatar/AthleteAvatar';
import { AthleteAvatarDropdown } from './AthleteAvatarDropdown';

export default function Navbar() {
  const { isOpen, onToggle } = useDisclosure();
  const { authenticatedUser, isFetched: athleteFetched } = useAuthenticated()

  // TODO: Probably want to factor this toast better?
  // const toast = useToast()
  // const toastID = "authenticated-user-toast"
  // useEffect(() => {
  //   if (fetchError) {
  //     const title = getErrorMessage(fetchError, "Authentication Error")
  //     const description = getErrorDetail(fetchError)
  //     toast({
  //       id: toastID,
  //       title: <>{title}</>,
  //       description: <>{description ? description : "Unkown error getting authenticated user"}</>,
  //       position: "bottom-right",
  //       isClosable: true,
  //       status: "error",
  //     })
  //   }
  // }, [toast, authenticatedUser, fetchError]);

  const connectURL = "oauth2/callback?redirect=" +
    (window.location.pathname ? encodeURIComponent(window.location.pathname) : "/")
  return (
    <Box>
      <Flex
        bg={useColorModeValue('white', 'gray.800')}
        color={useColorModeValue('gray.600', 'white')}
        minH={'60px'}
        py={{ base: 2 }}
        px={{ base: 4 }}
        borderBottom={1}
        borderStyle={'solid'}
        borderColor={useColorModeValue('gray.200', 'gray.900')}
        align={'center'}>
        <Flex
          flex={{ base: 1, md: 'auto' }}
          ml={{ base: -2 }}
          display={{ base: 'flex', md: 'none' }}>
          <IconButton
            onClick={onToggle}
            icon={
              isOpen ? <CloseIcon w={3} h={3} /> : <HamburgerIcon w={5} h={5} />
            }
            variant={'ghost'}
            aria-label={'Toggle Navigation'}
          />
        </Flex>
        <Flex flex={{ base: 1 }} justify={{ base: 'center', md: 'start' }}>
          <Box>
            <RouteLink to="/">
              {/* https://chakra-ui.com/docs/components/image/usage */}
              <Image maxHeight={"80px"} src="/logos/LogoTypeColorSquare.png" alt="Hugel Ranker" display={{ base: 'block', md: 'none' }} />
              <Image maxHeight={"80px"} src="/logos/LogoTypeColor.png" alt="Hugel Ranker" display={{ base: 'none', md: 'block' }} />
            </RouteLink>
          </Box>

          <Flex display={{ base: 'none', md: 'flex' }} ml={10}>
            <DesktopNav />
          </Flex>
        </Flex>

        <Stack
          flex={{ base: 1, md: 0 }}
          justify={'flex-end'}
          direction={'row'}
          spacing={6}>
          {authenticatedUser ?
            <AthleteAvatarDropdown athlete={authenticatedUser} />
            :
            <Link href={connectURL}>
              <IconButton
                aria-label={"strava sign in"}
                icon={<StravaConnect />}
                bg="transparent"
              />
            </Link>
          }
        </Stack>
      </Flex>

      <Collapse in={isOpen} animateOpacity>
        <MobileNav />
      </Collapse>
    </Box >
  );
}

const DesktopNav = () => {
  const linkColor = useColorModeValue('gray.600', 'gray.200');
  const linkHoverColor = useColorModeValue('gray.800', 'white');
  const popoverContentBgColor = useColorModeValue('white', 'gray.800');

  return (
    <Stack direction={'row'} spacing={4}>
      {NAV_ITEMS.map((navItem) => (
        <Box key={navItem.label}>
          <Popover trigger={'hover'} placement={'bottom-start'}>
            <PopoverTrigger>
              <Container
                p={2}
                fontSize={'sm'}
                fontWeight={500}
                color={linkColor}
                _hover={{
                  textDecoration: 'none',
                  color: linkHoverColor,
                }}>
                <RouteLink to={navItem.href ?? '#'}>
                  {navItem.label}
                </RouteLink>
              </Container>
            </PopoverTrigger>

            {navItem.children && (
              <PopoverContent
                border={0}
                boxShadow={'xl'}
                bg={popoverContentBgColor}
                p={4}
                rounded={'xl'}
                minW={'sm'}>
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
      role={'group'}
      display={'block'}
      p={2}
      rounded={'md'}
      _hover={{ bg: useColorModeValue('pink.50', 'gray.900') }}
    >
      <RouteLink
        to={href || '#'}
      >
        <Stack direction={'row'} align={'center'}>
          <Box>
            <Text
              transition={'all .3s ease'}
              _groupHover={{ color: 'pink.400' }}
              fontWeight={500}>
              {label}
            </Text>
            <Text fontSize={'sm'}>{subLabel}</Text>
          </Box>
          <Flex
            transition={'all .3s ease'}
            transform={'translateX(-10px)'}
            opacity={0}
            _groupHover={{ opacity: '100%', transform: 'translateX(0)' }}
            justify={'flex-end'}
            align={'center'}
            flex={1}>
            <Icon color={'pink.400'} w={5} h={5} as={ChevronRightIcon} />
          </Flex>
        </Stack>
      </RouteLink>
    </Container >
  );
};

const MobileNav = () => {
  return (
    <Stack
      bg={useColorModeValue('white', 'gray.800')}
      p={4}
      display={{ md: 'none' }}>
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
        href={href ?? '#'}
        justify={'space-between'}
        align={'center'}
        _hover={{
          textDecoration: 'none',
        }}>
        <Text
          fontWeight={600}
          color={useColorModeValue('gray.600', 'gray.200')}>
          {label}
        </Text>
        {children && (
          <Icon
            as={ChevronDownIcon}
            transition={'all .25s ease-in-out'}
            transform={isOpen ? 'rotate(180deg)' : ''}
            w={6}
            h={6}
          />
        )}
      </Flex>

      <Collapse in={isOpen} animateOpacity style={{ marginTop: '0!important' }}>
        <Stack
          mt={2}
          pl={4}
          borderLeft={1}
          borderStyle={'solid'}
          borderColor={useColorModeValue('gray.200', 'gray.700')}
          align={'start'}>
          {children &&
            children.map((child) => (
              <Container py={2} key={child.label}>
                <RouteLink to={child.href || '#'}>
                  {child.label}
                </RouteLink>
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
    label: 'Das Hugel',
    children: [
      {
        label: 'All Hugels',
        subLabel: 'See how you stack up in the Das Hugel Leaderboard',
        href: '/hugelboard',
      },
      {
        label: 'All Hugel Super Scores',
        subLabel: 'Have not done a Das Hugel?',
        href: '#',
      },
    ],
  },
  // {
  //   label: 'Find Work',
  //   children: [
  //     {
  //       label: 'Job Board',
  //       subLabel: 'Find your dream design job',
  //       href: '#',
  //     },
  //     {
  //       label: 'Freelance Projects',
  //       subLabel: 'An exclusive list for contract work',
  //       href: '#',
  //     },
  //   ],
  // },
  // {
  //   label: 'Learn Design',
  //   href: '#',
  // },
  // {
  //   label: 'Hire Designers',
  //   href: '#',
  // },
];