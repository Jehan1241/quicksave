import { useEffect, useState } from "react";
import { setTheme } from "./ToggleTheme";
import { useLocation, useNavigate } from "react-router-dom";
import { Route, Routes } from "react-router-dom";
import CustomTitleBar from "./components/CustomTitleBar/CustomTitleBar";
import AddGameManuallyDialog from "./components/Dialogs/AddGameManually";
import { Toaster } from "@/components/ui/toaster";
import Integrations from "./components/Dialogs/Integrations";
import GameView from "./components/GameView/GameView";
import { LibraryView } from "./components/LibraryView/LibraryView";
import WishlistDialog from "./components/Dialogs/WishListDialog";
import { useSortContext } from "./hooks/useSortContex";
import BackButtonListener from "./hooks/BackButtonListener";
import { attachSSEListener } from "./lib/attachSSEListener";

function App() {
  const {
    sortType,
    sortOrder,
    setSortType,
    setSortOrder,
    setMetaData,
    setTileSize,
    tileSize,
    sortStateUpdate,
    setSortStateUpdate,
    randomGameClicked,
    setRandomGameClicked,
    isAddGameDialogOpen,
    isIntegrationsDialogOpen,
    isWishlistAddDialogOpen,
    setIntegrationLoadCount,
  } = useSortContext();
  const location = useLocation();
  const [dataArray, setDataArray] = useState<any[]>([]);
  const [wishlistArray, setWishlistArray] = useState<any[]>([]);
  const [hiddenArray, setHiddenArray] = useState<any[]>([]);
  const [installedArray, setInstalledArray] = useState<any[]>([]);
  const [integrationsPreviouslyOpened, setIntegrationsPreviouslyOpened] =
    useState<boolean>(false);
  const [
    wishListAddDialogPreviouslyOpened,
    setWishlistAddDialogPreviouslyOpened,
  ] = useState<boolean>(false);
  const [addGameDialogHasBeenOpened, setAddGameDialogHasBeenOpened] =
    useState<boolean>(false);
  const navigate = useNavigate();

  const fetchData = async () => {
    console.log("Sending Get Basic Info");
    try {
      const response = await fetch(
        `http://localhost:8080/getBasicInfo?type=${sortType}&order=${sortOrder}&size=${tileSize}`
      );
      const json = await response.json();

      // Extract the Hidden UIDs
      const hiddenUIDs = json.HiddenUIDs || [];

      // Filter out games with Hidden UIDs and separate them into hiddenArray
      const filteredLibraryGames = Object.values(json.MetaData).filter(
        (item: any) => {
          return item.isDLC === 0 && !hiddenUIDs.includes(item.UID); // Exclude games with hidden UIDs
        }
      );

      const filteredWishlistGames = Object.values(json.MetaData).filter(
        (item: any) => {
          return item.isDLC === 1 && !hiddenUIDs.includes(item.UID); // Exclude games with hidden UIDs
        }
      );

      const hiddenGames = Object.values(json.MetaData).filter((item: any) =>
        hiddenUIDs.includes(item.UID)
      ); // Include games with hidden UIDs

      const installedGames = Object.values(json.MetaData).filter(
        (item: any) => item.InstallPath != ""
      );

      console.log(installedGames);

      // Update state with filtered data
      setDataArray(filteredLibraryGames);
      setMetaData(json.MetaData);
      setSortOrder(json.SortOrder);
      setSortType(json.SortType);
      setWishlistArray(filteredWishlistGames);
      setHiddenArray(hiddenGames); // Store games with hidden UIDs
      setInstalledArray(installedGames);

      //await writeJSON("mainCache", JSON.stringify(json));
      console.log(json);
    } catch (error) {
      console.error(error);
    }
  };

  function initTileSize() {
    const tileSize = Number(localStorage.getItem("tileSize"));
    if (tileSize !== 0) {
      setTileSize(tileSize);
    } else {
      setTileSize(35);
      localStorage.setItem("tileSize", "35");
    }
  }

  useEffect(() => {
    initTileSize();
    setTheme();
    const initFunc = async () => {
      await fetchData();
      // const steamCreds = await getSteamCreds();
      // const npsso = await getNpsso();
      // importSteamLibrary(
      //   steamCreds?.ID,
      //   steamCreds?.APIKey,
      //   () => {},
      //   setIntegrationLoadCount
      // );
      // importPlaystationLibrary(
      //   npsso,
      //   () => {},
      //   () => {},
      //   setIntegrationLoadCount
      // );
    };
    initFunc();

    attachSSEListener(fetchData);
  }, []);

  useEffect(() => {
    if (!randomGameClicked) return;

    setRandomGameClicked(false);

    let targetArray = [];
    if (location.pathname === "/") {
      targetArray = dataArray;
    } else if (location.pathname === "/wishlist") {
      targetArray = wishlistArray;
    } else if (location.pathname === "/gameview") {
      // Determine the source array from where the game was previously picked
      targetArray = dataArray.some((game) => game.UID === location.state?.data)
        ? dataArray
        : wishlistArray;
    }

    if (targetArray.length > 0) {
      const randomIndex = Math.floor(Math.random() * targetArray.length);
      const uid = targetArray[randomIndex].UID;
      navigate(`/gameview`, {
        state: { data: uid, hidden: false },
      });
    }
  }, [randomGameClicked]);

  useEffect(() => {
    if (sortStateUpdate === true) {
      console.log("Sort State Update");
      fetchData();
      setSortStateUpdate(false);
    }
  }, [sortStateUpdate]);

  // These 2 are used to check if dialogs have been opened atleast once
  // This is necessary as if you leave the dialog in load state it wont save otherwise
  useEffect(() => {
    if (isAddGameDialogOpen) {
      setAddGameDialogHasBeenOpened(true);
    }
  }, [isAddGameDialogOpen]);
  useEffect(() => {
    if (isIntegrationsDialogOpen) {
      setIntegrationsPreviouslyOpened(true);
    }
  }, [isIntegrationsDialogOpen]);
  useEffect(() => {
    if (isWishlistAddDialogOpen) {
      setWishlistAddDialogPreviouslyOpened(true);
    }
  }, [isWishlistAddDialogOpen]);

  return (
    <>
      <CustomTitleBar>
        <BackButtonListener />
        {addGameDialogHasBeenOpened && <AddGameManuallyDialog />}
        {integrationsPreviouslyOpened && <Integrations />}
        {wishListAddDialogPreviouslyOpened && <WishlistDialog />}
        <Routes>
          <Route
            element={
              <LibraryView data={dataArray} hidden={false} viewText="Library" />
            }
            path="/"
          />

          <Route
            element={
              <LibraryView
                data={wishlistArray}
                hidden={false}
                viewText="Wishlist"
              />
            }
            path="/wishlist"
          />

          <Route
            element={
              <LibraryView
                data={installedArray}
                hidden={false}
                viewText="Installed"
              />
            }
            path="/installed"
          />

          <Route
            element={
              <LibraryView data={hiddenArray} hidden={true} viewText="Hidden" />
            }
            path="/hidden"
          />

          <Route element={<GameView />} path="/gameview" />
        </Routes>
        <Toaster />
      </CustomTitleBar>
    </>
  );
}

export default App;
