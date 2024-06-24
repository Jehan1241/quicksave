import { IoAddCircle } from "react-icons/io5";
import { FaPlaystation, FaPlus, FaSteam } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom } from "react-icons/si";

function ImportPopUp(props) {
  return (
    <div
      onClick={props.importClickHandler}
      className="flex overflow-hidden fixed z-20 w-screen h-screen bg-black/80"
    >
      <div
        className={`flex flex-col p-5 m-auto text-white bg-gray-900 shadow-lg`}
      >
        <button
          className="text-gray-400 text-left text-m font-semibold  hover:text-lg  min-h-[30px]"
          onClick={props.fromSteamClickHandler}
        >
          FROM STEAM <FaSteam className="inline" size={20} />
        </button>
        <button
          onClick={props.addGameManuallyClickHandler}
          className="text-gray-400 text-left text-m font-semibold  hover:text-lg min-h-[30px]"
        >
          ADD GAME MANUALLY <IoAddCircle className="inline" size={20} />
        </button>
      </div>
    </div>
  );
}

export default ImportPopUp;
