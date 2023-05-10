import { Tooltip, TooltipProps, useStyleConfig } from "@chakra-ui/react";
import { FC, PropsWithChildren } from "react";

export type ConditionalTooltipProps = TooltipProps & {};

// ConditionalLink only applies a Link wrapper if href is not empty.
export const ConditionalTooltip: FC<
  PropsWithChildren<ConditionalTooltipProps>
> = ({ children, ...props }) => {
  if (props.label && props.label !== "") {
    return <Tooltip {...props}>{children}</Tooltip>;
  }
  return <>{children}</>;
};
