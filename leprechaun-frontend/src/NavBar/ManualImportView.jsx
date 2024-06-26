import { useState, useEffect } from "react";
import ListFoundGames from "../AddGameManually/ListFoundGames";

function SteamImportView(props) {
  const [searchClicked, setSearchClicked] = useState(false);
  const [data, setData] = useState("");
  const [toSearch, setToSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [clickCount, setClickCount] = useState(0);
  const [selectedPlatform, setSelectedPlatform] = useState("");
  const [timePlayed, setTimePlayed] = useState(0);

  useEffect(() => {
    if (searchClicked) {
      fetchData(toSearch);
    }
  }, [clickCount]);

  useEffect(() => {
    if (data !== "") {
      console.log(Object.values(data));
    }
  }, [data]);

  const searchClickHandler = () => {
    setClickCount(clickCount + 1);
    setSearchClicked(true);
    const value = document.getElementById("SearchBar").value;
    const time = document.getElementById("timePlayed").value;
    const platform = document.getElementById("Platform").value;
    setSelectedPlatform(platform);
    setTimePlayed(time);
    setLoading(true);
    setToSearch(value);
  };

  const fetchData = async (toSearch) => {
    try {
      const response = await fetch("http://localhost:8080/IGDBsearch", {
        method: "POST",
        headers: { "Content-type": "application/json" },
        body: JSON.stringify({ NameToSearch: toSearch }),
      });
      setData(await response.json());
      setLoading(false);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="flex flex-col gap-1 p-4 mt-2 w-full h-full text-base rounded-xl">
      <div className="flex flex-row gap-10 w-full h-full">
        <div className="flex flex-col gap-2 border-2">
          <div className="flex flex-row gap-2 items-center">
            <p>Title</p>
            <input
              id="SearchBar"
              className="px-1 ml-7 w-52 rounded-lg bg-gray-500/20"
            ></input>
          </div>
          <div className="flex flex-row gap-2 items-center">
            <p>Time Played</p>
            <input
              id="timePlayed"
              className="px-1 w-52 h-6 rounded-lg bg-gray-500/20"
            ></input>
          </div>
          <div className="flex flex-row gap-2 items-center">
            <p>Platform</p>
            <input
              id="Platform"
              className="px-1 w-52 rounded-lg bg-gray-500/20"
            ></input>
          </div>
        </div>
        <div className="ml-auto w-full border-4">
          {loading ? (
            <div className="text-left bg-gameView h-[28vh] w-auto overflow-scroll flex justify-center">
              <p>Loading...</p>
            </div>
          ) : (
            <ListFoundGames
              FoundGames={data}
              SelectedPlatform={selectedPlatform}
              onGameAdded={props.onGameAdded}
              time={timePlayed}
            />
          )}
        </div>
      </div>
      <div className="flex justify-end mt-auto">
        <button
          className="w-32 h-10 rounded-lg border-2 bg-gameView hover:bg-gray-500/20"
          onClick={searchClickHandler}
        >
          Search
        </button>
      </div>
    </div>
  );

  /* <div className="flex flex-col p-4 mt-2 w-full h-full text-base rounded-xl">
      <div className="flex flex-row gap-4">
        <div className="flex flex-col gap-4">
          <p>Title</p>
          <p>Hours Played</p>
          <p>Platform</p>
        </div>
        <div className="flex flex-col gap-4">
          <input
            id="SearchBar"
            className="px-1 rounded-lg bg-gray-500/20"
          ></input>
          <input
            id="timePlayed"
            className="px-1 rounded-lg bg-gray-500/20"
          ></input>
          <input
            id="Platform"
            className="px-1 rounded-lg bg-gray-500/20"
          ></input>
        </div>
        <div>
          {loading ? (
            <div className="text-left bg-gameView rounded-2xl h-[28vh] w-[30vw] overflow-scroll border-gray-700 flex justify-center hover:border-gray-500">
              <p>Loading...</p>
            </div>
          ) : (
            <ListFoundGames
              FoundGames={data}
              SelectedPlatform={selectedPlatform}
              onGameAdded={props.onGameAdded}
              time={timePlayed}
            />
          )}
        </div>
      </div>

      <div className="flex justify-end mt-auto">
        <button
          className="w-32 h-10 rounded-lg border-2 bg-primary"
          onClick={searchClickHandler}
        >
          Search
        </button>
      </div>
    </div> */
}

export default SteamImportView;
