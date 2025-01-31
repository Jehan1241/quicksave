import React, { useRef, useState, useEffect, useCallback } from "react";
import GridMaker from "./GridMaker";
import { useSortContext } from "@/SortContext";
import { ChevronDown, ChevronUp, Grid2X2, ListIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import DetialsMaker from "@/LibraryView/DetailsMaker";
import useIntersectionObserver from "@/hooks/useIntersectionObserver";
import LibraryHeader from "./ViewHeader";
import GridView from "./GridView";
import ListView from "./ListView";
import ViewHeader from "./ViewHeader";
interface libraryViewProps {
    data: any[];
}

export default function LibraryView({ data }: libraryViewProps) {
    const gridScrollPositionRef = useRef(0);
    const listScrollPositionRef = useRef(0);
    const listScrollRef = useRef<HTMLDivElement | null>(null);
    const [view, setView] = useState<string | null>(null);
    const { visibleItems, gridScrollRef } = useIntersectionObserver(data, view);

    const scrollHandler = () => {
        if (gridScrollRef.current) {
            const currentScrollPos = gridScrollRef.current.scrollTop;
            gridScrollPositionRef.current = currentScrollPos;
        }

        if (listScrollRef.current) {
            const currentScrollPos = listScrollRef.current.scrollTop;
            listScrollPositionRef.current = currentScrollPos;
        }
    };

    useEffect(() => {
        const savedGridScrollPos = sessionStorage.getItem("libraryGridScrollPosition");
        const savedListScrollPos = sessionStorage.getItem("libraryListScrollPosition");

        if (view === "grid" && savedGridScrollPos !== null && gridScrollRef.current) {
            const scrollPosition = parseInt(savedGridScrollPos, 10);
            gridScrollRef.current.scrollTop = scrollPosition;
        } else if (view === "list" && savedListScrollPos !== null && listScrollRef.current) {
            const scrollPosition = parseInt(savedListScrollPos, 10);
            listScrollRef.current.scrollTop = scrollPosition;
        }
        const layout = sessionStorage.getItem("layout");
        if (layout) {
            setView(layout);
        } else {
            setView("grid");
        }

        return () => {
            if (view === "grid") {
                sessionStorage.setItem(
                    "libraryGridScrollPosition",
                    gridScrollPositionRef.current.toString()
                );
            }
            if (view === "list") {
                sessionStorage.setItem(
                    "libraryListScrollPosition",
                    listScrollPositionRef.current.toString()
                );
            }
        };
    }, [view]);

    return (
        <div className="absolute flex h-full w-full flex-col justify-center">
            <ViewHeader view={view} setView={setView} text={"Library"} />
            {view === "grid" && (
                <GridView
                    data={data}
                    scrollHandler={scrollHandler}
                    gridScrollRef={gridScrollRef}
                    visibleItems={visibleItems}
                />
            )}
            {view === "list" && (
                <ListView
                    onScroll={onscroll}
                    listScrollRef={listScrollRef}
                    scrollHandler={scrollHandler}
                    data={data}
                />
            )}
        </div>
    );
}

// export default function LibraryView({ data }: libraryViewProps) {
//     const { tileSize } = useSortContext();
//     const tileSizeInt = Number(tileSize / 30);
//     const style = {
//         width: `calc(11rem * ${tileSizeInt})`,
//         height: `calc(16rem * ${tileSizeInt})`,
//     };
//     const { visibleItems, gridScrollRef } = useIntersectionObserver(data);

//     return (
//         <div className="flex h-full w-full flex-col justify-center">
//             <div
//                 ref={gridScrollRef}
//                 className="absolute flex h-full w-full select-none flex-wrap justify-center gap-8 overflow-y-auto pb-10 pt-4"
//             >
//                 {data.map((item, index) => {
//                     const cleanedName = item.Name.toLowerCase()
//                         .replace("'", "")
//                         .replace("’", "")
//                         .replace("®", "")
//                         .replace("™", "")
//                         .replace(":", "");

//                     const itemId = `item-${index}`;

//                     return (
//                         <div
//                             key={item.UID}
//                             id={itemId}
//                             className="flex items-center justify-center"
//                             style={style}
//                         >
//                             {visibleItems.has(itemId) && (
//                                 <GridMaker
//                                     cleanedName={cleanedName}
//                                     name={item.Name}
//                                     cover={item.CoverArtPath}
//                                     uid={item.UID}
//                                     platform={item.OwnedPlatform}
//                                     hidden={false}
//                                     style={style}
//                                 />
//                             )}
//                         </div>
//                     );
//                 })}
//             </div>
//         </div>
//     );
// }

// export default function LibraryView({ data }: libraryViewProps) {
//     const { tileSize } = useSortContext();
//     const tileSizeInt = Number(tileSize / 30);
//     const tileWidth = 11 * 16 * tileSizeInt; // Width of each tile
//     const tileHeight = 16 * 16 * tileSizeInt; // Height of each tile
//     const gap = 16; // Gap between tiles

//     const [visibleTiles, setVisibleTiles] = useState<number[]>([]);
//     const gridScrollRef = useRef<HTMLDivElement | null>(null);

//     // Define a buffer size (e.g., 2 rows above and below)
//     const bufferRows = 3;

//     // Calculate the number of tiles per row based on the container width
//     const calculateTilesPerRow = useCallback(() => {
//         if (!gridScrollRef.current) return 0;
//         const containerWidth = gridScrollRef.current.clientWidth;
//         const tilesPerRow = Math.floor(containerWidth / (tileWidth + gap));
//         return tilesPerRow;
//     }, [tileWidth, gap]);

//     // Calculate the visible tiles based on the scroll position
//     const calculateVisibleTiles = useCallback(() => {
//         if (!gridScrollRef.current) return;

//         const container = gridScrollRef.current;
//         const { scrollTop, clientHeight } = container;
//         const tilesPerRow = calculateTilesPerRow();

//         if (tilesPerRow === 0) return;

//         // Calculate the start and end rows based on the scroll position
//         const startRow = Math.max(0, Math.floor(scrollTop / tileHeight) - bufferRows);
//         const endRow = Math.ceil((scrollTop + clientHeight) / tileHeight) + bufferRows;
//         // console.log(startRow, endRow);

//         // Calculate the start and end indices of the visible tiles
//         const startIndex = startRow * tilesPerRow;
//         const endIndex = endRow * tilesPerRow + tilesPerRow;

//         setVisibleTiles(Array.from({ length: endIndex - startIndex }, (_, i) => startIndex + i));
//     }, [tileHeight, calculateTilesPerRow, bufferRows]);

//     // Attach a scroll event listener to the container
//     useEffect(() => {
//         const scrollContainer = gridScrollRef.current;
//         if (!scrollContainer) return;

//         calculateVisibleTiles();
//         scrollContainer.addEventListener("scroll", calculateVisibleTiles);

//         return () => {
//             scrollContainer.removeEventListener("scroll", calculateVisibleTiles);
//         };
//     }, [calculateVisibleTiles]);

//     // Handle window resize
//     useEffect(() => {
//         const handleResize = () => {
//             calculateVisibleTiles(); // Recalculate visible tiles when window is resized
//         };

//         window.addEventListener("resize", handleResize);

//         return () => {
//             window.removeEventListener("resize", handleResize);
//         };
//     }, [calculateVisibleTiles]);

//     // Adjust the visible tiles when the data or container size changes
//     useEffect(() => {
//         calculateVisibleTiles();
//     }, [data, calculateVisibleTiles]);

//     // Calculate the total height of the scrollable container
//     const tilesPerRow = calculateTilesPerRow();
//     const totalRows = Math.ceil(data.length / tilesPerRow);
//     const totalHeight = totalRows * tileHeight + (totalRows - 1) * gap;

//     return (
//         <div className="flex h-full w-full flex-col justify-center">
//             <div
//                 ref={gridScrollRef}
//                 className="absolute w-full overflow-y-auto"
//                 style={{ height: "100%" }}
//             >
//                 <div
//                     style={{
//                         height: `${totalHeight}px`,
//                         position: "relative",
//                     }}
//                 >
//                     {visibleTiles.map((index) => {
//                         if (index >= data.length) return null;

//                         const item = data[index];
//                         const cleanedName = item.Name.toLowerCase()
//                             .replace("'", "")
//                             .replace("’", "")
//                             .replace("®", "")
//                             .replace("™", "")
//                             .replace(":", "");

//                         const row = Math.floor(index / tilesPerRow);
//                         const col = index % tilesPerRow;

//                         return (
//                             <div
//                                 key={item.UID}
//                                 style={{
//                                     position: "absolute",
//                                     top: `${row * (tileHeight + gap)}px`,
//                                     left: `${col * (tileWidth + gap)}px`,
//                                     width: `${tileWidth}px`,
//                                     height: `${tileHeight}px`,
//                                 }}
//                             >
//                                 <GridMaker
//                                     cleanedName={cleanedName}
//                                     name={item.Name}
//                                     cover={item.CoverArtPath}
//                                     uid={item.UID}
//                                     platform={item.OwnedPlatform}
//                                     hidden={false}
//                                     style={{ width: tileWidth, height: tileHeight }}
//                                 />
//                             </div>
//                         );
//                     })}
//                 </div>
//             </div>
//         </div>
//     );
// }
