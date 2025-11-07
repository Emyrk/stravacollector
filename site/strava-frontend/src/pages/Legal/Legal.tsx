import { Box, Container, Heading, Text, VStack } from "@chakra-ui/react";
import { FC } from "react";

export const Legal: FC = () => {
  return (
    <Container maxW="4xl" py={10}>
      <VStack spacing={6} align="stretch">
        <Heading>Legal Disclaimer</Heading>
        
        <Box>
          <Heading size="md" mb={2}>Website Purpose</Heading>
          <Text>
            This website (dashugel.bike) is an unofficial, community-run results tracker 
            and informational resource. The site owner is not affiliated with, endorsed by, 
            or responsible for the organization or execution of Das Hügel or any other cycling 
            events referenced on this site.
          </Text>
        </Box>

        <Box>
          <Heading size="md" mb={2}>Limitation of Liability</Heading>
          <Text>
            The site owner assumes no liability for any injuries, damages, losses, or expenses 
            incurred by participants during any cycling event, ride, or activity. Cycling involves 
            inherent risks including but not limited to collisions, falls, equipment failure, and 
            traffic hazards.
          </Text>
        </Box>

        <Box>
          <Heading size="md" mb={2}>Assumption of Risk</Heading>
          <Text>
            By participating in any ride or event, you acknowledge that you do so at your own risk. 
            You are responsible for:
          </Text>
          <Text as="ul" pl={6} pt={2}>
            <li>Your own safety and the safety of others</li>
            <li>Ensuring your bicycle is in safe working condition</li>
            <li>Obeying all traffic laws and regulations</li>
            <li>Wearing appropriate safety equipment</li>
            <li>Riding within your skill level and physical capabilities</li>
          </Text>
        </Box>

        <Box>
          <Heading size="md" mb={2}>Data Accuracy</Heading>
          <Text>
            While we strive to provide accurate information, we make no warranties regarding the 
            accuracy, completeness, or reliability of ride data, routes, or results displayed on 
            this website. All data is sourced from Strava and other third-party services.
          </Text>
        </Box>

        <Box>
          <Heading size="md" mb={2}>Contact</Heading>
          <Text>
            For questions or concerns, please contact{" "}
            <Text as="a" href="mailto:help@dashugel.bike" color="brand.stravaOrange">
              help@dashugel.bike
            </Text>
          </Text>
        </Box>
      </VStack>
    </Container>
  );
};
