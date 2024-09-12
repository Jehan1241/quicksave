import { useState, useEffect } from "react";
import NavBar from "./NavBar/NavBar";
import LibraryView from "./LibraryView/LibraryView";
import { Route, Routes, useLocation } from "react-router-dom";
import GameView from "./GameView/GameView";

function App() {
  const [metaData, setMetaData] = useState([]);
  const [searchText, setSearchText] = useState("");
  const location = useLocation();
  const state = location.state;
  const [tileSize, setTileSize] = useState("");
  const [sortType, setSortType] = useState("default");
  const [sortOrder, setSortOrder] = useState("default");
  const [sse, setSse] = useState(null); // State to hold SSE connection

  const sortTypeChangeHandler = (type, order) => {
    setSortType(type);
    setSortOrder(order);
  };

  const NavBarInputChangeHandler = (e) => {
    const text = e.target.value;
    setSearchText(text.toLowerCase());
  };

  useEffect(() => {
    setTileSize(40);
  }, []);

  const sizeChangeHandler = (e) => {
    setTileSize(e.target.value);
  };

  useEffect(() => {
    const eventSource = new EventSource(
      "http://localhost:8080/sse-steam-updates"
    );

    eventSource.onmessage = (event) => {
      console.log("SSE message received:", event.data); // Log SSE message directly
      fetchData();
    };
    eventSource.onerror = (error) => {
      console.error("SSE Error:", error);
      // Handle SSE connection errors here
    };
    setSse(eventSource); // Save eventSource object in state
    return () => {
      eventSource.close(); // Clean up SSE connection on component unmount
    };
  }, []); // Empty dependency array ensures this effect runs only once

  const fetchData = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/getBasicInfo?type=${sortType}&order=${sortOrder}`
      );
      const json = await response.json();
      setMetaData(json.MetaData);
      setSortOrder(json.SortOrder);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [sortType, sortOrder]);

  const DataArray = Object.values(metaData);

  return (
    <div className="flex flex-col w-screen h-screen">
      <NavBar
        onGameAdded={fetchData}
        inputChangeHandler={NavBarInputChangeHandler}
        sizeChangeHandler={sizeChangeHandler}
        sortTypeChangeHandler={sortTypeChangeHandler}
        sortOrder={sortOrder}
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
          element={<GameView uid={state?.data} onDelete={fetchData} />}
          path="gameview"
        />
      </Routes>
    </div>
  );
}

export default App;
