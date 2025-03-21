import React from "react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useSortContext } from "@/hooks/useSortContex";
import { Button } from "@/components/ui/button";
import {
  ArrowDownWideNarrow,
  ChartNoAxesColumnDecreasing,
  ChartNoAxesColumnIncreasing,
} from "lucide-react";

export default function SortGames() {
  const { sortOrder, setSortOrder, sortType, setSortType, setSortStateUpdate } =
    useSortContext();

  const sortTypeClicked = (type: string) => {
    console.log("Sort Order Type To", type);
    setSortType(type);
    setSortStateUpdate(true);
  };

  const sortOrderClicked = (order: "ASC" | "DESC") => {
    console.log("Sort Order Set To", order);
    setSortOrder(order);
    setSortStateUpdate(true);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant={"outline"}
          className="f my-auto h-8 w-8 bg-topBarButtons hover:bg-topBarButtonsHover"
        >
          <ArrowDownWideNarrow size={18} strokeWidth={1} />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        <DropdownMenuLabel className="select-none">
          Sort Games
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <div className="flex select-none flex-col text-sm">
            <div className="flex flex-col text-sm">
              <Button
                className={`${
                  sortOrder === "ASC"
                    ? "bg-selectedSortItem text-selectedSortItemText hover:bg-selectedSortItemHover hover:text-selectedSortItemText"
                    : "bg-popover"
                } h-8 justify-start border-none rounded-sm pl-2`}
                variant={"outline"}
                onClick={() => sortOrderClicked("ASC")}
              >
                <ChartNoAxesColumnIncreasing size={18} /> Ascending
              </Button>
              <Button
                className={`${
                  sortOrder === "DESC"
                    ? "bg-selectedSortItem text-selectedSortItemText hover:bg-selectedSortItemHover hover:text-selectedSortItemText"
                    : "bg-popover"
                } h-8 justify-start border-none rounded-sm pl-2`}
                variant={"outline"}
                onClick={() => sortOrderClicked("DESC")}
              >
                <ChartNoAxesColumnDecreasing size={18} /> Descending
              </Button>
            </div>
          </div>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem
            onClick={() => sortTypeClicked("CustomTitle")}
            className={`${sortType === "CustomTitle" ? "bg-selectedSortItem text-selectedSortItemText focus:bg-selectedSortItemHover focus:text-selectedSortItemText" : ""}`}
          >
            <span>Title</span>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => sortTypeClicked("CustomTimePlayed")}
            className={`${sortType === "CustomTimePlayed" ? "bg-selectedSortItem text-selectedSortItemText focus:bg-selectedSortItemHover focus:text-selectedSortItemText" : ""}`}
          >
            <span>Hours Played</span>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => sortTypeClicked("CustomRating")}
            className={`${sortType === "CustomRating" ? "bg-selectedSortItem text-selectedSortItemText focus:bg-selectedSortItemHover focus:text-selectedSortItemText" : ""}`}
          >
            <span>Rating</span>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => sortTypeClicked("CustomReleaseDate")}
            className={`${sortType === "CustomReleaseDate" ? "bg-selectedSortItem text-selectedSortItemText focus:bg-selectedSortItemHover focus:text-selectedSortItemText" : ""}`}
          >
            <span>Release Date</span>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => sortTypeClicked("OwnedPlatform")}
            className={`${sortType === "OwnedPlatform" ? "bg-selectedSortItem text-selectedSortItemText focus:bg-selectedSortItemHover focus:text-selectedSortItemText" : ""}`}
          >
            <span>Platform</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
