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
    <div className="relative w-[calc(100vw-104px)] text-white font-mono h-screen overflow-scroll text-center bg-gameView ml-1 rounded-l-3xl">
      <img
        className="absolute z-0 w-full opacity-20"
        src="https://i.redd.it/imjkslu7qjia1.jpg"
      />
      {/* Horizontal FLEX Holder */}
      <div className="flex flex-row">
        {/* LOGO AND Description DIV */}
        <div className=" m-2 w-1/3 h-[calc(100vh-15px)] opacity-75 rounded-3xl">
          <div
            className={`flex flex-col p-2 h-1/3 rounded-t-3xl border-2 border-gray-500 ${
              infoVisible ? "" : "rounded-3xl"
            }`}
          >
            <img
              className="object-scale-down h-5/6"
              src="https://cdn.cloudflare.steamstatic.com/steam/apps/292030/logo.png?t=1693590448"
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
          {infoVisible ? (
            <DisplayInfo
              data={data}
              tags={tagsArray}
              companies={companiesArray}
            />
          ) : (
            ""
          )}
        </div>
        {/* Image Div */}
        {imageVisible ? <DisplayImage screenshots={screenshotsArray} /> : ""}
      </div>
    </div>
  );
}

export default GameView;
