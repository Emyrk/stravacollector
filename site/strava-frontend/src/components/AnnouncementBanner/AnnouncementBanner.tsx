import { Box, Container, Text, CloseButton, useDisclosure } from "@chakra-ui/react";
import { useColorModeValue } from "@chakra-ui/react";
import { FC } from "react";

export const AnnouncementBanner: FC = () => {
  const { isOpen, onClose } = useDisclosure({ defaultIsOpen: true });
  const bgColor = useColorModeValue("brand.stravaOrange", "orange.600");
  const textColor = "white";

  if (!isOpen) {
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