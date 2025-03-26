import React, { useEffect, useLayoutEffect, useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import MultipleSelector, { Option } from "@/components/ui/multiple-selector";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { fetchTagsDevsPlatforms } from "@/lib/api/addGameManuallyAPI";
import {
  clearAllFilters,
  deleteCurrentlyFiltered,
  handleFilterChange,
  hideCurrentlyFiltered,
  loadFilterState,
} from "@/lib/api/filterGamesAPI";
import { useToast } from "@/hooks/use-toast";
import { useSortContext } from "@/hooks/useSortContex";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { useLocation } from "react-router-dom";

export default function FilterGames({
  filterDialogOpen,
  setFilterDialogOpen,
}: any) {
  const { setFilterActive, setDeleteFilterGames, setHideFilterGames } =
    useSortContext();
  const [tagOptions, setTagOptions] = useState([]);
  const [platformOptions, setPlatformOptions] = useState([]);
  const [devOptions, setDevOptions] = useState([]);
  const [selectedPlatforms, setSelectedPlatforms] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedTags, setSelectedTags] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedDevs, setSelectedDevs] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedName, setSelectedName] = useState<
    { value: string; label: string }[]
  >([]);
  const [isLoaded, setIsLoaded] = useState<boolean>(false); // Flag to track if data has been loaded
  const { toast } = useToast();

  const OPTIONS: Option[] = [
    { label: "A", value: "A" },
    { label: "B", value: "B" },
    { label: "C", value: "C" },
    { label: "D", value: "D" },
    { label: "E", value: "E" },
    { label: "F", value: "F" },
    { label: "G", value: "G" },
    { label: "H", value: "H" },
    { label: "I", value: "I" },
    { label: "J", value: "J" },
    { label: "K", value: "K" },
    { label: "L", value: "L" },
    { label: "M", value: "M" },
    { label: "N", value: "N" },
    { label: "O", value: "O" },
    { label: "P", value: "P" },
    { label: "Q", value: "Q" },
    { label: "R", value: "R" },
    { label: "S", value: "S" },
    { label: "T", value: "T" },
    { label: "U", value: "U" },
    { label: "V", value: "V" },
    { label: "W", value: "W" },
    { label: "X", value: "X" },
    { label: "Y", value: "Y" },
    { label: "Z", value: "Z" },
  ];

  const clearFilters = () => {
    clearAllFilters(
      setSelectedDevs,
      setSelectedName,
      setSelectedPlatforms,
      setSelectedTags
    );
    setFilterActive(false);
  };

  const loadFilterOptionsAndState = async () => {
    fetchTagsDevsPlatforms(setTagOptions, setDevOptions, setPlatformOptions);
    loadFilterState(
      setIsLoaded,
      setSelectedDevs,
      setSelectedPlatforms,
      setSelectedTags,
      setSelectedName
    );
  };

  useEffect(() => {
    if (!filterDialogOpen) return;
    loadFilterOptionsAndState();
  }, [filterDialogOpen]); //so that filters have latest state always

  useEffect(() => {
    if (isLoaded) {
      handleFilterChange(
        selectedPlatforms,
        selectedTags,
        selectedName,
        selectedDevs,
        toast
      );
    }
    const filtersActive =
      selectedPlatforms.length > 0 ||
      selectedTags.length > 0 ||
      selectedName.length > 0 ||
      selectedDevs.length > 0;

    setFilterActive(filtersActive);
  }, [selectedTags, selectedPlatforms, selectedDevs, selectedName]);

  return (
    <Sheet open={filterDialogOpen} onOpenChange={setFilterDialogOpen}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Filter</SheetTitle>
        </SheetHeader>
        <div className="my-4 flex flex-col gap-4">
          <div>
            <Button
              className="w-full"
              variant={"outline"}
              onClick={clearFilters}
            >
              Clear All Filters
            </Button>
          </div>
          <div className="flex flex-row items-center gap-4">
            <Label className="w-32 text-center">Platform</Label>

            <MultipleSelector
              options={platformOptions}
              value={selectedPlatforms}
              onChange={(e: any) => {
                setSelectedPlatforms(e);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Platforms"
              emptyIndicator={
                <p className="text-center text-sm">no results found.</p>
              }
            />
          </div>
          <div className="flex flex-row items-center gap-4">
            <Label className="w-32 text-center">Name</Label>

            <MultipleSelector
              defaultOptions={OPTIONS}
              maxSelected={1}
              value={selectedName}
              onChange={(e: any) => {
                setSelectedName(e);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Name"
              emptyIndicator={
                <p className="text-center text-sm">no results found.</p>
              }
            />
          </div>
          <div className="flex flex-row items-center gap-4">
            <Label className="w-32 text-center">Tags</Label>

            <MultipleSelector
              options={tagOptions}
              value={selectedTags}
              onChange={(e: any) => {
                setSelectedTags(e);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Tags"
              emptyIndicator={
                <p className="text-center text-sm">no results found.</p>
              }
            />
          </div>
          <div className="flex flex-row items-center gap-4">
            <Label className="w-32 text-center">Developer</Label>

            <MultipleSelector
              options={devOptions}
              value={selectedDevs}
              onChange={(e: any) => {
                setSelectedDevs(e);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Platforms"
              emptyIndicator={
                <p className="text-center text-sm">no results found.</p>
              }
            />
          </div>
        </div>
        <SheetFooter>
          <Dialog>
            <DialogTrigger asChild>
              <Button className="w-full" variant={"destructive"}>
                Delete Visible Games
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>Delete Games</DialogTitle>
                <DialogDescription>
                  This will delete all currently visible games. Running a
                  library import will re-import them, if you don't want this
                  consider hiding them instead.
                </DialogDescription>
              </DialogHeader>

              <DialogFooter>
                <Button
                  variant={"destructive"}
                  onClick={() => setDeleteFilterGames(true)}
                >
                  Confirm Delete
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
          <Dialog>
            <DialogTrigger asChild>
              <Button className="w-full" variant={"default"}>
                Hide Visible Games
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>Hide Games</DialogTitle>
                <DialogDescription>
                  This will hide all currently visible games. Running a library
                  import will not unhide them.
                </DialogDescription>
              </DialogHeader>

              <DialogFooter>
                <Button
                  variant={"destructive"}
                  onClick={() => setHideFilterGames(true)}
                >
                  Confirm Hide
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
