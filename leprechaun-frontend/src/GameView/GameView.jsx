import { useState } from "react";
import { IoIosInformationCircleOutline } from "react-icons/io";
import { IoImagesOutline } from "react-icons/io5";
import DisplayInfo from "./DisplayInfo";
import DisplayImage from "./DisplayImage";

function GameView(props) {
  const data = props.data;
  const tags = props.tags;
  const tagsArray = Object.values(tags);
  const companiesArray = Object.values(props.companies);
  const screenshotsArray = Object.values(props.screenshots);

  const [imageVisible, setImageVisible] = useState(true);
  const [infoVisible, setInfoVisible] = useState(true);

  const informationClickHandler = () => {
    setInfoVisible(!infoVisible);
  };
  const imageClickHandler = () => {
    setImageVisible(!imageVisible);
  };

  return (
    <div className="overflow-scroll relative h-screen font-mono text-center text-white bg-gameView bg-[src]">
      {/* <img
        className="absolute z-0 w-full opacity-20"
        src={"/leprechaun-backend/" + screenshotsArray[0]}
      /> */}
      {/* Horizontal FLEX Holder */}
      <div className="flex flex-row">
        {/* LOGO AND Description DIV */}
        <div className=" m-2 w-1/3 h-[calc(100vh-15px)] rounded-3xl ">
          <div
            className={`flex flex-col p-2 h-1/3 rounded-t-3xl border-2 border-gray-500 backdrop-blur-md bg-black/20 ${
              infoVisible ? "" : "rounded-3xl"
            }`}
          >
            <img
              className="object-scale-down m-2 h-5/6 text-5xl font-extrabold"
              src="https://cdn.cloudflare.steamstatic.com/steam/apps/292030/logo.png?t=1693590448"
              alt={data.Name}
            />
            <div>
              <button className="mx-1 w-4/12 h-12 rounded-2xl border-2 border-gray-500 bg-primary">
                Play
              </button>
              <button className="mx-1 w-4/12 h-12 rounded-2xl border-2 border-gray-500 bg-primary">
                Customize
              </button>
              <button
                onClick={informationClickHandler}
                className="mx-1 w-8 h-8 rounded-2xl border-2 border-gray-500 bg-primary"
              >
                <IoIosInformationCircleOutline className="inline" size={20} />
              </button>
              <button
                onClick={imageClickHandler}
                className="mx-1 w-8 h-8 rounded-2xl border-2 border-gray-500 bg-primary"
              >
                <IoImagesOutline className="inline" size={18} />
              </button>
            </div>
          </div>
          <DisplayInfo
            data={data}
            tags={tagsArray}
            companies={companiesArray}
            visible={infoVisible}
          />
        </div>
        {/* Image Div */}
        <DisplayImage screenshots={screenshotsArray} visible={imageVisible} />
      </div>
    </div>
  );
}

export default GameView;
