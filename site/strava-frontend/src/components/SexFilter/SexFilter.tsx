import { FC } from "react";
import { Flex, ButtonGroup, Button } from "@chakra-ui/react";

export type SexFilter = "all" | "M" | "F";

interface SexFilterProps {
  value: SexFilter;
  onChange: (sex: SexFilter) => void;
}

export const SexFilterButtons: FC<SexFilterProps> = ({ value, onChange }) => {
  const getButtonStyles = (isActive: boolean) => ({
    bg: isActive ? "brand.stravaOrange" : "transparent",
    color: isActive ? "white" : "gray.400",
    borderColor: isActive ? "brand.stravaOrange" : "gray.600",
    _hover: {
      bg: isActive ? "brand.stravaOrange" : "gray.700",
      borderColor: "brand.stravaOrange",
      "& + button": {
        borderLeftColor: "brand.stravaOrange",
      },
    },
  });

  return (
    <Flex justifyContent="center" pt={4}>
      <ButtonGroup isAttached variant="outline" size="md">
        <Button
          {...getButtonStyles(value === "all")}
          onClick={() => onChange("all")}
        >
          All
        </Button>
        <Button
          {...getButtonStyles(value === "M")}
          onClick={() => onChange("M")}
        >
          Men
        </Button>
        <Button
          {...getButtonStyles(value === "F")}
          onClick={() => onChange("F")}
        >
          Women
        </Button>
      </ButtonGroup>
    </Flex>
  );
};