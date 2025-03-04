import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import React, { useState, type ReactNode } from "react";
import { PiBookLight, PiListHeartLight } from "react-icons/pi";

import { Filter } from "lucide-react";
import { Dices } from "lucide-react";
import { useSortContext } from "@/hooks/useSortContex";
import { Slider } from "@/components/ui/slider";
import { useNavigate, useLocation } from "react-router-dom";
import QuicksaveMenu from "./QuicksaveMenu";
import SortGames from "./SortGames";
import FilterGames from "./FilterGames";
import WindowButtons from "./WindowsButtons";

export default function CustomTitleBar({ children }: { children: ReactNode }) {
  const [filterDialogOpen, setFilterDialogOpen] = useState(false);
  const { setSearchText, setRandomGameClicked } = useSortContext(); // Access context

  const location = useLocation();
  const navigate = useNavigate();
  const page = location.pathname;
  console.log("path", page);

  console.log(location.pathname);
  const handleViewClick = (view: "" | "wishlist" | "hidden") => {
    navigate(`/${view}`, { replace: true });
    console.log(`${view} View Clicked`);
  };

  // gloabal context vars
  const { tileSize, setTileSize, setSortStateUpdate } = useSortContext();

  // This one updates only on UI
  const sizeChangeHandler = (newSize: number[]) => {
    setTileSize(newSize[0]);
  };

  // This one commits to DB and triggers on release of mouse
  const sizeChangeHandlerCommit = (newSize: number[]) => {
    setTileSize(newSize[0]);
    localStorage.setItem("tileSize", String(newSize[0]));
  };

  return (
    <>
      <div className="flex h-screen w-screen flex-row">
        <div className="flex h-full w-14 flex-col bg-Sidebar">
          <div className="m-auto flex h-12 w-14">
            <QuicksaveMenu handleViewClick={handleViewClick} />
          </div>
          <div className="h-full w-14">
            <div className="my-4 flex flex-col items-center justify-start gap-4 align-middle">
              <Button
                variant={"ghost"}
                onClick={() => handleViewClick("")}
                className={`group h-auto hover:bg-transparent ${
                  page === "/"
                    ? "rounded-none border-r border-leftbarIcons"
                    : ""
                }`}
              >
                <PiBookLight
                  className={`group-hover:scale-125 text-leftbarIcons ${
                    page === "/" ? "scale-150 group-hover:scale-150" : ""
                  }`}
                  size={22}
                />
              </Button>
              <Button
                variant={"ghost"}
                onClick={() => handleViewClick("wishlist")}
                className={`group h-auto hover:bg-transparent text-leftbarIcons ${
                  page === "/wishlist"
                    ? "rounded-none border-r border-leftbarIcons"
                    : ""
                }`}
              >
                <PiListHeartLight
                  className={`group-hover:scale-125 ${
                    page === "/wishlist"
                      ? "scale-150 group-hover:scale-150"
                      : ""
                  }`}
                  size={22}
                />
              </Button>
            </div>
          </div>
        </div>
        <div className="flex h-full w-full flex-col bg-Sidebar">
          <div className="bg flex flex-row bg-Sidebar">
            <div className="flex h-10 w-full flex-row justify-between p-1 bg-sidebar">
              <div className="flex w-full flex-row bg-Sidebar">
                <div className="draglayer h-full flex-1 bg-Sidebar"></div>
                <div className="relative flex h-full w-[50rem] max-w-[60vw] flex-row gap-3 bg-sidebar">
                  <Input
                    onChange={(e) => {
                      setSearchText(e.target.value);
                      if (location.pathname == "/gameview") {
                        navigate(-1);
                      }
                    }}
                    className="my-auto h-8 bg-topBarButtons"
                    placeholder="Search"
                  />
                  <SortGames />
                  {filterDialogOpen && (
                    <FilterGames
                      filterDialogOpen={filterDialogOpen}
                      setFilterDialogOpen={setFilterDialogOpen}
                    />
                  )}
                  <Button
                    variant={"outline"}
                    onClick={() => {
                      setFilterDialogOpen(!filterDialogOpen);
                    }}
                    className="my-auto h-8 w-8 bg-topBarButtons hover:bg-topBarButtonsHover"
                  >
                    <Filter size={18} strokeWidth={1} />
                  </Button>
                  <Button
                    onClick={() => setRandomGameClicked(true)}
                    variant={"outline"}
                    className="my-auto h-8 w-8 bg-topBarButtons hover:bg-topBarButtonsHover"
                  >
                    <Dices size={18} strokeWidth={1} />
                  </Button>
                  <Slider
                    className="w-80"
                    value={[tileSize]}
                    onValueChange={sizeChangeHandler}
                    onValueCommit={sizeChangeHandlerCommit}
                    step={5}
                    min={15}
                    max={100}
                  />
                </div>
                <div className="draglayer h-full flex-1 bg-Sidebar"></div>
              </div>

              <WindowButtons />
            </div>
          </div>
          <div
            draggable={false}
            className="relative h-full w-full rounded-tl-xl bg-content"
          >
            {children}
          </div>
        </div>
      </div>
    </>
  );
}
