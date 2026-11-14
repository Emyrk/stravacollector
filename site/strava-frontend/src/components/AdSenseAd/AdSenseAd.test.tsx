import { render } from "../../test-utils";
import { AdSenseAd } from "./AdSenseAd";

beforeEach(() => {
  window.adsbygoogle = [];
});

test("does not render or initialize ads while AdSense is disabled", () => {
  const { container, rerender } = render(<AdSenseAd slot="9118265430" />);

  expect(container.querySelector("ins.adsbygoogle")).not.toBeInTheDocument();
  expect(window.adsbygoogle).toHaveLength(0);

  rerender(<AdSenseAd slot="9118265430" />);

  expect(container.querySelector("ins.adsbygoogle")).not.toBeInTheDocument();
  expect(window.adsbygoogle).toHaveLength(0);
});
