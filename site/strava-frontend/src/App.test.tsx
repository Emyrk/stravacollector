import {
  CURRENT_HUGEL_ROUTE_YEAR,
  CURRENT_HUGEL_YEAR,
} from "./hugel";

test("configures the 2026 leaderboard with the temporary 2025 route", () => {
  expect(CURRENT_HUGEL_YEAR).toBe(2026);
  expect(CURRENT_HUGEL_ROUTE_YEAR).toBe(2025);
});
