import {
  CURRENT_HUGEL_ROUTE_YEAR,
  CURRENT_HUGEL_YEAR,
  CURRENT_HUGEL_ROUTE_URL,
  CURRENT_HUGEL_LITE_ROUTE_YEAR,
  CURRENT_HUGEL_LITE_ROUTE_URL,
} from "./hugel";

test("configures the 2026 full leaderboard with the published route", () => {
  expect(CURRENT_HUGEL_YEAR).toBe(2026);
  expect(CURRENT_HUGEL_ROUTE_YEAR).toBe(2026);
  expect(CURRENT_HUGEL_ROUTE_URL).toBe("https://www.strava.com/routes/3554628211898515602");
});

test("configures the confirmed 2026 Lite route without a temporary-year warning", () => {
  expect(CURRENT_HUGEL_LITE_ROUTE_YEAR).toBe(CURRENT_HUGEL_YEAR);
  expect(CURRENT_HUGEL_LITE_ROUTE_URL).toBe("https://www.strava.com/routes/3555623609985396008");
});
