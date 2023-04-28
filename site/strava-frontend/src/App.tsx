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
} from "@chakra-ui/react"
import { ColorModeSwitcher } from "./ColorModeSwitcher"
import { Logo } from "./Logo"
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link as RouteLink
} from "react-router-dom";
import { HugelBoard } from "./pages/HugelBoard";
import { Landing } from "./pages/Landing/Landing";
import Navbar from "./components/Navbar/Navbar";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthenticatedProvider } from "./contexts/Authenticated";

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


export const App = () => (
  <QueryClientProvider client={queryClient}>
    <AuthenticatedProvider>
      <ChakraProvider theme={theme}>
        <Navbar />
        <Router>
          <Routes>
            {/* Navbar and statics */}
            <Route path="/" element={<Landing />} />
            <Route path="/hugelboard" element={<HugelBoard />} />
          </Routes>
        </Router>
      </ChakraProvider>
    </AuthenticatedProvider>
  </QueryClientProvider>
)
