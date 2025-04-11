import { useEffect, useState } from "react";
import { setTheme } from "./lib/ToggleTheme";
import { data, Navigate, useLocation, useNavigate } from "react-router-dom";
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
import { fetchData } from "./lib/api/fetchBasicInfo";
import { initTileSize, setLastPath } from "./lib/initTileSize";
import { pickRandomGame } from "./lib/pickRandomGame";
import { useNavigationContext } from "./hooks/useNavigationContext";
import {
  importPlaystationLibrary,
  importSteamLibrary,
} from "./lib/api/libraryImports";
import { getNpsso, getSteamCreds } from "./lib/api/getCreds";
import {
  deleteCurrentlyFiltered,
  hideCurrentlyFiltered,
  unHideCurrentlyFiltered,
} from "./lib/api/filterGamesAPI";
import { UpdateManager } from "./components/AutoUpdate/AutoUpdateManager";
import { toast } from "./hooks/use-toast";
import Settings from "./components/Dialogs/settings/Settings";
function App() {
  const {
    sortType,
    sortOrder,
    setSortType,
    setSortOrder,
    setMetaData,
    setTileSize,
    sortStateUpdate,
    setSortStateUpdate,
    randomGameClicked,
    setRandomGameClicked,
    setIntegrationLoadCount,
    playingGame,
    deleteFilterGames,
    setDeleteFilterGames,
    hideFilterGames,
    setHideFilterGames,
  } = useSortContext();
  const location = useLocation();
  const [dataArray, setDataArray] = useState<any[]>([]);
  const [wishlistArray, setWishlistArray] = useState<any[]>([]);
  const [hiddenArray, setHiddenArray] = useState<any[]>([]);
  const [installedArray, setInstalledArray] = useState<any[]>([]);

  const navigate = useNavigate();

  const visibleUIDs = dataArray.map((game) => game.UID);
  console.log(visibleUIDs);

  const updateData = () => {
    fetchData(
      sortType,
      sortOrder,
      setDataArray,
      setMetaData,
      setSortOrder,
      setSortType,
      setWishlistArray,
      setHiddenArray,
      setInstalledArray
    );
  };

  useEffect(() => {
    if (hideFilterGames === null) return;

    if (hideFilterGames === "hide") {
      switch (location.pathname) {
        case "/library":
          hideCurrentlyFiltered(dataArray);
          break;
        case "/wishlist":
          hideCurrentlyFiltered(wishlistArray);
          break;
        case "/installed":
          hideCurrentlyFiltered(installedArray);
          break;
        default:
          break;
      }
    } else if (hideFilterGames === "unhide") {
      unHideCurrentlyFiltered(hiddenArray);
    }
    setHideFilterGames(null);
  }, [hideFilterGames]);

  useEffect(() => {
    if (!deleteFilterGames) return;

    switch (location.pathname) {
      case "/library":
        deleteCurrentlyFiltered(dataArray);
        break;
      case "/wishlist":
        deleteCurrentlyFiltered(wishlistArray);
        break;
      case "/installed":
        deleteCurrentlyFiltered(installedArray);
        break;
      case "/hidden":
        deleteCurrentlyFiltered(hiddenArray);
        break;
      default:
        console.warn("Unknown path for deletion:", location.pathname);
    }
    setDeleteFilterGames(false);
  }, [deleteFilterGames]);

  useEffect(() => {
    initTileSize(setTileSize);
    setLastPath(navigate);
    setTheme();
    const initFunc = async () => {
      await updateData();
      if (import.meta.env.MODE === "production") {
        const steamCreds = await getSteamCreds();
        const npsso = await getNpsso();
        importSteamLibrary(
          steamCreds?.ID,
          steamCreds?.APIKey,
          () => {},
          setIntegrationLoadCount,
          toast
        );
        importPlaystationLibrary(
          npsso,
          () => {},
          () => {},
          setIntegrationLoadCount,
          toast
        );
      }
    };
    initFunc();
    attachSSEListener(updateData);
  }, []);

  const { lastLibraryPath } = useNavigationContext();

  useEffect(() => {
    console.log("Rand", randomGameClicked);
    if (!randomGameClicked) return;
    setRandomGameClicked(false);
    const randomUID = pickRandomGame(
      lastLibraryPath,
      dataArray,
      wishlistArray,
      installedArray
    );
    console.log(randomUID);
    if (randomUID) {
      navigate(`/gameview`, {
        state: { data: randomUID, hidden: false },
      });
    }
  }, [randomGameClicked]);

  useEffect(() => {
    if (sortStateUpdate === true) {
      console.log("Sort State Update");
      updateData();
      setSortStateUpdate(false);
    }
  }, [sortStateUpdate]);

  useEffect(() => {
    if (playingGame === "") {
      const updateSteam = async () => {
        const steamCreds = await getSteamCreds();

        await importSteamLibrary(
          steamCreds?.ID,
          steamCreds?.APIKey,
          () => {},
          setIntegrationLoadCount,
          () => {}
        );
      };

      updateSteam();
    } else if (playingGame != null) {
      window.windowFunctions.updatePlayingGame(playingGame);
    }
  }, [playingGame]);

  useEffect(() => {
    if (location.pathname !== "/gameview")
      localStorage.setItem("lastPath", location.pathname);
  }, [location.pathname]);

  return (
    <>
      <CustomTitleBar>
        <BackButtonListener />
        <AddGameManuallyDialog />
        <Integrations />
        <WishlistDialog />
        <Settings />
        <Routes>
          <Route path="/" element={<Navigate to="/library" replace />} />
          <Route
            element={
              <LibraryView data={dataArray} hidden={false} viewText="Library" />
            }
            path="/library"
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
        <UpdateManager />
      </CustomTitleBar>
    </>
  );
}

export default App;
