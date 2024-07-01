import { useState, useEffect, useRef } from "react";
import { useClickAway } from "react-use";
import { useNavigate } from "react-router-dom";
import { FaSortAmountUpAlt } from "react-icons/fa";
import { IoIosArrowRoundBack, IoIosArrowRoundForward } from "react-icons/io";
import { GoPlus } from "react-icons/go";
import ImportPopUp from "./ImportPopUp";

function NavBar(props) {
  const [importClicked, setImportClicked] = useState(false);
  const [libraryClicked, setLibraryClicked] = useState(false);
  const [sortClicked, setSortClicked] = useState(false);
  const navigate = useNavigate();
  const dropdownRef = useRef(null);

  /* Custom Hook from React-Use Library for Clickaway */
  useClickAway(dropdownRef, () => {
    setSortClicked(false);
  });

  const sortClickHandler = () => {
    setSortClicked(!sortClicked);
  };

  const sortOptionSelect = async (type) => {
    console.log(type);
    props.sortTypeChangeHandler(type);
    try {
      const response = await fetch(`http://localhost:8080/sort?type=${type}`);
    } catch (error) {
      console.error(error);
    }
  };

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

        <div className="flex flex-row gap-2 justify-center items-center text-white">
          <input
            onChange={props.inputChangeHandler}
            className="px-3 my-auto h-7 rounded-xl bg-gray-600/20"
            placeholder="Search"
          />
          <div className="flex relative flex-col">
            <button>
              <FaSortAmountUpAlt
                size={18}
                onClick={sortClickHandler}
                className={`duration-150 ease-in-out ${
                  sortClicked ? "rotate-180" : ""
                }`}
              />
            </button>
            {sortClicked ? (
              <div
                ref={dropdownRef}
                className="flex absolute top-10 flex-col gap-2 p-1 text-sm rounded-lg bg-gameView"
              >
                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect("Name")}
                >
                  Alphabetical
                </button>

                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect("TimePlayed")}
                >
                  Time Played
                </button>

                <button
                  className="px-3 py-2 rounded-lg hover:bg-gray-600/30"
                  onClick={() => sortOptionSelect("AggregatedRating")}
                >
                  Rating
                </button>
              </div>
            ) : null}
          </div>
        </div>
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
          onGameAdded={props.onGameAdded}
          importClickHandler={importClickHandler}
          fromSteamClickHandler={fromSteamClickHandler}
          addGameManuallyClickHandler={addGameManuallyClickHandler}
        />
      ) : null}
    </>
  );
}

export default NavBar;
