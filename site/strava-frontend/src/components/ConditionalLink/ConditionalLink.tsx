import { Link, LinkProps, useStyleConfig } from "@chakra-ui/react";
import { FC, PropsWithChildren } from "react";

export type ConditionalLinkProps = LinkProps & {};

// ConditionalLink only applies a Link wrapper if href is not empty.
export const ConditionalLink: FC<PropsWithChildren<ConditionalLinkProps>> = ({
  children,
  ...props
}) => {
  if (props.href && props.href !== "") {
    return <Link {...props}>{children}</Link>;
  }
  return <>{children}</>;
};
