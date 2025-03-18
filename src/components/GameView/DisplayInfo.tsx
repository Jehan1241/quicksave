import { useEffect, useRef } from "react";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";

export function DisplayInfo({ data, tags, companies }: any) {
  const selectedDevs = useRef<string[]>([]);
  const selectedPlats = useRef<string[]>([]);
  const selectedTags = useRef<string[]>([]);
  const selectedName = useRef<string[]>([]);

  const loadFilterState = async () => {
    try {
      console.log("Sending Load Filters");
      const response = await fetch("http://localhost:8080/LoadFilters");
      const data = await response.json();
      console.log(data);

      if (data.developers) {
        selectedDevs.current = data.developers;
      }
      if (data.platform) {
        selectedPlats.current = data.platform;
      }
      if (data.name) {
        selectedName.current = data.name;
      }
      if (data.tags) {
        selectedTags.current = data.tags;
      }

      console.log("AAAA", selectedTags, selectedPlats);
    } catch (error) {
      console.error("Error fetching filter:", error);
    }
  };

  const navigate = useNavigate();

  useEffect(() => {
    loadFilterState();
  }, []);

  const handleFilterChange = async () => {
    console.log(selectedPlats.current, selectedTags.current);
    const filter = {
      tags: selectedTags.current,
      name: selectedName.current,
      platforms: selectedPlats.current,
      devs: selectedDevs.current,
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

  return (
    <div className="mt-2 flex h-full flex-col gap-4 overflow-y-auto pr-1 text-left">
      <div className="flex flex-row items-center justify-start">
        <p className="flex flex-col items-start gap-2 text-base font-medium">
          Platform
          <Button
            onClick={(e) => {
              selectedPlats.current = [data?.OwnedPlatform];
              handleFilterChange();
              navigate(-1);
            }}
            className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText"
          >
            {data?.OwnedPlatform}
          </Button>
        </p>
      </div>
      <div className="flex flex-col gap-2 text-base">
        <p className="text-left text-base">Tags</p>
        <div className="flex flex-wrap gap-2 rounded-md text-center">
          {tags.map((item: any, index: any) => (
            <Button
              className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText text-xs"
              key={index}
              onClick={(e) => {
                selectedTags.current = [...selectedTags.current, item];
                handleFilterChange();
                navigate(-1);
              }}
            >
              {item}
            </Button>
          ))}
        </div>
      </div>

      <div className="flex flex-col items-start justify-center gap-2 text-base">
        <p>Developers And Publishers</p>
        <div className="flex flex-wrap gap-2 rounded-md text-center">
          {companies.map((item: any, index: any) => (
            <Button
              className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText text-xs"
              draggable={false}
              key={index}
              onClick={(e) => {
                selectedDevs.current = [...selectedDevs.current, item];
                handleFilterChange();
                navigate(-1);
              }}
            >
              {item}
            </Button>
          ))}
        </div>
      </div>

      <div className="flex flex-col items-start justify-start gap-2">
        Description
        <div className="flex h-full flex-col">
          <p
            dangerouslySetInnerHTML={{ __html: data?.Description }}
            className="text-sm"
          ></p>
        </div>
      </div>
    </div>
  );
}
