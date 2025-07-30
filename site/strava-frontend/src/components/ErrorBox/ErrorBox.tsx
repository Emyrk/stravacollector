import { Box, Heading, Button, useTheme, Text, Collapse, useDisclosure, Code } from "@chakra-ui/react";
import { Link as RouteLink } from "react-router-dom";
import { FC } from "react";

export const ErrorBox: FC<{ error: string, detail?: unknown}> = ({ error, detail }) => {
  const theme = useTheme();
  const { isOpen, onToggle } = useDisclosure()

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        pt: "50px",
        // minHeight: "100vh",
      }}
    >
      <Heading pb="50px" color={theme.colors.brand.stravaOrange}>
        Error
      </Heading>

      <Text>{error}</Text>
    
      {
        detail !== undefined && detail !== null && (<>
            <br />
            <Code colorScheme='red'>{detail.toString()}</Code>
          </>
        )
      }
      
    </Box>
  );
};
