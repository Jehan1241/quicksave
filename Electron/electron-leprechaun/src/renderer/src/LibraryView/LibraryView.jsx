import { useRef, useState, useEffect } from "react";
import GridMaker from "./GridMaker";

function LibraryView(props) {
  const data = props.data;
  const searchText = props.searchText;
  const scrollRef = useRef();
  const [scrollPosition, setScrollPosition] = useState(0);

  const scrollHandler = () => {
    const CurrentScrollPos = scrollRef.current.scrollTop;
    localStorage.setItem("ScrollPosition", CurrentScrollPos);
    setScrollPosition(CurrentScrollPos);
  };

  useEffect(() => {
    const savedScrollPos = localStorage.getItem("ScrollPosition");
    if (savedScrollPos !== null) {
      const scrollPosition = parseInt(savedScrollPos, 10);
      scrollRef.current.scrollTop = scrollPosition;
      setScrollPosition(scrollPosition);
    }
  }, []);

  return (
    <div
      className="overflow-y-auto min-h-screen text-center bg-gradient-to-br from-neutral-900/50 via-amber-950/70 to-neutral-900/50"
      onScroll={scrollHandler}
      ref={scrollRef}
    >
      <div className="overflow-y-auto min-h-screen text-center bg-gradient-to-r from-purple-900/20 via-neutral-900/50 to-purple-900/20">
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
              key={item.UID}
              name={item.Name}
              cover={item.CoverArtPath}
              uid={item.UID}
              platform={item.OwnedPlatform}
              tileSize={props.tileSize}
              scrollPosition={scrollPosition}
            />
          ) : null
        )}
      </div>
    </div>
  );
}

export default LibraryView;
