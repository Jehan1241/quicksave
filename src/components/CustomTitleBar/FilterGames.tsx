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
  const [isLoaded, setIsLoaded] = useState(false); // Flag to track if data has been loaded

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

  const loadFilterOptionsAndState = async () => {
    try {
      // Fetch tags
      const tagsResponse = await fetch("http://localhost:8080/getAllTags");
      const tagsData = await tagsResponse.json();

      // Transform the tags into key-value pairs
      const tagsAsKeyValuePairs = tagsData.tags.map((tag: any) => ({
        value: tag,
        label: tag,
      }));
      setTagOptions(tagsAsKeyValuePairs);

      // Fetch developers
      const devsResponse = await fetch(
        "http://localhost:8080/getAllDevelopers"
      );
      const devsData = await devsResponse.json();
      console.log(devsData);

      // Transform the developers into key-value pairs
      const devsAsKeyValuePairs = devsData.devs.map((dev: any) => ({
        value: dev,
        label: dev,
      }));
      setDevOptions(devsAsKeyValuePairs);

      // Fetch platforms
      const platformsResponse = await fetch(
        "http://localhost:8080/getAllPlatforms"
      );
      const platformsData = await platformsResponse.json();
      console.log(platformsData);

      // Transform the platforms into key-value pairs
      const platsAsKeyValuePairs = platformsData.platforms.map((plat: any) => ({
        value: plat,
        label: plat,
      }));
      setPlatformOptions(platsAsKeyValuePairs);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
    loadFilterState();
  };

  const loadFilterState = async () => {
    setIsLoaded(false);
    /* setSelectedDevs([]);
        setSelectedName([]);
        setSelectedPlatforms([]);
        setSelectedTags([]); */
    try {
      console.log("Sending Load Filters");
      const response = await fetch("http://localhost:8080/LoadFilters");
      const data = await response.json();
      console.log(data);
      if (data.developers) {
        const devsAsKeyValuePairs = data.developers.map((dev: any) => ({
          value: dev,
          label: dev,
        }));
        setSelectedDevs(devsAsKeyValuePairs);
      }
      if (data.platform) {
        const platsAsKeyValuePairs = data.platform.map((plat: any) => ({
          value: plat,
          label: plat,
        }));
        setSelectedPlatforms(platsAsKeyValuePairs);
      }
      if (data.tags) {
        const tagsAsKeyValuePairs = data.tags.map((tag: any) => ({
          value: tag,
          label: tag,
        }));
        setSelectedTags(tagsAsKeyValuePairs);
      }
      if (data.name) {
        const nameAsKeyValuePairs = data.name.map((name: any) => ({
          value: name,
          label: name,
        }));
        setSelectedName(nameAsKeyValuePairs);
      }
    } catch (error) {
      console.error("Error fetching filter:", error);
      setIsLoaded(true);
    }
    setIsLoaded(true);
  };

  const handleFilterChange = async () => {
    const platformValues = selectedPlatforms.map((platform) => platform.value);
    const tagValues = selectedTags.map((tag) => tag.value);
    const nameValues = selectedName.map((name) => name.value);
    const devValues = selectedDevs.map((dev) => dev.value);

    const filter = {
      tags: tagValues,
      name: nameValues,
      platforms: platformValues,
      devs: devValues,
    };

    try {
      // Send the filter as a POST request
      console.log("Sending Set Filter");
      const response = await fetch("http://localhost:8080/setFilter", {
        method: "POST",
        headers: {
          "Content-Type": "application/json", // Make sure to set the Content-Type to application/json
        },
        body: JSON.stringify(filter), // Convert the filter object to JSON
      });

      const data = await response.json();
      console.log(data); // Log the response from the server
    } catch (error) {
      console.error("Error fetching filter:", error);
    }
  };

  const clearAllFilters = async () => {
    setSelectedDevs([]);
    setSelectedName([]);
    setSelectedPlatforms([]);
    setSelectedTags([]);
    try {
      console.log("Sending Clear All Filters");
      // Send the filter as a POST request
      const response = await fetch("http://localhost:8080/clearAllFilters");
      const data = await response.json();
    } catch (error) {
      console.error("Error fetching filter:", error);
    }
  };

  useEffect(() => {
    loadFilterOptionsAndState();
  }, []);

  useEffect(() => {
    if (isLoaded) {
      handleFilterChange();
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
              onClick={clearAllFilters}
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
