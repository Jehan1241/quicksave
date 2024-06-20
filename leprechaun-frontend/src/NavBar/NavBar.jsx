import { useState } from "react";
import { IoAddCircle } from "react-icons/io5";
import { FaPlaystation, FaPlus, FaSteam } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom } from "react-icons/si";
import { useNavigate } from "react-router-dom";
import { IoIosArrowRoundBack, IoIosArrowRoundForward } from "react-icons/io";
import { GoPlus } from "react-icons/go";

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
      <div className="flex absolute z-30 flex-row justify-between items-center px-5 mx-10 mt-2 w-[calc(100vw-80px)] h-14 rounded-2xl shadow-md backdrop-blur-xl bg-purple-600/20">
        {/* Arrow Div */}
        <div className="flex justify-center items-center text-white">
          <button className="rounded-full hover:bg-purple-600/30">
            <IoIosArrowRoundBack size={40} onClick={() => navigate(-1)} />
          </button>
          <button className="rounded-full hover:bg-purple-600/30">
            <IoIosArrowRoundForward size={40} onClick={() => navigate(1)} />
          </button>
        </div>

        <input
          onChange={props.inputChangeHandler}
          className="px-3 my-auto h-7 text-white rounded-xl bg-purple-600/20"
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
            className={`text-left text-white rounded-full hover:bg-purple-500/30`}
          >
            <GoPlus className={`inline p-2 rounded-full`} size={40} />
          </button>
        </div>
      </div>
      {/* POP BAR */}
      {importClicked ? (
        <div
          onClick={importClickHandler}
          className="flex overflow-hidden fixed z-20 w-screen h-screen bg-black/70"
        >
          <div className={`flex flex-col p-5 m-auto bg-gray-900 shadow-lg`}>
            <button
              className="text-gray-400 text-left text-m font-semibold  hover:text-lg hover:text-white min-h-[30px]"
              onClick={fromSteamClickHandler}
            >
              FROM STEAM <FaSteam className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold  hover:text-lg hover:text-white min-h-[30px]">
              FROM PLAYSTATION <FaPlaystation className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold  hover:text-lg hover:text-white min-h-[30px]">
              FROM EPIC <SiEpicgames className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold  hover:text-lg hover:text-white min-h-[30px]">
              FROM GOG <SiGogdotcom className="inline" size={20} />
            </button>
            <button
              onClick={addGameManuallyClickHandler}
              className="text-gray-400 text-left text-m font-semibold  hover:text-lg hover:text-white min-h-[30px]"
            >
              ADD GAME MANUALLY <IoAddCircle className="inline" size={20} />
            </button>
          </div>
        </div>
      ) : null}
    </>
  );
}

export default NavBar;
