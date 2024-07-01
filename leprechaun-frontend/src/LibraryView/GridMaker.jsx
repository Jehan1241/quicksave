import { useState } from "react";
import { FaPlaystation, FaSteamSymbol } from "react-icons/fa";
import { SiEpicgames, SiGogdotcom, SiSteam } from "react-icons/si";
import { useNavigate } from "react-router-dom";

function GridMaker(props) {
  const [imageError, setImageError] = useState(false);
  const navigate = useNavigate();
  const tileSize = props.tileSize / 30;

  const tileCilckHandler = () => {
    console.log(props.uid);
    navigate(`gameview`, {
      state: { data: props.uid },
    });
  };

  const style = {
    width: `calc(11rem * ${tileSize})`,
    height: `calc(16rem * ${tileSize})`,
  };

  return (
    <div
      onClick={tileCilckHandler}
      className="inline-flex justify-center items-center mx-5 mt-5 mb-5 transition duration-300 ease-in-out hover:scale-105"
      style={style}
    >
      <div className="flex relative flex-col group">
        <img
          src={
            !imageError
              ? `http://localhost:8080/cover-art/${props.cover}`
              : "http://localhost:8080/cover-art/default/default.jpg"
          }
          onError={() => setImageError(true)}
          style={style}
          className="object-cover text-sm text-white"
        />

        <p className="overflow-hidden text-xs text-center text-white truncate whitespace-nowrap opacity-0 transition duration-300 ease-in-out group-hover:opacity-100">
          {props.platform === "PlayStation 5" && (
            <FaPlaystation size={25} className="absolute -bottom-4 -right-5" />
          )}
          {props.platform === "Steam" && (
            <FaSteamSymbol size={25} className="absolute -right-4 -bottom-4" />
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
