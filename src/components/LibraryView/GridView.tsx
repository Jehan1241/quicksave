import React, { useEffect, useMemo, useState } from "react";
import GridMaker from "./GridMaker";
import { useSortContext } from "@/hooks/useSortContex";

interface GridViewProps {
  data: any;
  scrollHandler: () => void;
  gridScrollRef: React.RefObject<HTMLDivElement>;
  visibleItems: Set<string>;
  hidden: boolean;
}

export default function GridView({
  data,
  scrollHandler,
  gridScrollRef,
  visibleItems,
  hidden,
}: GridViewProps) {
  const { searchText } = useSortContext();
  const { tileSize } = useSortContext();
  const tileSizeInt = Number(tileSize / 30);
  const style = {
    width: `calc(11rem * ${tileSizeInt})`,
    height: `calc(16rem * ${tileSizeInt})`,
  };

  const TileRoundness = localStorage.getItem("tile-roundness") || "lg";

  return (
    <div
      onScroll={scrollHandler}
      ref={gridScrollRef}
      className="px-4 flex h-full w-full select-none flex-wrap justify-center gap-8 overflow-y-auto pb-10 pt-4 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-0 focus-visible:rounded-sm"
    >
      {data.map((item: any, index: any) => {
        const cleanedName = item.Name.toLowerCase()
          .replace("'", "")
          .replace("’", "")
          .replace("®", "")
          .replace("™", "")
          .replace(":", "");
        if (
          cleanedName.includes(searchText.replace("'", "").toLocaleLowerCase())
        ) {
          const itemId = `item-${index}`;
          return (
            <div
              key={item.UID}
              id={itemId}
              className="flex items-center justify-center"
              style={style}
            >
              {visibleItems.has(itemId) && (
                <GridMaker
                  data={item}
                  style={style}
                  hidden={hidden}
                  roundness={TileRoundness}
                />
              )}
            </div>
          );
        }
      })}
    </div>
  );
}
