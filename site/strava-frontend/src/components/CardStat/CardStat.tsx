import { Flex, Text, FlexProps, useStyleConfig } from "@chakra-ui/react";
import { type } from "os";
import { FC, PropsWithChildren } from "react";

export type CardStatProps = FlexProps & {
  title: string;
  value: string;
  // valueUnit?: string;
};

export const CardStat: FC<CardStatProps> = ({ title, value, ...props }) => {
  props = {
    ...props,
    flexDirection: props.flexDirection || "column",
    alignItems: props.alignItems || "center",
    justifyContent: props.justifyContent || "center",
  };

  // let valueString: string;
  // switch (typeof value) {
  //   case "number":
  //     valueString = value.toLocaleString()
  //     break;
  //   case "string":
  //     valueString = value;
  //     break;
  // }

  return (
    <Flex {...props}>
      <Text color="brand.cardStatTitle" fontSize={"13px"}>
        {title}
      </Text>
      <Text color="brand.cardStatValue" fontSize={"16px"}>
        {value}
      </Text>
    </Flex>
  );
};
