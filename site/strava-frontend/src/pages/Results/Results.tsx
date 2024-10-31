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
} from "@chakra-ui/react";
import { FC, ReactElement } from "react";
import { Link as RouteLink } from "react-router-dom";
import {
  StravaConnect,
  StravaConnectHref,
} from "../../components/Navbar/StravaConnect";
import { faFlagCheckered, faList } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { LandingCard } from "../Landing/Landing";
export const Results: FC<{}> = ({}) => {
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
            />
            <LandingCard
              heading={"Ride Das Hügel"}
              icon={<FontAwesomeIcon icon={faFlagCheckered} size="2x" />}
              description={
                "Challenge yourself on this epic ride and see how you compare to others!"
              }
              hrefText={"Route"}
              href={"/route/das-hugel"}
            />
            <LandingCard
              heading={"Explore Your Stats"}
              icon={<FontAwesomeIcon icon={faList} size="2x" />}
              description={"See results for all participating athletes."}
              hrefText={"Results"}
              href={"/hugelboard"}
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
