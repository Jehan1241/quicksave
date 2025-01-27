import React, { useRef, useState, useEffect } from "react";
import GridMaker from "./GridMaker";
import { useSortContext } from "@/SortContext";
import {
    ChevronDown,
    ChevronUp,
    Grid2X2,
    Grid2X2Icon,
    LibraryBig,
    ListCheckIcon,
    ListIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import DetialsMaker from "@/DetailsView/DetailsMaker";
import { GrCatalogOption } from "react-icons/gr";
interface GridViewProps {
    data: any[];
}

export default function GridView({ data }: GridViewProps) {
    const scrollRef = useRef<HTMLDivElement | null>(null);
    const [scrollPosition, setScrollPosition] = useState(0);
    const [view, setView] = useState("grid");

    const { searchText, sortType, sortOrder, setSortOrder, setSortType, setSortStateUpdate } =
        useSortContext();

    const handleSortChange = (incomingSort: string) => {
        if (incomingSort != sortType) {
            setSortType(incomingSort);
        } else {
            if (sortOrder == "ASC") {
                setSortOrder("DESC");
            }
            if (sortOrder == "DESC") {
                setSortOrder("ASC");
            }
        }
        setSortStateUpdate(true);
    };

    const scrollHandler = () => {
        if (scrollRef.current) {
            const CurrentScrollPos = scrollRef.current.scrollTop;
            localStorage.setItem("ScrollPosition", CurrentScrollPos.toString()); // Convert number to string
            setScrollPosition(CurrentScrollPos);
        }
    };

    useEffect(() => {
        const savedScrollPos = localStorage.getItem("ScrollPosition");
        if (savedScrollPos !== null && scrollRef.current) {
            const scrollPosition = parseInt(savedScrollPos, 10);
            scrollRef.current.scrollTop = scrollPosition;
            setScrollPosition(scrollPosition);
        }
    }, []);

    return (
        <>
            <div
                className="absolute flex h-full w-full flex-col justify-center"
                onScroll={scrollHandler}
                ref={scrollRef}
            >
                <div className="mx-5 flex items-center justify-between p-2 text-xl font-bold tracking-wide">
                    <div className="flex items-center gap-2">Library</div>
                    <div className="flex gap-2">
                        <Button
                            className="h-8 w-8"
                            onClick={() => setView("grid")}
                            variant={"ghost"}
                        >
                            <Grid2X2 strokeWidth={1.7} size={20} />
                        </Button>
                        <Button
                            className="h-8 w-8"
                            onClick={() => setView("list")}
                            variant={"ghost"}
                        >
                            <ListIcon size={20} />
                        </Button>
                    </div>
                </div>
                {view === "grid" && (
                    <div className="flex h-full w-full select-none flex-wrap justify-center gap-8 overflow-y-auto pb-10 pt-4">
                        {data.map((item, key) => {
                            const cleanedName = item.Name.toLowerCase()
                                .replace("'", "")
                                .replace("’", "")
                                .replace("®", "")
                                .replace("™", "")
                                .replace(":", "");
                            if (
                                cleanedName.includes(
                                    searchText.replace("'", "").toLocaleLowerCase()
                                )
                            ) {
                                return (
                                    <GridMaker
                                        cleanedName={cleanedName}
                                        key={item.UID}
                                        name={item.Name}
                                        cover={item.CoverArtPath}
                                        uid={item.UID}
                                        platform={item.OwnedPlatform}
                                    />
                                );
                            }
                        })}
                    </div>
                )}
                {view === "list" && (
                    <div className="mb-5 h-full w-full select-none overflow-y-auto pt-4">
                        <div className="mx-10 flex h-10 justify-between gap-4 rounded-sm bg-background px-5">
                            <div className="flex w-1/4 items-center justify-center">
                                <Button
                                    onClick={() => {
                                        handleSortChange("CustomTitle");
                                    }}
                                    variant={"ghost"}
                                    className="h-8 w-full"
                                >
                                    Title
                                    {sortType == "CustomTitle" && sortOrder == "DESC" && (
                                        <ChevronDown className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                    {sortType == "CustomTitle" && sortOrder == "ASC" && (
                                        <ChevronUp className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                </Button>
                            </div>
                            <div className="flex w-60 items-center justify-center bg-transparent text-center">
                                <Button
                                    onClick={() => {
                                        handleSortChange("OwnedPlatform");
                                    }}
                                    variant={"ghost"}
                                    className="h-8 w-full"
                                >
                                    Platform
                                    {sortType == "OwnedPlatform" && sortOrder == "DESC" && (
                                        <ChevronDown size={22} strokeWidth={0.9} />
                                    )}
                                    {sortType == "OwnedPlatform" && sortOrder == "ASC" && (
                                        <ChevronUp size={22} strokeWidth={0.9} />
                                    )}
                                </Button>
                            </div>
                            <div className="flex w-60 items-center justify-center text-center">
                                <Button
                                    variant={"ghost"}
                                    className="h-8 w-full"
                                    onClick={() => {
                                        handleSortChange("CustomRating");
                                    }}
                                >
                                    Rating
                                    {sortType == "CustomRating" && sortOrder == "DESC" && (
                                        <ChevronDown className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                    {sortType == "CustomRating" && sortOrder == "ASC" && (
                                        <ChevronUp className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                </Button>
                            </div>
                            <div className="flex w-60 cursor-pointer items-center justify-center text-center">
                                <Button
                                    onClick={() => {
                                        handleSortChange("CustomTimePlayed");
                                    }}
                                    variant={"ghost"}
                                    className="h-8 w-full"
                                >
                                    Hours Played
                                    {sortType == "CustomTimePlayed" && sortOrder == "DESC" && (
                                        <ChevronDown className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                    {sortType == "CustomTimePlayed" && sortOrder == "ASC" && (
                                        <ChevronUp className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                </Button>
                            </div>
                            <div className="flex w-60 items-center justify-center text-center">
                                <Button
                                    onClick={() => {
                                        handleSortChange("CustomReleaseDate");
                                    }}
                                    variant={"ghost"}
                                    className="h-8 w-full"
                                >
                                    Release Date
                                    {sortType == "CustomReleaseDate" && sortOrder == "DESC" && (
                                        <ChevronDown className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                    {sortType == "CustomReleaseDate" && sortOrder == "ASC" && (
                                        <ChevronUp className="mx-2" size={22} strokeWidth={0.9} />
                                    )}
                                </Button>
                            </div>
                        </div>
                        {data.map((item, key) => {
                            const cleanedName = item.Name.toLowerCase()
                                .replace("'", "")
                                .replace("’", "")
                                .replace("®", "")
                                .replace("™", "")
                                .replace(":", "");
                            if (
                                cleanedName.includes(
                                    searchText.replace("'", "").toLocaleLowerCase()
                                )
                            ) {
                                return (
                                    <DetialsMaker
                                        cleanedName={cleanedName}
                                        key={item.UID}
                                        name={item.Name}
                                        cover={item.CoverArtPath}
                                        uid={item.UID}
                                        platform={item.OwnedPlatform}
                                        timePlayed={item.TimePlayed}
                                        rating={item.AggregatedRating}
                                        releaseDate={item.ReleaseDate}
                                    />
                                );
                            }
                        })}
                    </div>
                )}
            </div>
        </>
    );
}
