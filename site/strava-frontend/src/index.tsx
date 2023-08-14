import {
  ChakraProvider,
  ColorModeScript,
  StyleConfig,
  ThemeConfig,
  extendTheme,
  createLocalStorageManager,
  theme as defaultTheme,
  createMultiStyleConfigHelpers,
  AlertProps,
} from "@chakra-ui/react";
import { alertAnatomy } from "@chakra-ui/anatomy";
import * as React from "react";
import * as ReactDOM from "react-dom/client";
import { App } from "./App";
import reportWebVitals from "./reportWebVitals";
import * as serviceWorker from "./serviceWorker";
import { mode } from "@chakra-ui/theme-tools";

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(alertAnatomy.keys);

const container = document.getElementById("root");
if (!container) throw new Error("Failed to find the root element");
const root = ReactDOM.createRoot(container);

// console.log("defaultTheme", defaultTheme);
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
          transform: "translateY(-3px)",
        },
      }),
    },
  },
  Link: {
    variants: {
      stravaLink: ({ colorMode }) => ({
        ...defaultTheme.components.Link.baseStyle,
        // Apply transitions to the underlying image
        img: {
          transition: "all .1s ease",
        },
        _hover: {
          img: {
            transform: "scale(1.1)",
          },
        },
      }),
    },
  },
  // Alert: {
  //   variants: {
  //     moreOpaque: (props) => {
  //       const { colorScheme: c } = props;

  //       return {
  //         bg: `${c}.500`,
  //         color: `${c}.50`,
  //       };
  //     },
  // },
  Alert: {
    baseStyle: definePartsStyle((props: AlertProps) => {
      const { status } = props;
      const num = "600";

      const successBase = status === "success" && {
        container: {
          background: `green.${num}`,
        },
      };

      const warningBase = status === "warning" && {
        container: {
          background: `yellow.${num}`,
        },
      };

      const errorBase = status === "error" && {
        container: {
          background: `red.${num}`,
        },
      };

      const infoBase = status === "info" && {
        container: {
          background: `blue.${num}`,
        },
      };

      return {
        ...successBase,
        ...warningBase,
        ...errorBase,
        ...infoBase,
      };
    }),
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
    // styles: {
    //   global: () => ({
    //     body: {
    //     },
    //   }),
    // },
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
