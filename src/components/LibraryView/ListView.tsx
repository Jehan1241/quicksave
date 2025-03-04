import { Button } from "@/components/ui/button";
import { useSortContext } from "@/hooks/useSortContex";
import { ChevronDown, ChevronUp } from "lucide-react";
import React from "react";
import DetialsMaker from "./DetailsMaker";

export default function ListView({
  listScrollRef,
  scrollHandler,
  data,
  hidden,
}: any) {
  const {
    sortType,
    sortOrder,
    setSortOrder,
    setSortType,
    setSortStateUpdate,
    searchText,
  } = useSortContext();

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

  return (
    <div
      onScroll={scrollHandler}
      ref={listScrollRef}
      className="mb-5 h-full w-full select-none overflow-y-auto pt-4"
    >
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
      {data.map((item: any, key: any) => {
        const cleanedName = item.Name.toLowerCase()
          .replace("'", "")
          .replace("’", "")
          .replace("®", "")
          .replace("™", "")
          .replace(":", "");
        if (
          cleanedName.includes(searchText.replace("'", "").toLocaleLowerCase())
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
              hidden={hidden}
            />
          );
        }
      })}
    </div>
  );
}
