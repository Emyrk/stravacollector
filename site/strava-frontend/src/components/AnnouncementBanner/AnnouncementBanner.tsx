import { Box, Container, Text, CloseButton, useDisclosure } from "@chakra-ui/react";
import { useColorModeValue } from "@chakra-ui/react";
import { FC } from "react";

export const AnnouncementBanner: FC = () => {
  const { isOpen, onClose } = useDisclosure({ defaultIsOpen: true });
  const bgColor = useColorModeValue("brand.stravaOrange", "orange.600");
  const textColor = "white";

   // Hide banner after 9am Central on November 8, 2024
  const cutoffDate = new Date("2025-11-08T09:00:00-06:00"); // -06:00 is Central Time
  const now = new Date();
  
  if (now >= cutoffDate || !isOpen) {
    return null;
  }

  return (
    <Box
      bg={bgColor}
      color={textColor}
      py={3}
      px={4}
      position="relative"
      w="100%"
    >
      <Container maxW="7xl" display="flex" alignItems="center" justifyContent="center">
        <Text fontSize={{ base: "sm", md: "md" }} fontWeight="bold" textAlign="center">
          🚴 The ride will happen on November 8th, 2025 with rollout at 7:15am! 🚴
        </Text>
        <CloseButton
          position="absolute"
          right={4}
          onClick={onClose}
          aria-label="Close announcement"
        />
      </Container>
    </Box>
  );
};