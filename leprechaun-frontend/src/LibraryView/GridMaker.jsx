import { useState } from "react";
import { FaPlaystation, FaSteamSymbol } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom, SiSteam } from "react-icons/si";
import { useNavigate } from "react-router-dom";

function GridMaker(props) {
  const [imageError, setImageError] = useState(false);
  const navigate = useNavigate();

  const tileCilckHandler = () => {
    navigate(`GameView/${props.uid}`);
  };

  return (
    <div
      onClick={tileCilckHandler}
      className="inline-flex justify-center items-center mx-5 mt-5 mb-5 w-44 h-64 hover:rounded-2xl hover:scale-110"
    >
      <div className="flex relative flex-col group hover:rounded-2xl">
        <img
          /* src={
            !imageError
              ? `http://localhost:8080/cover-art/${props.cover}`
              : "/leprechaun-backend/coverArt/default/default.jpg"
          } */
          src={`http://localhost:8080/cover-art/${props.cover}`}
          onError={() => setImageError(true)}
          className={`object-cover w-44 h-64 text-sm text-white group-hover:rounded-2xl`}
        />
        <span className="text-white">{props.name}</span>

        <p className="hidden overflow-hidden text-xs text-center text-white truncate whitespace-nowrap group-hover:block">
          {props.platform === "PS5" && (
            <FaPlaystation size={25} className="absolute -bottom-4 -right-5" />
          )}
          {props.platform === "Steam" && (
            <FaSteamSymbol
              size={25}
              className="absolute -right-2 -bottom-2 rounded-full bg-gameView"
            />
          )}
          {props.platform === "GOG" && (
            <SiGogdotcom size={25} className="absolute -right-4 -bottom-4" />
          )}
          {props.platform === "Epic" && (
            <SiEpicgames
              size={25}
              className="absolute -right-4 -bottom-5 border"
            />
          )}
        </p>
      </div>
    </div>
  );
}

export default GridMaker;
