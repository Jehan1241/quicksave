import { useState } from "react";

import { useNavigate } from "react-router-dom";
import { IoIosArrowRoundBack, IoIosArrowRoundForward } from "react-icons/io";
import { GoPlus } from "react-icons/go";
import ImportPopUp from "./ImportPopUp";

function NavBar(props) {
  const [importClicked, setImportClicked] = useState(false);
  const [libraryClicked, setLibraryClicked] = useState(false);
  const navigate = useNavigate();

  const importClickHandler = () => {
    setImportClicked(!importClicked);
    setLibraryClicked(false);
    if (libraryClicked) {
      navigate("/");
    }
    console.log(importClicked);
  };

  const libraryClickHandler = () => {
    console.log("AA");
    setImportClicked(false);
    navigate("/");
  };

  /* In the next 2 funcs event.stopPropagation() needed to stop click through into next elements */
  const addGameManuallyClickHandler = (event) => {
    event.stopPropagation();
    setImportClicked(false);
    navigate("/AddGameManually");
  };

  const fromSteamClickHandler = (event) => {
    event.stopPropagation();
    setImportClicked(false);
    navigate("/AddGameSteam");
  };

  return (
    <>
      {/* Main BAR */}
      <div className="flex absolute z-30 flex-row justify-between items-center px-5 mx-10 mt-2 w-[calc(100vw-80px)] h-14 rounded-2xl shadow-md backdrop-blur-xl bg-primary/75">
        {/* Arrow Div */}
        <div className="flex justify-center items-center text-white">
          <button className="rounded-full hover:bg-gray-600/30">
            <IoIosArrowRoundBack size={40} onClick={() => navigate(-1)} />
          </button>
          <button className="rounded-full hover:bg-gray-600/30">
            <IoIosArrowRoundForward size={40} onClick={() => navigate(1)} />
          </button>
        </div>

        <input
          onChange={props.inputChangeHandler}
          className="px-3 my-auto h-7 text-white rounded-xl bg-gray-600/20"
          placeholder="Search"
        />
        {/* Add and Slider Div */}
        <div className="flex gap-3">
          <input
            defaultValue={40}
            onChange={props.sizeChangeHandler}
            type="range"
            min={25}
            max={80}
          ></input>
          <button
            onClick={importClickHandler}
            className={`text-left text-white rounded-full hover:bg-gray-600/30`}
          >
            <GoPlus className={`inline p-2 rounded-full`} size={40} />
          </button>
        </div>
      </div>
      {/* POP BAR */}
      {importClicked ? (
        <ImportPopUp
          importClickHandler={importClickHandler}
          fromSteamClickHandler={fromSteamClickHandler}
          addGameManuallyClickHandler={addGameManuallyClickHandler}
        />
      ) : null}
    </>
  );
}

export default NavBar;
