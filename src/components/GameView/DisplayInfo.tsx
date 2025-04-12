import { useEffect, useRef, useState } from "react";
import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";
import { useSortContext } from "@/hooks/useSortContex";
import { useNavigationContext } from "@/hooks/useNavigationContext";
import { handleFilterChange, loadFilterState } from "@/lib/api/filterGamesAPI";

export function DisplayInfo({ data, tags, companies }: any) {
  const [selectedDevs, setSelectedDevs] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedPlats, setSelectedPlats] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedTags, setSelectedTags] = useState<
    { value: string; label: string }[]
  >([]);
  const [selectedName, setSelectedName] = useState<
    { value: string; label: string }[]
  >([]);
  const { lastLibraryPath } = useNavigationContext();

  const navigate = useNavigate();

  useEffect(() => {
    loadFilterState(
      () => {},
      setSelectedDevs,
      setSelectedPlats,
      setSelectedTags,
      setSelectedName
    );
  }, []);

  const { setFilterActive } = useSortContext();

  const clickHandler = (
    plats = selectedPlats,
    tags = selectedTags,
    devs = selectedDevs,
    name = selectedName
  ) => {
    handleFilterChange(plats, tags, name, devs, setFilterActive);
    navigate(lastLibraryPath);
  };

  return (
    <div className="mt-2 flex h-full flex-col gap-4 overflow-y-auto pr-1 text-left">
      {/* Platforms */}
      <div className="flex flex-row items-center justify-start">
        <p className="flex flex-col items-start gap-2 text-base font-medium">
          Platform
          <Button
            onClick={() => {
              setSelectedPlats(() => {
                const updatedPlats = [
                  { value: data?.OwnedPlatform, label: data?.OwnedPlatform },
                ];
                clickHandler(
                  updatedPlats,
                  selectedTags,
                  selectedDevs,
                  selectedName
                );
                return updatedPlats;
              });
            }}
            className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText"
          >
            {data?.OwnedPlatform}
          </Button>
        </p>
      </div>

      {/* Tags */}
      <div className="flex flex-col gap-2 text-base">
        <p className="text-left text-base">Tags</p>
        <div className="group relative">
          <div
            className={`flex flex-wrap gap-2 rounded-md text-center transition-all duration-75 max-h-[60px] group-hover:max-h-[1000px] overflow-hidden`}
          >
            {tags.map((item: any, index: number) => (
              <Button
                className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText text-xs"
                key={index}
                onClick={() => {
                  setSelectedTags((prev) => {
                    const updatedTags = [...prev, { value: item, label: item }];
                    clickHandler(
                      selectedPlats,
                      updatedTags,
                      selectedDevs,
                      selectedName
                    );
                    return updatedTags;
                  });
                }}
              >
                {item}
              </Button>
            ))}
          </div>
        </div>
      </div>

      {/* Developers and Publishers */}
      <div className="flex flex-col items-start justify-center gap-2 text-base">
        <p>Developers And Publishers</p>
        <div className="flex flex-wrap gap-2 rounded-md text-center">
          {companies.map((item: any, index: number) => (
            <Button
              className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText text-xs"
              key={index}
              onClick={() => {
                setSelectedDevs((prev) => {
                  const updatedDevs = [...prev, { value: item, label: item }];
                  clickHandler(
                    selectedPlats,
                    selectedTags,
                    updatedDevs,
                    selectedName
                  );
                  return updatedDevs;
                });
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
