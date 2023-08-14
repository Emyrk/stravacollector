import { Box, defineStyleConfig, useStyleConfig } from "@chakra-ui/react";
import { FC } from "react";
import "./Loading.css";

const LoadingCSS = defineStyleConfig({
  baseStyle: {
    background: "#555",
    width: "20px",
    height: "20px",
    ".loop": {
      background: "#555",
      width: "20px",
      height: "20px",
    },
  },
});

export const Loading: FC = () => {
  return (
    <>
      <Box>
        <Box id="loop" className="center"></Box>
        <div id="bike-wrapper" className="center">
          <div id="bike" className="centerBike"></div>
          <div id="bike" className="centerBike"></div>
        </div>
      </Box>
    </>
  );
};
