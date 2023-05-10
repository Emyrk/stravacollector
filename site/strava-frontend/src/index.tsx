import {
  ChakraProvider,
  ColorModeScript,
  StyleConfig,
  ThemeConfig,
  extendTheme,
  createLocalStorageManager,
  theme as defaultTheme,
} from "@chakra-ui/react";
import * as React from "react";
import * as ReactDOM from "react-dom/client";
import { App } from "./App";
import reportWebVitals from "./reportWebVitals";
import * as serviceWorker from "./serviceWorker";

const container = document.getElementById("root");
if (!container) throw new Error("Failed to find the root element");
const root = ReactDOM.createRoot(container);

const customComponents: Record<string, StyleConfig> = {
  Box: {
    variants: {
      responsiveCard: ({ colorMode }) => ({
        // bgColor={"#272c35"}
        // bgColor={"#3b3f48"}
        // https://www.sessions.edu/color-calculator/
        bgGradient:
          "linear(135deg, brand.primaryCard 0%, brand.secondaryCard 200%)",
        borderRadius: "10px",
        boxShadow: "rgb(20, 20, 20) 0px 3px 6px",
        transition: "all .25s ease",
        _hover: {
          bgGradient:
            "linear(135deg, brand.primaryCard 0%, brand.secondaryCard 35%)",
          boxShadow: "rgb(20, 20, 20) 0px 5px 10px",
          marginTop: "-3px",
          marginBottom: "-3px",
        },
      }),
    },
  },
  // Flex: {
  //   variants: {
  //     cardStat: ({ colorMode }) => ({
  //       // defaultTheme.components.Flex
  //       flexDirection: "column",
  //       alignItems: "center",
  //       justifyContent: "center",
  //     }),
  //   },
  // },
  Text: {
    variants: {
      // used as <Text variant="minor">
      minor: ({ colorMode }) => ({
        color: colorMode === "dark" ? "whiteAlpha.500" : "blackAlpha.500",
      }),
    },
  },
  Tabs: {
    baseStyle: {
      tab: {
        _selected: {
          color: "#fc4c02",
        },
      },
    },
  },
  Tag: {
    baseStyle: {
      container: {
        color: "white",
      },
    },
    variants: {},
  },
};

export const config: ThemeConfig = {
  initialColorMode: "dark",
  useSystemColorMode: false,
};

const theme = extendTheme(
  {
    components: { ...customComponents },
    colors: {
      brand: {
        primary: "#ebebeb",
        stravaOrange: "#fc4c02",
        primaryCard: "#3b3f48",
        secondaryCard: "#48403b",
        cardStatTitle: "#a7afbe",
        cardStatValue: "white",
      },
    },
  },
  { config }
);

export default theme;

const manager = createLocalStorageManager("strava-chakra-ui-color-mode");

root.render(
  <React.StrictMode>
    <ChakraProvider theme={theme} colorModeManager={manager}>
      <ColorModeScript initialColorMode={theme.config.initialColorMode} />
      <App />
    </ChakraProvider>
  </React.StrictMode>
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://cra.link/PWA
serviceWorker.unregister();

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
