import React, { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { FaPlay } from "react-icons/fa";
import { HideDialog } from "./HideDialog";
import { DeleteDialog } from "./DeleteDialog";
import { DisplayInfo } from "./DisplayInfo";
import { EditDialog } from "./EditDialog/EditDialog";
import { CarouselSection } from "./CarouselSection";
import { DateTimeRatingSection } from "./DateTimeRatingSection";
import { SettingsDropdown } from "./SettingsDropdown";
import { useSortContext } from "@/hooks/useSortContex";
import { fetchData } from "@/lib/api/fetchBasicInfo";
import { getGameDetails, launchGame } from "@/lib/api/GameViewAPI";

export default function GameView() {
  const navigate = useNavigate();
  const location = useLocation();
  const uid = location.state.data;
  const hidden = location.state.hidden;
  const preloadData = location.state.preloadData;

  const [companies, setCompanies] = useState("");
  const [tags, setTags] = useState("");
  const [screenshots, setScreenshots] = useState("");
  const [metadata, setMetadata] = useState<any>();
  const [editDialogOpen, setEditDialogOpen] = useState<boolean>(false);
  const [editDialogSelectedTab, setEditDialogSelectedTab] = useState<
    "metadata" | "images" | "path"
  >("metadata");
  const [hideDialogOpen, setHideDialogOpen] = useState<boolean>(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);

  const updateDetails = () => {
    getGameDetails(uid, setCompanies, setTags, setMetadata, setScreenshots);
  };

  const playGame = () => {
    launchGame(
      uid,
      setEditDialogOpen,
      setEditDialogSelectedTab,
      setCompanies,
      setTags,
      setMetadata,
      setScreenshots
    );
  };

  // Its on UID change to accomodate randomGamesClicked
  useEffect(() => {
    console.log("Preload Data", preloadData);
    if (preloadData === undefined) {
      setCompanies(preloadData.companies ?? { 0: "unknown" });
      setTags(preloadData.tags ?? { 0: "unknown" });
      setScreenshots(preloadData.screenshots ?? {});

      setMetadata({
        AggregatedRating: preloadData.metadata?.AggregatedRating ?? 0,
        Description:
          preloadData.metadata?.Description ?? "No description available.",
        Name: preloadData.metadata?.Name ?? "Unknown",
        OwnedPlatform: preloadData.metadata?.OwnedPlatform ?? "Unknown",
        ReleaseDate: preloadData.metadata?.ReleaseDate ?? "?-?-?",
      });
    } else {
      getGameDetails(uid, setCompanies, setTags, setMetadata, setScreenshots);
    }
  }, [uid]);

  const tagsArray = Object.values(tags);
  const companiesArray = Object.values(companies);
  let screenshotsArray = Object.values(screenshots);
  const { cacheBuster } = useSortContext();
  screenshotsArray = screenshotsArray.map((screenshot) => {
    return `http://localhost:8080/screenshots${screenshot}?t=${cacheBuster}`;
  });
  let timePlayed = metadata?.TimePlayed?.toFixed(1);
  if (timePlayed < 0) timePlayed = "0.0";
  const isWishlist = metadata?.isDLC;
  const rating = metadata?.AggregatedRating?.toFixed(1);
  let releaseDate = metadata?.ReleaseDate;

  if (releaseDate === "1970-01-01") {
    releaseDate = "unknown";
    // Reformat to "dd-mm-yyyy"
  } else if (releaseDate) {
    const [year, month, day] = releaseDate.split("-");
    releaseDate = `${day}-${month}-${year}`;
  }

  return (
    <>
      {/* Avoid undefined URL error */}
      {screenshotsArray.length > 0 && (
        <img
          className="absolute z-0 h-full w-full rounded-2xl object-cover opacity-20 blur-md"
          src={screenshotsArray[0]}
        />
      )}
      <div className="absolute z-10 flex h-full w-full flex-col overflow-y-hidden px-6 py-8 text-center select-none">
        <header className="mx-8 mb-2 text-left text-3xl font-semibold">
          {metadata?.Name}
        </header>
        <div className="mx-8 mb-4 flex h-full flex-row gap-10 overflow-hidden">
          <div className="flex h-full w-1/3 flex-col overflow-y-auto">
            <div className="mt-2 flex w-full flex-row items-center gap-4 text-base font-normal xl:flex-row">
              <div className="flex gap-2">
                <Button
                  onClick={playGame}
                  disabled={isWishlist === 0 ? false : true}
                  className="h-10 lg:w-20 xl:w-40 2xl:w-48 bg-playButton hover:bg-playButtonHover text-playButtonText"
                >
                  <FaPlay /> Play
                </Button>

                <SettingsDropdown
                  uid={uid}
                  hidden={hidden}
                  setEditDialogOpen={setEditDialogOpen}
                  setHideDialogOpen={setHideDialogOpen}
                  setDeleteDialogOpen={setDeleteDialogOpen}
                />

                <EditDialog
                  uid={uid}
                  editDialogSelectedTab={editDialogSelectedTab}
                  setEditDialogSelectedTab={setEditDialogSelectedTab}
                  editDialogOpen={editDialogOpen}
                  setEditDialogOpen={setEditDialogOpen}
                  getGameDetails={updateDetails}
                  coverArtPath={metadata?.CoverArtPath}
                  screenshotsArray={screenshotsArray}
                  platform={metadata?.OwnedPlatform}
                  tags={tagsArray}
                  companies={companiesArray}
                />
                <HideDialog
                  uid={uid}
                  hideDialogOpen={hideDialogOpen}
                  setHideDialogOpen={setHideDialogOpen}
                />
                <DeleteDialog
                  uid={uid}
                  deleteDialogOpen={deleteDialogOpen}
                  setDeleteDialogOpen={setDeleteDialogOpen}
                />
              </div>

              <DateTimeRatingSection
                releaseDate={releaseDate}
                rating={rating}
                isWishlist={isWishlist}
                timePlayed={timePlayed}
              />
            </div>

            <DisplayInfo
              data={metadata}
              tags={tagsArray}
              companies={companiesArray}
            />
          </div>
          <CarouselSection screenshotsArray={screenshotsArray} />
        </div>
      </div>
    </>
  );
}
