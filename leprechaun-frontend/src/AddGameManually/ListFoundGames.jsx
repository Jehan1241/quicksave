import { useState } from "react";

function ListFoundGames(props) {
  const [commiting, setCommiting] = useState(false);

  const selectedGameClickHandler = async (appid) => {
    setCommiting(true);
    console.log(props.time);
    console.log(appid);
    console.log(props.SelectedPlatform);
    try {
      const response = await fetch("http://localhost:8080/InsertGameInDB", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          key: appid,
          platform: props.SelectedPlatform,
          time: props.time,
        }),
      });
      props.onGameAdded();
      setCommiting(false);
    } catch (error) {
      console.error("Error:", error);
      setCommiting(false);
    }
  };

  if (props.FoundGames === "") {
    return;
  } else {
    const data = JSON.parse(props.FoundGames.foundGames);
    if (Object.keys(data).length === 0) {
      return (
        <div className="mt-4 font-mono text-left my-2 bg-primary p-5 rounded-2xl h-[calc(100vh-400px)] w-1/2 overflow-scroll border-2 border-gray-700 flex justify-center hover:border-gray-500">
          No Games Found
        </div>
      );
    } else {
      return commiting ? (
        <div className="mt-4 font-mono text-left my-2 bg-primary p-5 rounded-2xl h-[calc(100vh-400px)] w-1/2 overflow-scroll border-2 border-gray-700 flex justify-center hover:border-gray-500">
          Adding Game...
        </div>
      ) : (
        <div className="mt-4 font-mono text-left my-2 bg-primary p-5 rounded-2xl h-[calc(100vh-400px)] w-1/2 overflow-scroll border-2 border-gray-700 flex justify-center hover:border-gray-500">
          <ul className="w-11/12 text-white">
            {Object.values(data).map((game) => (
              <li
                className="p-2 my-2 border-b-2 border-gray-700 hover:font-bold hover:scale-105 hover:border-white"
                key={game.appid}
                onClick={() => selectedGameClickHandler(game.appid)}
              >
                <div className="">{game.name}</div>
                <div className="text-right">
                  {new Date(game.date).getFullYear()}
                </div>
              </li>
            ))}
          </ul>
        </div>
      );
    }
  }
}

export default ListFoundGames;
