import { Location } from "react-router-dom";

interface Game {
  UID: string;
}

export const pickRandomGame = (
  lastLibraryPath: string,
  dataArray: Game[],
  wishlistArray: Game[],
  installedArray: Game[]
): string | null => {
  let targetArray: Game[] = [];

  if (lastLibraryPath === "/library") {
    targetArray = dataArray;
  } else if (lastLibraryPath === "/wishlist") {
    targetArray = wishlistArray;
  } else if (lastLibraryPath === "/installed") {
    targetArray = installedArray;
  }

  return targetArray.length > 0
    ? targetArray[Math.floor(Math.random() * targetArray.length)].UID
    : null;
};
