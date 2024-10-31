import { Box, Heading, Button, useTheme, Text } from "@chakra-ui/react";
import { Link as RouteLink } from "react-router-dom";
import { FC } from "react";

export const ErrorBox: FC<{ error: string }> = ({ error }) => {
  const theme = useTheme();

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        pt: "50px",
        // minHeight: "100vh",
      }}
    >
      <Heading pb="50px" color={theme.colors.brand.stravaOrange}>
        Error
      </Heading>

      <Text>{error}</Text>
    </Box>
  );
};
