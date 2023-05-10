import { Box, ChakraProps, useStyleConfig } from "@chakra-ui/react";
import { FC, PropsWithChildren } from "react";

export type ResponsiveCardProps = ChakraProps & {};

export const ResponsiveCard: FC<PropsWithChildren<ResponsiveCardProps>> = ({
  children,
  ...props
}) => {
  const styles = useStyleConfig("Box", { variant: "responsiveCard" });

  return (
    <Box __css={styles} {...props}>
      {children}
    </Box>
  );
};
