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
    sortStateUpdate,
    setSortStateUpdate,
    randomGameClicked,
    setRandomGameClicked,
  } = useSortContext();
  const location = useLocation();
  const [dataArray, setDataArray] = useState<any[]>([]);
  const [wishlistArray, setWishlistArray] = useState<any[]>([]);
  const [hiddenArray, setHiddenArray] = useState<any[]>([]);
  const [installedArray, setInstalledArray] = useState<any[]>([]);

  const navigate = useNavigate();

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
    attachSSEListener(updateData);
  }, []);

  useEffect(() => {
    if (!randomGameClicked) return;
    setRandomGameClicked(false);
    const randomUID = pickRandomGame(location, dataArray, wishlistArray);
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

  return (
    <>
      <CustomTitleBar>
        <BackButtonListener />
        <AddGameManuallyDialog />
        <Integrations />
        <WishlistDialog />
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
                hidden={true}
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
