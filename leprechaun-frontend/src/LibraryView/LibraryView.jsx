import { CiSearch } from "react-icons/ci";
import GridMaker from "./GridMaker";
import GridMaker2 from "./GridMaker";
import { FaFilter } from "react-icons/fa";
import { useState } from "react";

function LibraryView(props) {
  const data = props.data;
  const [searchText, setSearchText] = useState("");
  const searchBarHandler = () => {
    setSearchText(document.getElementById("searchBar").value);
  };

  return (
    <div className="relative w-[calc(100vw-104px)] h-screen overflow-scroll text-center bg-gameView ml-1 rounded-l-3xl p-8">
      <div className="fixed top-0 left-1/2 z-10 h-12 drop-shadow-lg -translate-x-1/2">
        <input
          id="searchBar"
          onChange={searchBarHandler}
          className="m-2 w-72 h-9 text-white rounded-lg border-2 border-gray-700 bg-primary"
        />
        <button className="m-2 w-9 h-9 text-white rounded-lg border-2 border-gray-700 bg-primary">
          <CiSearch className="inline" />
        </button>
        <button className="m-2 w-9 h-9 text-white rounded-lg border-2 border-gray-700 bg-primary">
          <FaFilter className="inline" />
        </button>
      </div>
      {data.map((item) =>
        item.Name.toLowerCase().includes(searchText) ? (
          <GridMaker
            name={item.Name}
            cover={item.CoverArtPath}
            uid={item.UID}
            platform={item.OwnedPlatform}
          />
        ) : null
      )}
    </div>
  );
}

export default LibraryView;
