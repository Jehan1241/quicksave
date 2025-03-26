import { useSortContext } from "@/hooks/useSortContex";
import { useState } from "react";
import { Input } from "../ui/input";
import { useLocation, useNavigate } from "react-router-dom";
import SortGames from "./SortGames";
import FilterGames from "./FilterGames";
import { Button } from "../ui/button";
import { Dices, Filter, X } from "lucide-react";
import { Slider } from "../ui/slider";
import IntegrationsLoading from "./IntegrationsLoading";
import WindowButtons from "./WindowsButtons";
import { useNavigationContext } from "@/hooks/useNavigationContext";

export default function TopBar() {
  const {
    tileSize,
    setTileSize,
    searchText,
    setSearchText,
    setRandomGameClicked,
    randomGameClicked,
  } = useSortContext();
  const [filterDialogOpen, setFilterDialogOpen] = useState(false);

  const { lastLibraryPath } = useNavigationContext();

  const navigate = useNavigate();
  const location = useLocation();

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
    <div
      className="flex w-full flex-col bg-Sidebar items-center"
      style={{ WebkitAppRegion: "drag" } as any}
    >
      <div className="bg flex flex-row p-1 items-center w-full">
        {/* Left Spacer */}
        <div className="hiden xl:flex-1"></div>
        {/* Centered Search & Buttons */}
        <div
          className="relative flex h-full w-[50rem] max-w-[60vw] flex-row gap-3 bg-sidebar justify-center"
          style={{ WebkitAppRegion: "no-drag" } as any}
        >
          <div className="relative w-full">
            <Input
              value={searchText}
              onChange={(e) => {
                setSearchText(e.target.value);
                if (location.pathname == "/gameview") {
                  navigate(lastLibraryPath, { replace: true });
                }
              }}
              className="my-auto h-8 bg-topBarButtons"
              placeholder="Search"
            />
            <button
              className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-accent-foreground"
              onClick={() => {
                setSearchText("");
              }}
            >
              <X size={20} />
            </button>
          </div>
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
            onClick={() => {
              if (randomGameClicked) return; // Prevents spamming from breaking random games
              setRandomGameClicked(true);
            }}
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
        {/* Centered IntegrationsLoading in remaining space */}
        <div className="flex-1 flex justify-center">
          <IntegrationsLoading />
        </div>
        <WindowButtons />
      </div>
    </div>
  );
}
