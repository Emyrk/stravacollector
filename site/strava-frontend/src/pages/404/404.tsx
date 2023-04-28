import { Box, Button, Heading, useTheme } from "@chakra-ui/react";
import { FC } from "react"
import {
  Link as RouteLink
} from "react-router-dom";

export const NotFound: FC = () => {
  const theme = useTheme()
  console.log(theme.colors.brand)

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
        minHeight: '100vh',
        backgroundColor: theme.colors.brand.primary,
      }}
    >
      {/* <Typography variant="h1" style={{ color: 'white' }}>
        404
      </Typography>
      <Typography variant="h6" style={{ color: 'white' }}>
        The page you’re looking for doesn’t exist.
      </Typography>
      <Button variant="contained">Back Home</Button> */}
      <Heading color={theme.colors.brand.secondary}>404 Not Found</Heading>

      <RouteLink to="/">
        <Button size={"lg"} textColor={theme.colors.brand.primary} marginTop={7} backgroundColor={theme.colors.brand.secondary}>Back Home</Button>
      </RouteLink>

    </Box >
  );
}