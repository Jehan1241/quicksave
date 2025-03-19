import { Location } from "react-router-dom";

interface Game {
  UID: string;
}

export const pickRandomGame = (
  location: Location,
  dataArray: Game[],
  wishlistArray: Game[]
): string | null => {
  let targetArray: Game[] = [];

  if (location.pathname === "/") {
    targetArray = dataArray;
  } else if (location.pathname === "/wishlist") {
    targetArray = wishlistArray;
  } else if (location.pathname === "/gameview") {
    targetArray = dataArray.some((game) => game.UID === location.state?.data)
      ? dataArray
      : wishlistArray;
  }

  return targetArray.length > 0
    ? targetArray[Math.floor(Math.random() * targetArray.length)].UID
    : null;
};
