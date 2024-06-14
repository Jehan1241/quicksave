import { useState } from "react";
import { IoAddCircle, IoLibraryOutline } from "react-icons/io5";
import { CiImport } from "react-icons/ci";
import { FaPlaystation, FaSteam } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom } from "react-icons/si";
import { useNavigate } from "react-router-dom";
import AddGameManually from "../AddGameManually/AddGameManually";

function NavBar() {
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
    navigate("/LibraryView");
  };

  const addGameManuallyClickHandler = (event) => {
    event.stopPropagation();
    setImportClicked(false);
    navigate("/AddGameManually");
  };

  const fromSteamClickHandler = () => {
    navigate("/AddGameSteam");
  };

  return (
    <>
      {/* Main BAR */}
      <div className="flex flex-row justify-between items-center px-20 py-2 w-screen shadow-lg bg-primary">
        <button
          onClick={libraryClickHandler}
          className={`text-gray-400 hover:text-xl text-left hover:text-white my-3 min-h-16 ${
            libraryClicked ? "border-r-4 border-white translate-x-3" : "pl-2"
          }`}
        >
          <IoLibraryOutline
            size={50}
            className={`inline transition duration-75 hover:scale-125 ${
              libraryClicked ? "text-white scale-125 translate-x-3" : ""
            }`}
          />
        </button>
        <input className="my-auto h-8"></input>
        <button
          onClick={importClickHandler}
          className={`my-3 text-left text-gray-400 hover:text-xl hover:text-white min-h-16`}
        >
          <CiImport
            className={`inline transition duration-75 hover:scale-125`}
            size={50}
          />
        </button>
      </div>
      {/* POP BAR */}

      {importClicked ? (
        <div
          onClick={libraryClickHandler}
          className="flex overflow-hidden fixed z-20 w-screen h-screen bg-black/70"
        >
          <div className={`flex flex-col p-5 m-auto bg-gray-900 shadow-lg`}>
            <button
              className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]"
              onClick={fromSteamClickHandler}
            >
              FROM STEAM <FaSteam className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]">
              FROM PLAYSTATION <FaPlaystation className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]">
              FROM EPIC <SiEpicgames className="inline" size={20} />
            </button>
            <button className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]">
              FROM GOG <SiGogdotcom className="inline" size={20} />
            </button>
            <button
              onClick={addGameManuallyClickHandler}
              className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]"
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
