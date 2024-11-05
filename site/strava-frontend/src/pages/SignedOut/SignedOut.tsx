import {
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  Box,
  Button,
  Heading,
  useTheme,
} from "@chakra-ui/react";
import { FC } from "react";
import { Link as RouteLink } from "react-router-dom";

export const SignedOut: FC = () => {
  const theme = useTheme();

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        minHeight: "100vh",
        backgroundColor: "#1a202c", //theme.colors.brand.primary,
      }}
    >
      <Alert
        status="success"
        variant="subtle"
        flexDirection="column"
        alignItems="center"
        justifyContent="center"
        textAlign="center"
        height="200px"
        backgroundColor="transparent"
      >
        <AlertIcon boxSize="40px" mr={0} />
        <AlertTitle mt={4} mb={1} fontSize="lg">
          Signed Out
        </AlertTitle>
        <AlertDescription maxWidth="sm">
          You have been sucessfully signed out.
        </AlertDescription>
        <RouteLink to="/">
          <Button
            size={"lg"}
            textColor={theme.colors.brand.primary}
            marginTop={7}
            backgroundColor={theme.colors.brand.secondary}
          >
            Back Home
          </Button>
        </RouteLink>
      </Alert>
    </Box>
  );
};
