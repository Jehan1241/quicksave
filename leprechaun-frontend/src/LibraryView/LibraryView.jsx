import { CiSearch } from "react-icons/ci";
import GridMaker from "./GridMaker";
import GridMaker2 from "./GridMaker";
import { FaFilter } from "react-icons/fa";
import { useState } from "react";

function LibraryView(props) {
  const data = props.data;
  const searchText = props.searchText;

  return (
    <div className="overflow-y-auto min-h-screen text-center bg-gameView">
      {/* Spacer Div */}
      <div className="h-16"></div>
      {data.map((item) =>
        item.Name.toLowerCase()
          .replace("'", "")
          .replace("’", "")
          .replace("®", "")
          .replace("™", "")
          .replace(":", "")
          .includes(searchText) ? (
          <GridMaker
            name={item.Name}
            cover={item.CoverArtPath}
            uid={item.UID}
            platform={item.OwnedPlatform}
            tileSize={props.tileSize}
          />
        ) : null
      )}
    </div>
  );
}

export default LibraryView;
