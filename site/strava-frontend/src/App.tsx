import * as React from "react"
import {
  ChakraProvider,
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  extendTheme,
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  useColorModeValue,
} from "@chakra-ui/react"
import { Logo } from "./Logo"
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link as RouteLink,
  Outlet
} from "react-router-dom";
import { HugelBoard } from "./pages/HugelBoard/HugelBoard";
import { Landing } from "./pages/Landing/Landing";
import Navbar from "./components/Navbar/Navbar";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthenticatedProvider } from "./contexts/Authenticated";
import { FC } from "react";
import { NotFound } from "./pages/404/404";
import { SignedOut } from "./pages/SignedOut/SignedOut";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      cacheTime: 0,
      refetchOnWindowFocus: false,
      networkMode: "offlineFirst",
    },
  },
})


const theme = extendTheme({
  components: {
    Tabs: {
      baseStyle: {
        tab: {
          _selected: {
            color: "#fc4c02",
          }
        }
      },
    }
  },
  colors: {
    brand: {
      primary: "#ebebeb",
      stravaOrange: "#fc4c02",
    },
  },
})

export const App = () => {
  return <QueryClientProvider client={queryClient}>
    <AuthenticatedProvider>
      <ChakraProvider theme={theme}>
        <Router>
          <Routes>
            <Route element={<IncludeNavbar />}>
              {/* Navbar and statics */}
              <Route path="/" element={<Landing />} />
              <Route path="/hugelboard" element={<HugelBoard />} />
              <Route path="/signed-out" element={<SignedOut />} />
            </Route>
            <Route path='*' element={<NotFound />} />
          </Routes>
        </Router>
      </ChakraProvider>
    </AuthenticatedProvider>
  </QueryClientProvider>
}

export const IncludeNavbar: FC = () => {
  const bgSrc = useColorModeValue("dark", "light")
  return <>
    <Box h={'100svh'} background={`url(/hugel_route_lines_${bgSrc}.svg)`} backgroundPosition={'center'} backgroundRepeat={'no-repeat'} backgroundSize='cover'  >
      <Navbar />
      <Outlet />
    </Box>
  </>
}
