import { useState } from "react";
import { BiCaretRight } from "react-icons/bi";
import { IoAdd, IoAddCircle, IoLibraryOutline } from "react-icons/io5";
import { CiImport } from "react-icons/ci";
import { FaPlaystation, FaSteam } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom } from "react-icons/si";
import { IconContext } from "react-icons/lib";

function NavBar() {
  const [importClicked, setImportClicked] = useState(false);
  const [libraryClicked, setLibraryClicked] = useState(false);

  const importClickHandler = () => {
    setImportClicked(!importClicked);
    setLibraryClicked(false);
    console.log(importClicked);
  };
  const libraryClickHandler = () => {
    setImportClicked(false);
    setLibraryClicked(!libraryClicked);
    console.log(libraryClicked);
  };

  return (
    <div className="flex flex-row">
      {/* Main BAR */}
      <div className="flex flex-col bg-zinc-950 h-screen w-24 shadow-lg p-2 z-20 rounded-r-3xl">
        <button
          onClick={libraryClickHandler}
          className={`text-gray-400 hover:text-xl text-left hover:text-white my-3 min-h-16 ${
            libraryClicked ? "translate-x-3 border-r-4 border-white" : "pl-2"
          }`}
        >
          <IoLibraryOutline
            size={50}
            className={`inline transition duration-150 hover:scale-125 ${
              libraryClicked ? "scale-125 text-white translate-x-3" : ""
            }`}
          />
        </button>
        <button
          onClick={importClickHandler}
          className={`text-gray-400 hover:text-xl text-left hover:text-white my-3 min-h-16 ${
            importClicked ? "translate-x-3 border-r-4 border-white" : "pl-2"
          }`}
        >
          <CiImport
            className={`inline transition duration-150 hover:scale-125 ${
              importClicked ? "scale-125 text-white translate-x-4" : ""
            }`}
            size={50}
          />
        </button>
      </div>

      {/* POP BAR */}
      <div
        className={`bg-primary z-0 flex flex-col h-screen w-64 p-5 rounded-3xl ml-1 ${
          importClicked ? "translate-x-0" : "-translate-x-96"
        } ease-in-out duration-300`}
      >
        <button className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]">
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
        <button className="text-gray-400 text-left text-m font-semibold font-mono hover:text-lg hover:text-white min-h-[30px]">
          ADD GAME MANUALLY <IoAddCircle className="inline" size={20} />
        </button>
      </div>
    </div>
  );
}

export default NavBar;
