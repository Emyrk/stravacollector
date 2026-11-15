export const CURRENT_HUGEL_YEAR: number = 2026;
export const CURRENT_HUGEL_ROUTE_YEAR: number = 2026;
export const CURRENT_HUGEL_ROUTE_URL = "https://www.strava.com/routes/3554628211898515602";

// The published 2026 Lite route retains the same scoring segments as 2025.
export const CURRENT_HUGEL_LITE_ROUTE_YEAR: number = 2026;
export const CURRENT_HUGEL_LITE_ROUTE_URL = "https://www.strava.com/routes/3555623609985396008";

// Alternative for the 2026 full-route Cat Mountain requirement only.
export const CAT_MOUNTAIN_ORIGINAL = 891091;
export const CAT_MOUNTAIN_ALTERNATE = 784585;
export const catMountainAlternativeName = (segmentID: string | number): string | undefined =>
  String(segmentID) === String(CAT_MOUNTAIN_ALTERNATE) ? "Cat Mountain Dr" : undefined;
