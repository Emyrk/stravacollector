import * as React from "react";
import {
  ChakraProvider,
  Box,
  Text,
  Link,
  VStack,
  Code,
  Grid,
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  useColorModeValue,
  StyleConfig,
  StyleFunctionProps,
  ThemeConfig,
  Container,
} from "@chakra-ui/react";
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link as RouteLink,
  Outlet,
} from "react-router-dom";
import { HugelBoard } from "./pages/HugelBoard/HugelBoard";
import { Landing } from "./pages/Landing/Landing";
import Navbar from "./components/Navbar/Navbar";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthenticatedProvider } from "./contexts/Authenticated";
import { FC } from "react";
import { NotFound } from "./pages/404/404";
import { SignedOut } from "./pages/SignedOut/SignedOut";
import { SuperHugelBoard } from "./pages/HugelBoard/SuperHugelBoard";
import { Activity } from "./components/Activity/Activity";
import { ChallengeRoute } from "./pages/ChallengeRoute/ChallengeRoute";
import { Footer } from "./components/Footer/Footer";
import { AthletePage } from "./pages/AthletePage/AthletePage";
import { Results } from "./pages/Results/Results";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      cacheTime: 0,
      refetchOnWindowFocus: false,
      networkMode: "offlineFirst",
    },
  },
});

export const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthenticatedProvider>
        <Router>
          <Routes>
            <Route element={<Background />}>
              <Route element={<IncludeNavbar />}>
                {/* Navbar and statics */}
                <Route path="/" element={<Landing />} />
                <Route path="/athlete/:athlete_id" element={<AthletePage />} />
                <Route path="/hugelboard" element={<HugelBoard />} />
                <Route path="/superhugelboard" element={<SuperHugelBoard />} />
                <Route path="/results" element={<Results />} />
                <Route path="/activity/:activity_id" element={<Activity />} />
                <Route path="/route/:name" element={<ChallengeRoute />} />
                <Route path="/signed-out" element={<SignedOut />} />
              </Route>
              <Route path="*" element={<NotFound />} />
            </Route>
          </Routes>
        </Router>
      </AuthenticatedProvider>
    </QueryClientProvider>
  );
};

export const Background: FC = () => {
  const bgSrc = useColorModeValue("dark", "light");
  return (
    <>
      <Box
        h={"100svh"}
        background={`url(/hugel_route_lines_${bgSrc}.svg)`}
        backgroundPosition={"center"}
        backgroundRepeat={"no-repeat"}
        backgroundSize="cover"
      >
        <Outlet />
      </Box>
    </>
  );
};

export const IncludeNavbar: FC = () => {
  const bgSrc = useColorModeValue("dark", "light");
  return (
    <>
      <Box>
        <Box minH={"90vh"}>
          <Navbar />
          <Outlet />
        </Box>

        <Footer />
      </Box>
    </>
  );
};
