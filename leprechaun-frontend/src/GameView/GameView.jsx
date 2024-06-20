import { useEffect, useState } from "react";
import DisplayInfo from "./DisplayInfo";
import DisplayImage from "./DisplayImage";

function GameView(props) {
  const [companies, setCompanies] = useState("");
  const [tags, setTags] = useState("");
  const [screenshots, setScreenshots] = useState("");
  const [metadata, setMetadata] = useState("");

  const fetchData = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/GameDetails?uid=${props.uid}`
      );
      const json = await response.json();

      // Destructure metadata from the JSON response
      const { companies, tags, screenshots, m: metadata } = json.metadata;

      // Set state correctly
      setCompanies(companies[props.uid]); // Access companies by UID
      setTags(tags[props.uid]); // Access tags by UID
      setMetadata(metadata[props.uid]); // Access metadata by UID
      setScreenshots(screenshots[props.uid]); // Access screenshots by UID
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const UID = props.uid;
  /*   const data = props.data; */
  /*   const tags = props.tags;
  const tagsArray = Object.values(tags);
  const companiesArray = Object.values(props.companies);
  const screenshotsArray = Object.values(props.screenshots); */

  const tagsArray = Object.values(tags);
  const companiesArray = Object.values(companies);
  const screenshotsArray = Object.values(screenshots);

  return (
    <>
      <img
        className="absolute top-0 right-0 w-screen h-screen opacity-20 blur-md"
        src={"http://localhost:8080/screenshots/" + screenshots[0]}
      />
      <div className="overflow-y-auto relative h-screen text-center text-white">
        {/* Spacer Div */}
        <div className="h-[4%]"></div>
        {/* Name and Time Flex */}
        <div className="flex flex-row justify-between items-center px-3 mx-10 my-2 h-24 rounded-2xl">
          <div className="text-3xl font-bold">{metadata.Name}</div>
          <div className="text-2xl">
            Time Played : {metadata.TimePlayed} Hrs
          </div>
        </div>
        {/* Horizontal FLEX Holder */}
        <div className="flex flex-row mx-10">
          {/* LOGO AND Description DIV */}
          <div className="m-2 w-1/3 h-full rounded-3xl">
            <DisplayInfo
              data={metadata}
              tags={tagsArray}
              companies={companiesArray}
            />
          </div>
          {/* Image Div */}
          <DisplayImage screenshots={screenshotsArray} />
        </div>
      </div>
    </>
  );
}

export default GameView;
