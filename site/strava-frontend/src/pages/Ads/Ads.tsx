import { Container, Heading, Text, VStack } from "@chakra-ui/react";
import { FC } from "react";
import { AdSenseAd } from "../../components/AdSenseAd/AdSenseAd";

export const Ads: FC = () => {
  return (
    <Container maxW="4xl" py={10}>
      <VStack spacing={8} align="stretch">
        <Heading>AdSense Test</Heading>
        <Text>
          This page is used to verify that responsive Google ads render correctly
          in production.
        </Text>
        <AdSenseAd slot="9118265430" />
      </VStack>
    </Container>
  );
};
