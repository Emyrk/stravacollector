import * as React from "react"
import {
  ChakraProvider,
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  theme,
  extendTheme,
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
} from "@chakra-ui/react"
import { ColorModeSwitcher } from "./ColorModeSwitcher"
import { Logo } from "./Logo"
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link as RouteLink,
  Outlet
} from "react-router-dom";
import { HugelBoard } from "./pages/HugelBoard";
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



export const App = () => {
  const theme = extendTheme({
    colors: {
      brand: {
        primary: "#ebebeb",
        secondary: "#fc4c02",
      },
    },
  })

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
  return <>
    <Navbar />
    <Outlet />
  </>
}
