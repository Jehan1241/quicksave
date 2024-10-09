import { useState } from "react";
import SteamImportView from "./SteamImportView";
import ManualImportView from "./ManualImportView";

function ImportPopUp(props) {
  const [steamClicked, setSteamClicked] = useState(false);
  const [manuallyClicked, setManuallyClicked] = useState(true);

  const fromSteamClickHandler = () => {
    setSteamClicked(true);
    setManuallyClicked(false);
  };

  const manuallyClickHandler = () => {
    setManuallyClicked(true);
    setSteamClicked(false);
  };

  return (
    <div
      className="flex overflow-hidden fixed z-20 w-screen h-screen text-white bg-black/80"
      onClick={props.importClickHandler}
    >
      <div
        className="flex flex-col p-5 m-auto w-2/3 h-2/3 text-2xl bg-gameView"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="pb-2 mb-2 border-b-2">Import Games</div>
        <div className="flex flex-row">
          <button
            onClick={manuallyClickHandler}
            className={`mx-2 h-10 text-lg ${
              manuallyClicked ? "border-b-2" : ""
            }`}
          >
            Manually
          </button>
          <button
            onClick={fromSteamClickHandler}
            className={`mx-2 h-10 text-lg ${steamClicked ? "border-b-2" : ""}`}
          >
            From Steam
          </button>
        </div>
        {steamClicked ? (
          <SteamImportView onGameAdded={props.onGameAdded} />
        ) : null}
        {manuallyClicked ? (
          <ManualImportView onGameAdded={props.onGameAdded} />
        ) : null}
      </div>
    </div>
  );
}

export default ImportPopUp;
