import { Container, Heading, Text, VStack } from "@chakra-ui/react";
import { FC } from "react";
import { AdSenseAd } from "../../components/AdSenseAd/AdSenseAd";

export const Ads: FC = () => {
  return (
    <Container maxW="4xl" py={10}>
      <VStack spacing={8} align="stretch">
        <Heading>AdSense Test</Heading>
        <Text>
          Google AdSense is temporarily disabled.
        </Text>
        <AdSenseAd slot="9118265430" />
      </VStack>
    </Container>
  );
};
