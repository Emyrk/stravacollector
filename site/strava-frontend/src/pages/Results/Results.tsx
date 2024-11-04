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
import {
  faList,
  faMountain,
  faMountainCity,
  faMound,
} from "@fortawesome/free-solid-svg-icons";
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
          py={{ base: 12, md: 16 }}
        >
          <Heading
            fontWeight={600}
            fontSize={{ base: "3xl", sm: "4xl", md: "6xl" }}
            lineHeight={"110%"}
          >
            Tour Das Hügel Results
          </Heading>
          <Text color={"gray.300"} maxW={"3xl"}>
            Congratulations on finishing the Tour Das Hügel! Explore our
            leaderboards below to see how you stack up against other athletes.
          </Text>
          <Stack spacing={6} direction={{ base: "column", md: "row" }}>
            <LandingCard
              heading={"2024 Das Hügel"}
              icon={<Image src="/img/icons/Results2.png" />}
              description={
                "See how you did on the Hügel this year and find out who won our superlatives."
              }
              hrefText={"Results"}
              href={"/hugelboard/2024"}
            />
            <LandingCard
              heading={"2024 Hügel Lite"}
              icon={<Image src="/img/icons/Results1.png" />}
              description={
                "Full Hügel not in the cards this year? See how you did on the first 40 miles!"
              }
              hrefText={"Results"}
              href={"/hugelboard/2024"}
            />
            {/* <LandingCard
              heading={"Super Hügel"}
              icon={<Image src="/img/icons/Results3.png" />}
              description={
                "Improve your time on each individual segment to become the ultimate Hügeler."
              }
              hrefText={"Results"}
              href={"/superhugelboard"}
            /> */}
          </Stack>
          <Stack
            textAlign={"center"}
            align={"center"}
            spacing={{ base: 8, md: 10 }}
            py={{ base: 8, md: 12 }}
          >
            <Heading
              fontWeight={400}
              fontSize={{ base: "1xl", sm: "2xl", md: "4xl" }}
              lineHeight={"110%"}
            >
              Tour Das Hügel Hall of Fame
            </Heading>

            <Heading
              fontWeight={400}
              fontSize={{ base: "xl" }}
              lineHeight={"110%"}
            >
              <Link as={RouteLink} to={"/hugelboard/2023"}>
                <Button variant={"link"} color="brand.stravaOrange" size={"xl"}>
                  2023 Das Hügel
                </Button>
              </Link>
            </Heading>
          </Stack>
        </Stack>
      </Container>
    </Flex>
  );
};
