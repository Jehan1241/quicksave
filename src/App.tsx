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
import { fetchData } from "./lib/api/fetchBasicInfo";
import { initTileSize } from "./lib/initTileSize";
import { pickRandomGame } from "./lib/pickRandomGame";

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
    setCacheBuster,
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

  const updateData = () => {
    fetchData(
      sortType,
      sortOrder,
      tileSize,
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
    initTileSize(setTileSize);
    setTheme();
    const initFunc = async () => {
      await updateData();
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

    attachSSEListener(() => {
      updateData();
    }, setCacheBuster);
  }, []);

  useEffect(() => {
    if (!randomGameClicked) return;
    setRandomGameClicked(false);
    const targetArray = pickRandomGame(location, dataArray, wishlistArray);
    if (targetArray) {
      navigate(`/gameview`, {
        state: { data: targetArray, hidden: false },
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
