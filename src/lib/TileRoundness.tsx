const TILE_ROUNDNESS_KEY = "tile-roundness";

export function setTileRoundness(roundness: string) {
  localStorage.setItem(TILE_ROUNDNESS_KEY, roundness);
}

export function getTileRoundness() {
  const roundness = localStorage.getItem("tile-roundness");
  switch (roundness) {
    case "none":
      return "None";
    case "sm":
      return "Small";
    case "md":
      return "Medium";
    case "lg":
      return "Large";
    case "xl":
      return "Extra Large";
    case "2xl":
      return "2xl";
    case "3xl":
      return "3xl";
    default:
      return "Large";
  }
}
