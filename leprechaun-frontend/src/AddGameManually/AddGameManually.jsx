import { useState, useEffect } from "react";
import ListFoundGames from "./ListFoundGames";

function AddGameManually(props) {
  const [searchClicked, setSearchClicked] = useState(false);
  const [data, setData] = useState("");
  const [toSearch, setToSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [clickCount, setClickCount] = useState(0);
  const [platforms, setPlatforms] = useState("");
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

  const getPlatforms = async () => {
    try {
      const response = await fetch("http://localhost:8080/GetPlatforms");
      const json = await response.json();
      setPlatforms(json.platforms);
      console.log("Platforms");
    } catch (error) {
      console.error(error);
    }
  };

  const optionSelectHandler = (platform) => {
    setSelectedPlatform(platform);
  };

  useEffect(() => {
    getPlatforms();
  }, []);

  return (
    <div className="overflow-hidden relative p-12 ml-1 h-screen text-center text-white rounded-l-3xl bg-gameView">
      {/* Spacer DIV */}
      <div className="h-[4%]" />
      <div className="p-6 w-1/2 rounded-2xl border-2 border-gray-700 bg-primary hover:border-gray-500">
        <p className="text-2xl text-left text-white">Import Game Manually</p>
        <div className="m-6 text-left">
          <p className="inline m-2 text-white">Enter Game Title</p>
          <input
            id="SearchBar"
            className="p-1 ml-9 w-72 rounded-xl border-2 border-gray-700 bg-gameView"
          />
          <br />
          <br />
          <p className="inline m-2">Platform</p>
          <select
            className="p-2 ml-28 w-72 rounded-xl border-2 border-gray-700 bg-gameView"
            onChange={(e) => optionSelectHandler(e.target.value)}
          >
            {Object.values(platforms).map((platform) => (
              <option id={platform.id} key={platform.id} value={platform.name}>
                {platform.name}
              </option>
            ))}
          </select>
          <br />
          <br />
          <p className="inline m-2 text-white">Hours Played</p>
          <input
            id="timePlayed"
            onChange={(e) => {
              setTimePlayed(e.target.value);
            }}
            className="p-1 ml-20 w-24 rounded-xl border-2 border-gray-700 bg-gameView"
            onInput={(e) => {
              const value = e.target.value;
              const regex = /^[0-9]+$/;
              if (!regex.test(value)) {
                e.target.value = value.replace(/[^0-9]/g, "");
              }
            }}
          />
          <div className="text-right">
            <button
              onClick={searchClickHandler}
              className={`h-10 w-32 rounded-lg bg-gameView border-2 border-gray-700 hover:font-extrabold hover:border-white ${
                searchClicked ? "text" : ""
              }`}
            >
              Search
            </button>
          </div>
        </div>
      </div>
      <div>
        {loading ? (
          <div className="mt-4  text-left my-2 p-5 bg-primary  rounded-2xl h-[calc(100vh-400px)] w-1/2 overflow-scroll border-2 border-gray-700 flex justify-center hover:border-gray-500">
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
  );
}

export default AddGameManually;
