import { render } from "../../test-utils";
import { AdSenseAd } from "./AdSenseAd";

beforeEach(() => {
  window.adsbygoogle = [];
});

test("renders and initializes a responsive ad once", () => {
  const { container, rerender } = render(<AdSenseAd slot="9118265430" />);

  const ad = container.querySelector("ins.adsbygoogle");
  expect(ad).toHaveAttribute("data-ad-client", "ca-pub-8208259743822818");
  expect(ad).toHaveAttribute("data-ad-slot", "9118265430");
  expect(ad).toHaveAttribute("data-ad-format", "auto");
  expect(ad).toHaveAttribute("data-full-width-responsive", "true");
  expect(window.adsbygoogle).toHaveLength(1);

  rerender(<AdSenseAd slot="9118265430" />);

  expect(window.adsbygoogle).toHaveLength(1);
});
