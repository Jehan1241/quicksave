import { useState, useEffect } from "react";
import NavBar from "./NavBar/NavBar";
import LibraryView from "./LibraryView/LibraryView";
import { Route, Routes } from "react-router-dom";
import AddGameManually from "./AddGameManually/AddGameManually";
import AddGameSteam from "./AddGameSteam/AddGameSteam";
import GameView from "./GameView/GameView";

function App() {
  const [metaData, setMetaData] = useState([]);
  const [tags, setTags] = useState([]);
  const [companies, setCompanies] = useState([]);
  const [screenshots, setScreenshots] = useState([]);

  const fetchData = async () => {
    try {
      const response = await fetch("http://localhost:8080/getBasicInfo");
      const json = await response.json();
      setMetaData(json.MetaData);
      setTags(json.Tags);
      setCompanies(json.Companies);
      setScreenshots(json.Screenshots);
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
    <div className="flex overflow-hidden flex-col w-screen h-screen bg-[url('leprechaun-backend/screenshots/a70700f55306ef99764b30c8de7dd78d/a70700f55306ef99764b30c8de7dd78d-0.jpeg')]">
      <NavBar />
      <Routes>
        <Route element={<LibraryView data={DataArray} />} path="LibraryView" />
        <Route
          element={<AddGameManually onGameAdded={fetchData} />}
          path="AddGameManually"
        />
        <Route element={<AddGameSteam />} path="AddGameSteam" />
        {DataArray.map((item) => (
          <Route
            element={
              <GameView
                screenshots={screenshots[item.UID] ? screenshots[item.UID] : {}}
                companies={companies[item.UID] ? companies[item.UID] : {}}
                tags={tags[item.UID] ? tags[item.UID] : {}}
                data={item}
              />
            }
            path={`/LibraryView/GameView/${item.UID}`}
          />
        ))}
      </Routes>
    </div>
  );
}

export default App;
