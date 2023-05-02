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
} from '@chakra-ui/react';
import {
  Link as RouteLink,
} from "react-router-dom";
import {
  ChevronDownIcon,
  ChevronRightIcon,
} from '@chakra-ui/icons';
import { StravaConnect } from './StravaConnect';
import { useAuthenticated } from '../../contexts/Authenticated';
import { getErrorMessage, getErrorDetail } from '../../api/rest';
import { useEffect } from 'react';
import { AthleteAvatar } from '../AthleteAvatar/AthleteAvatar';
import { AthleteAvatarDropdown } from './AthleteAvatarDropdown';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrophy } from '@fortawesome/free-solid-svg-icons'
import { ColorModeSwitcher } from '../ColorModeSwitcher/ColorModeSwitcher';

const Navbar: React.FC = () => {
  const { isOpen, onToggle } = useDisclosure();
  const theme = useTheme()
  console.log(theme.colors)

  return <>
    <Flex w='100%' maxW={'7xl'} m={'1rem auto 0'} justifyContent='space-between' alignItems={'center'} p={3} pb={0}>
      <Box>
        <RouteLink to="/">
          {/* https://chakra-ui.com/docs/components/image/usage */}
          <Image maxHeight={"80px"} src="/logos/LogoTypeColorSquare.png" alt="Hugel Ranker" display={{ base: 'block', md: 'none' }} />
          <Image maxHeight={"80px"} src="/logos/LogoTypeColor.png" alt="Hugel Ranker" display={{ base: 'none', md: 'block' }} />
        </RouteLink>
      </Box>
      <ColorModeSwitcher justifySelf="flex-end" />
      <Flex alignItems={'center'} gap={2}>
        <DesktopNav display={{ base: 'none', md: 'block' }} />
        <StravaConnect />
      </Flex>
    </Flex>
    <MobileNav display={{ base: 'block', md: 'none' }} />
  </>
}

export default Navbar


const DesktopNav: React.FC<{ display: { base: string, md: string } }> = ({ display }) => {
  const linkColor = useColorModeValue('gray.600', 'gray.200');
  const linkHoverColor = useColorModeValue('gray.800', 'white');
  const popoverContentBgColor = useColorModeValue('white', 'gray.800');

  return (
    <Stack direction={'row'} spacing={4} display={display}>
      {NAV_ITEMS.map((navItem, index) => (
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
                  <Tag p={3} display={'flex'} gap={2}>
                    <FontAwesomeIcon icon={faTrophy} />
                    {navItem.label}
                  </Tag>
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

const MobileNav: React.FC<{ display: { base: string, md: string } }> = ({ display }) => {
  return (
    <Stack
      bg={useColorModeValue('white', 'gray.800')}
      p={4}
      display={display}>
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
        alignItems={'center'}
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