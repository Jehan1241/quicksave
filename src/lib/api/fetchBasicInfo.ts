import { showErrorToast } from "../toastService";
import { handleApiError } from "./apiErrors";

export const fetchData = async (
  sortType: string,
  sortOrder: string,
  setDataArray: (data: any[]) => void,
  setMetaData: (meta: any) => void,
  setSortOrder: (order: "ASC" | "DESC" | "default") => void,
  setSortType: (type: string) => void,
  setWishlistArray: (data: any[]) => void,
  setHiddenArray: (data: any[]) => void,
  setInstalledArray: (data: any[]) => void
) => {
  console.log("Sending Get Basic Info");
  try {
    const response = await fetch(
      `http://localhost:8080/getBasicInfo?type=${sortType}&order=${sortOrder}`
    );
    if (!response.ok) await handleApiError(response);
    const json = await response.json();

    const hiddenUIDs = json.HiddenUIDs || [];

    const filteredLibraryGames = Object.values(json.MetaData).filter(
      (item: any) => !hiddenUIDs.includes(item.UID) && item.isDLC === 0
    );
    const filteredWishlistGames = Object.values(json.MetaData).filter(
      (item: any) => !hiddenUIDs.includes(item.UID) && item.isDLC === 1
    );
    const hiddenGames = Object.values(json.MetaData).filter((item: any) =>
      hiddenUIDs.includes(item.UID)
    );
    const installedGames = Object.values(json.MetaData).filter(
      (item: any) => item.InstallPath !== "" && !hiddenUIDs.includes(item.UID)
    );

    console.log(installedGames);

    setDataArray(filteredLibraryGames);
    setMetaData(json.MetaData);
    setSortOrder(json.SortOrder);
    setSortType(json.SortType);
    setWishlistArray(filteredWishlistGames);
    setHiddenArray(hiddenGames);
    setInstalledArray(installedGames);

    console.log(json);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to load filters!", String(error));
  }
};
