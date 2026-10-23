import { Box } from "@chakra-ui/react";
import { FC, useEffect, useRef } from "react";

declare global {
  interface Window {
    adsbygoogle: Record<string, never>[];
  }
}

interface AdSenseAdProps {
  slot: string;
}

export const AdSenseAd: FC<AdSenseAdProps> = ({ slot }) => {
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) {
      return;
    }

    initialized.current = true;
    window.adsbygoogle = window.adsbygoogle || [];
    window.adsbygoogle.push({});
  }, []);

  return (
    <Box minH="100px" width="100%" overflow="hidden">
      <ins
        className="adsbygoogle"
        style={{ display: "block" }}
        data-ad-client="ca-pub-8208259743822818"
        data-ad-slot={slot}
        data-ad-format="auto"
        data-full-width-responsive="true"
      />
    </Box>
  );
};
