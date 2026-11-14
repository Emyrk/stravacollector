import { screen } from "@testing-library/react";
import { render } from "../../test-utils";
import { Ads } from "./Ads";

beforeEach(() => {
  window.adsbygoogle = [];
});

test("keeps the test page accessible without serving ads", () => {
  const { container } = render(<Ads />);

  expect(
    screen.getByRole("heading", { name: "AdSense Test" })
  ).toBeInTheDocument();
  expect(container.querySelector('[data-ad-slot="9118265430"]')).not.toBeInTheDocument();
  expect(window.adsbygoogle).toHaveLength(0);
});
