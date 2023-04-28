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

export const App = () => (
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
)
