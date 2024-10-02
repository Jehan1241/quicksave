import { useState } from "react";

function AddGamePath(props) {
  const [gamePath, setGamePath] = useState(null);

  return (
    <div
      className="flex overflow-hidden fixed z-20 w-screen h-screen text-white bg-black/80"
      onClick={props.addGamePathClickHandler}
    >
      <div
        className="flex flex-col p-5 m-auto w-1/3 h-1/3 text-2xl bg-gameView"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="pb-2 mb-2 border-b-2">Add Game Path</div>

        <div className="flex flex-row gap-2 p-4 mt-2 w-full h-full text-base rounded-xl">
          <p>Game Path</p>
          <input
            id="gamePath"
            className="px-1 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
          ></input>
        </div>
        <button
          className="ml-auto w-32 h-16 text-base rounded-lg border-2 bg-primary"
          onClick={() =>
            props.sendGamePathtoDB(document.getElementById("gamePath").value)
          }
        >
          Add Path
        </button>
      </div>
    </div>
  );
}

export default AddGamePath;
