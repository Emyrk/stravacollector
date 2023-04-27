import { Link, Text, Box, Grid, VStack, Code } from "@chakra-ui/react"
import { FC } from "react"
import { ColorModeSwitcher } from "../../ColorModeSwitcher"
import { Logo } from "../../Logo"
import {
  Link as RouteLink
} from "react-router-dom";

export const Landing: FC = () => {
  return <Box textAlign="center" fontSize="xl">
    <Grid minH="100vh" p={3}>
      <ColorModeSwitcher justifySelf="flex-end" />
      <VStack spacing={8}>
        <Logo h="40vmin" pointerEvents="none" />
        <Text>
          Edit <Code fontSize="xl">src/App.tsx</Code> and save to reload.
        </Text>
        <RouteLink to="/hugelboard">Hugelboard</RouteLink>
        <Link
          color="teal.500"
          href="https://chakra-ui.com"
          fontSize="2xl"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn Chakra
        </Link>
      </VStack>
    </Grid>
  </Box>
}