import { fireEvent, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { render } from "../../test-utils";
import { MobileNav2 } from "./Navbar";
import { CURRENT_HUGEL_YEAR } from "../../hugel";

// Navigation does not need authentication or real API requests.
jest.mock("../../api/rest", () => ({}));

test("menu links to current-year full and Lite results", async () => {
  render(
    <MemoryRouter>
      <MobileNav2 display={{ base: "block", md: "block" }} />
    </MemoryRouter>
  );
  fireEvent.click(screen.getByRole("button", { name: "Options" }));
  const full = await screen.findByText(`${CURRENT_HUGEL_YEAR} Das Hügel Results`);
  const lite = screen.getByText(`${CURRENT_HUGEL_YEAR} Das Hügel Lite Results`);
  expect(full.closest("a")).toHaveAttribute("href", `/hugelboard/${CURRENT_HUGEL_YEAR}`);
  expect(lite.closest("a")).toHaveAttribute("href", `/hugelboard/${CURRENT_HUGEL_YEAR}?lite=true`);
  expect(screen.queryByText("2025 Das Hügel Results")).not.toBeInTheDocument();
  expect(screen.getByText("Other Results").closest("a")).toHaveAttribute("href", "/results");
});
