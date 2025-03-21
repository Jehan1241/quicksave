import React, { useEffect, useState } from "react";
import {
  Sheet,
  SheetContent,
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
  handleFilterChange,
  loadFilterState,
} from "@/lib/api/filterGamesAPI";
import { useToast } from "@/hooks/use-toast";

export default function FilterGames({
  filterDialogOpen,
  setFilterDialogOpen,
}: any) {
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
    loadFilterOptionsAndState();
  }, []);

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
      </SheetContent>
    </Sheet>
  );
}
