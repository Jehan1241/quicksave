import { useState } from "react";
import { IoAddCircle, IoLibraryOutline } from "react-icons/io5";
import { CiImport } from "react-icons/ci";
import { FaPlaystation, FaSteam } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom } from "react-icons/si";
import { useNavigate } from "react-router-dom";

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
    setImportClicked(false);
    if (libraryClicked) {
      navigate("/");
    } else {
      navigate("/LibraryView");
    }
    setLibraryClicked(!libraryClicked);
    console.log(libraryClicked);
  };

  const addGameManuallyClickHandler = () => {
    importClickHandler();
    navigate("/AddGameManually");
  };

  const fromSteamClickHandler = () => {
    importClickHandler();
    navigate("/AddGameSteam");
  };

  return (
    <div className="flex flex-row">
      {/* Main BAR */}
      <div className="flex z-20 flex-col p-2 w-24 h-screen rounded-r-3xl shadow-lg bg-zinc-950">
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
        <button
          onClick={importClickHandler}
          className={`text-gray-400 hover:text-xl text-left hover:text-white my-3 min-h-16 ${
            importClicked ? "border-r-4 border-white translate-x-3" : "pl-2"
          }`}
        >
          <CiImport
            className={`inline transition duration-75 hover:scale-125 ${
              importClicked ? "text-white scale-125 translate-x-4" : ""
            }`}
            size={50}
          />
        </button>
      </div>

      {/* POP BAR */}
      <div
        className={`bg-primary z-10 fixed flex flex-col h-screen w-64 p-5 rounded-3xl ml-1 ${
          importClicked ? "translate-x-24" : "-translate-x-96"
        } ease-in-out duration-300`}
      >
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
  );
}

export default NavBar;
