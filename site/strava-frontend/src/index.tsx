import { ChakraProvider, ColorModeScript, StyleConfig, ThemeConfig, extendTheme, createLocalStorageManager } from "@chakra-ui/react"
import * as React from "react"
import * as ReactDOM from "react-dom/client"
import { App } from "./App"
import reportWebVitals from "./reportWebVitals"
import * as serviceWorker from "./serviceWorker"


const container = document.getElementById("root")
if (!container) throw new Error('Failed to find the root element');
const root = ReactDOM.createRoot(container)

const customComponents: Record<string, StyleConfig> = {
  Text: {
    variants: {
      // used as <Text variant="minor">
      minor: ({ colorMode }) => ({
        color: colorMode === "dark" ? "whiteAlpha.500" : "blackAlpha.500",
      })
    },
  },
  Tabs: {
    baseStyle: {
      tab: {
        _selected: {
          color: "#fc4c02",
        }
      }
    },
  },
  Tag: {
    baseStyle: {
      container: {
        color: "white",
      }
    },
    variants: {},
  }
}


export const config: ThemeConfig = {
  initialColorMode: "dark",
  useSystemColorMode: false,
};

const theme = extendTheme({
  components: { ...customComponents },
  colors: {
    brand: {
      primary: "#ebebeb",
      stravaOrange: "#fc4c02",
    },
  },
}, { config })

export default theme

const manager = createLocalStorageManager("strava-chakra-ui-color-mode")

root.render(
  < React.StrictMode >
  <ChakraProvider theme={theme} colorModeManager={manager} >
    <ColorModeScript initialColorMode={theme.config.initialColorMode} />
    <App />
    </ChakraProvider>
  </React.StrictMode >
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://cra.link/PWA
serviceWorker.unregister()

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()

