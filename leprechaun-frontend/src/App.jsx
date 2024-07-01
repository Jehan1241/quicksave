import { useState, useEffect } from "react";
import NavBar from "./NavBar/NavBar";
import LibraryView from "./LibraryView/LibraryView";
import { Route, Routes, useLocation } from "react-router-dom";
import AddGameManually from "./AddGameManually/AddGameManually";
import AddGameSteam from "./AddGameSteam/AddGameSteam";
import GameView from "./GameView/GameView";

function App() {
  const [metaData, setMetaData] = useState([]);
  const [searchText, setSearchText] = useState("");
  const location = useLocation();
  const state = location.state;
  const [tileSize, setTileSize] = useState("");

  const NavBarInputChangeHanlder = (e) => {
    const text = e.target.value;
    setSearchText(text.toLowerCase());
  };

  /* To initially set Tile Size to xxx */
  useEffect(() => {
    setTileSize(40);
  }, []);

  const sizeChangeHandler = (e) => {
    setTileSize(e.target.value);
  };

  const fetchData = async () => {
    console.log("fetch");
    try {
      const response = await fetch("http://localhost:8080/getBasicInfo");
      const json = await response.json();
      setMetaData(json.MetaData);
      console.log(json.Screenshots);
      console.log("Run");
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const DataArray = Object.values(metaData);

  return (
    <div className="flex flex-col w-screen h-screen">
      <NavBar
        onGameAdded={fetchData}
        inputChangeHandler={NavBarInputChangeHanlder}
        sizeChangeHandler={sizeChangeHandler}
      />
      <Routes>
        <Route
          element={
            <LibraryView
              tileSize={tileSize}
              searchText={searchText}
              data={DataArray}
            />
          }
          path="/"
        />
        <Route
          element={<AddGameManually onGameAdded={fetchData} />}
          path="AddGameManually"
        />
        <Route element={<AddGameSteam />} path="AddGameSteam" />
        <Route
          element={<GameView uid={state?.data} onDelete={fetchData} />}
          path="gameview"
        />
      </Routes>
    </div>
  );
}

export default App;
