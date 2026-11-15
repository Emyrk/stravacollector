import { screen } from "@testing-library/react";
import { render } from "../../test-utils";
import { catMountainAlternativeName } from "../../hugel";
import { CatMountainAlternatives } from "./CatMountainAlternatives";

test("explains both accepted Cat Mountain segments and fastest-only scoring", () => {
  render(<CatMountainAlternatives />);
  expect(screen.getByRole("link", { name: "Cat Mountain (Hugel Route)" })).toHaveAttribute("href", "https://www.strava.com/segments/891091");
  expect(screen.getByRole("link", { name: "Cat Mountain Dr" })).toHaveAttribute("href", "https://www.strava.com/segments/784585");
  expect(screen.getByText(/Only the fastest effort across these two segments counts, once/)).toBeInTheDocument();
});

test("labels an actual winning alternate effort without relabeling other segments", () => {
  expect(catMountainAlternativeName("784585")).toBe("Cat Mountain Dr");
  expect(catMountainAlternativeName(784585)).toBe("Cat Mountain Dr");
  expect(catMountainAlternativeName(891091)).toBeUndefined();
});
