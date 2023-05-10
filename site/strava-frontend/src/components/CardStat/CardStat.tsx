import {
  Flex,
  Text,
  FlexProps,
  useStyleConfig,
  TextProps,
} from "@chakra-ui/react";
import { type } from "os";
import { FC, PropsWithChildren } from "react";

export type CardStatProps = FlexProps & {
  title: string;
  value: string;
  titleProps?: TextProps;
  valueProps?: TextProps;
};

export const CardStat: FC<CardStatProps> = ({
  title,
  titleProps,
  value,
  valueProps,
  ...props
}) => {
  props = {
    ...props,
    flexDirection: props.flexDirection || "column",
    alignItems: props.alignItems || "center",
    justifyContent: props.justifyContent || "center",
  };

  return (
    <Flex {...props}>
      <Text color="brand.cardStatTitle" fontSize={"0.85em"} {...titleProps}>
        {title}
      </Text>
      <Text color="brand.cardStatValue" fontSize={"1em"} {...valueProps}>
        {value}
      </Text>
    </Flex>
  );
};
