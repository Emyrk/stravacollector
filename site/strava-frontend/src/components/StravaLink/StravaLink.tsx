import { ImageProps, Image, Link, Tooltip } from "@chakra-ui/react";
import { FC } from "react";
import { ConditionalTooltip } from "../ConditionalTooltip/ConditionalTooltip";

export type StravaLinkProps = ImageProps & {
  href: string;
  target?: string;
  tooltip?: string;
};

export const StravaLink: FC<StravaLinkProps> = ({
  href,
  target,
  tooltip,
  ...props
}) => {
  return (
    <ConditionalTooltip label={tooltip || ""} aria-label="Strava logo tooltip">
      <Link href={href} target={target} variant={"stravaLink"}>
        <Image
          src={"/logos/stravalogo.png"}
          height={"34px"}
          width={"34px"}
          maxWidth={"34px"}
          {...props}
        />
      </Link>
    </ConditionalTooltip>
  );
};
