import {
  CURRENT_HUGEL_ROUTE_YEAR,
  CURRENT_HUGEL_YEAR,
  CURRENT_HUGEL_ROUTE_URL,
  CURRENT_HUGEL_LITE_ROUTE_YEAR,
} from "./hugel";

test("configures the 2026 full leaderboard with the published route", () => {
  expect(CURRENT_HUGEL_YEAR).toBe(2026);
  expect(CURRENT_HUGEL_ROUTE_YEAR).toBe(2026);
  expect(CURRENT_HUGEL_ROUTE_URL).toBe("https://www.strava.com/routes/3554628211898515602");
});

test("keeps Lite on the 2025 definition pending confirmation", () => {
  expect(CURRENT_HUGEL_LITE_ROUTE_YEAR).toBe(2025);
});
