import { Link, Text } from "@chakra-ui/react";
import { CAT_MOUNTAIN_ORIGINAL, CAT_MOUNTAIN_ALTERNATE } from "../../hugel";

export const CatMountainAlternatives = () => (
  <Text pt="1em">
    For the 2026 full route, Cat Mountain accepts either{" "}
    <Link href={`https://www.strava.com/segments/${CAT_MOUNTAIN_ORIGINAL}`}>
      Cat Mountain (Hugel Route)
    </Link>{" "}
    or{" "}
    <Link href={`https://www.strava.com/segments/${CAT_MOUNTAIN_ALTERNATE}`}>
      Cat Mountain Dr
    </Link>.
    Only the fastest effort across these two segments counts, once.
  </Text>
);
