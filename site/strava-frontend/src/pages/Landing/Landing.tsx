import {
  Link,
  Flex,
  Text,
  Box,
  Heading,
  useColorModeValue,
  Button,
  Container,
  Stack,
  useStyleConfig,
  Image,
} from "@chakra-ui/react";
import { FC, ReactElement } from "react";
import { Link as RouteLink } from "react-router-dom";
import {
  StravaConnect,
  StravaConnectHref,
} from "../../components/Navbar/StravaConnect";
import { faFlagCheckered, faList } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const bgColors = {
  light: "antiquewhite",
  dark: "saddlebrown",
};

export const Landing: FC = () => {
  return (
    <Flex
      w={"100vw"}
      // h={"100vh"}
      // backgroundImage={"url(/img/hill-rider.svg)"}
      backgroundSize={"cover"}
      backgroundPosition={"center center"}
    >
      <Container maxW={"5xl"}>
        <Stack
          textAlign={"center"}
          align={"center"}
          spacing={{ base: 8, md: 10 }}
          py={{ base: 20, md: 28 }}
        >
          <Heading
            fontWeight={600}
            fontSize={{ base: "3xl", sm: "4xl", md: "6xl" }}
            lineHeight={"110%"}
          >
            Welcome to the{" "}
            <Text as={"span"} color={"brand.stravaOrange"}>
              Tour Das Hügel
            </Text>
          </Heading>
          <Text color={"gray.300"} maxW={"3xl"}>
            Austin, Texas is home to the notorious Tour das Hügel, a challenging
            bike route with over 110 miles of strenuous climbs that reach a
            total elevation gain of over 11,500 feet. To help cyclists stay safe
            and get the most out of their experience, we offer a timing system
            that exclusively records times on the iconic climbs (segments).
            Because this is an unsanctioned event held on open roads, our system
            ensures riders are not penalized for things like stop signs or
            traffic congestion.
          </Text>
          <Stack spacing={6} direction={{ base: "column", md: "row" }}>
            <LandingCard
              heading={"Connect to Strava"}
              icon={<StravaConnect useSquareLogo={true} />}
              description={
                "Link your Strava account to our site to explore your stats."
              }
              hrefText={"Connect"}
              href={StravaConnectHref()}
              hrefRealLink={true}
            />
            <LandingCard
              heading={"Ride Das Hügel"}
              icon={<Image src="/img/icons/Home2.png" />}
              description={
                "Challenge yourself on this epic ride and see how you compare to others!"
              }
              hrefText={"Route"}
              href={"/route/das-hugel"}
            />
            <LandingCard
              heading={"Explore Your Stats"}
              icon={<Image src="/img/icons/Home3.png" />}
              description={
                "Check out the results for all athletes who have participated."
              }
              hrefText={"Results"}
              href={"/results"}
            />

            {/* <LandingCard
              heading={"Wait to Sync"}
              icon={<FontAwesomeIcon icon={faSpinner} size="2x" />}
              description={
                "Our API may be slow to load due to our current alpha stage. We are actively working on improving its speed and performance. Please come back after signing in and your data should be available shortly. Thank you for your patience!"
              }
              href={"#"}
            /> */}
            {/* <Button
              rounded={"full"}
              px={6}
              colorScheme={"orange"}
              bg={"orange.400"}
              _hover={{ bg: "orange.500" }}
            >
              Get started
            </Button>
            <Button rounded={"full"} px={6}>
              Learn more
            </Button> */}
          </Stack>
          <Flex w={"full"}>
            {/* <chakra.img
            src="/img/hill-rider.svg"
            height={{ sm: "24rem", lg: "28rem" }}
            mt={{ base: 12, sm: 16 }}
          /> */}
          </Flex>
        </Stack>
      </Container>
    </Flex>
  );
};

interface LandingCardProps {
  heading: string;
  description: string;
  icon: ReactElement;
  href: string;
  hrefRealLink?: boolean;
  hrefText: string;
}

export const LandingCard = ({
  heading,
  description,
  icon,
  href,
  hrefText,
  hrefRealLink,
}: LandingCardProps) => {
  const styles = useStyleConfig("Box", { variant: "responsiveCard" });
  let link = <></>;
  if (href !== "") {
    let props = { to: href } as Record<string, any>;
    if (hrefRealLink) {
      props = { href: href };
    }
    link = (
      <Link as={hrefRealLink ? Link : RouteLink} {...props}>
        <Button variant={"link"} color="brand.stravaOrange" size={"sm"}>
          {hrefText}
        </Button>
      </Link>
    );
  }

  return (
    <Box
      __css={styles}
      maxW={{ base: "full", md: "340px" }}
      w={"full"}
      borderWidth="1px"
      borderRadius="lg"
      overflow="hidden"
      p={5}
    >
      <Flex direction="column" gap={"1em"}>
        <Flex
          direction={"row"}
          align={"center"}
          gap={"1em"}
          width={"100%"}
          justifyContent={"center"}
        >
          <Flex w={"48px"} h={"48px"} align={"center"} justify={"center"}>
            {icon}
          </Flex>
          <Heading size="md">{heading}</Heading>
        </Flex>
        <Box mt={0.5}>
          <Text mt={1} fontSize={"sm"}>
            {description}
          </Text>
        </Box>
        {link}
      </Flex>
    </Box>
  );
};

export const Landing2: FC = () => {
  const bg = useColorModeValue(bgColors.light, bgColors.dark);

  return (
    <Box textAlign="center" fontSize="xl" maxW={"7xl"} m={"1rem auto"}>
      <Heading>Welcome to Tour Das Hügel!</Heading>
      <Box bg={bg} p={5} borderRadius={3} m={"0 auto"} maxW="2xl" w="100%">
        <Text>
          Austin, Texas is home to the notorious Tour das Hügel — a challenging
          bike route boasting over 100 miles of treacherous hills and climbs
          totaling over 13,000 feet. The event is unsanctioned and is done on
          open roads. To help cyclists stay safe and get the most out of their
          Tour das Hügel experience, our site provides a Gran Fondo type timing
          system. The solution only times riders on the famous climbs, allowing
          them to avoid being penalized for stop signs or traffic jams.
        </Text>
      </Box>

      <Flex
        w="100%"
        maxW="4xl"
        m="0 auto"
        justifyContent="space-between"
        p={3}
        gap={"1rem"}
        flexDir={{ base: "column", md: "row" }}
      >
        <Card
          header={"Connect With Strava"}
          text="Link your Strava account to our site to get started and unlock exclusive benefits."
        />
        <Card
          header={"Wait to Sync"}
          text={`Our API may be slow to load due to our current alpha stage. We are actively working on improving its speed and performance. Please come back after signing in and your data should be available shortly. Thank you for your patience!`}
        />
        <Card
          header={"Climb the Leaderboard"}
          text="Track your scores against others and challenge yourself to be the best on the Hügel!"
        />
      </Flex>
    </Box>
  );
};

// Connect with Strava
// Wait for your activities to sync
// See how you stack up!

const Card: React.FC<{ header: string; text: string }> = ({ header, text }) => {
  const bg = useColorModeValue(bgColors.light, bgColors.dark);
  return (
    <Box bg={bg} p={3} w="100%" borderRadius={3} textAlign={"left"}>
      <Text fontWeight={700}>{header}</Text>
      <Text fontSize="0.9rem" opacity={0.6}>
        {text}
      </Text>
    </Box>
  );
};
