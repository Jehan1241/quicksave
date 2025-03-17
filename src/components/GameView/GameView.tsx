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

export default function GameView() {
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
  const navigate = useNavigate();

  const playGame = async () => {
    console.log("Play Game Clicked");
    try {
      const response = await fetch(
        `http://localhost:8080/LaunchGame?uid=${uid}`
      );
      const json = await response.json();
      const launchStatus = json.LaunchStatus;
      if (launchStatus === "ToAddPath") {
        setEditDialogSelectedTab("path");
        setEditDialogOpen(true);
      }
      if (launchStatus === "Launched") {
        fetchData();
      }
    } catch (error) {
      console.log(error);
    }
  };

  // Its on UID change to accomodate randomGamesClicked
  useEffect(() => {
    console.log("Preload Data", preloadData);
    if (preloadData != undefined) {
      setCompanies(preloadData.companies);
      setTags(preloadData.tags);
      setMetadata(preloadData.metadata);
      setScreenshots(preloadData.screenshots);
    } else {
      fetchData();
    }
  }, [uid]);

  const fetchData = async () => {
    console.log("Sending Get Game Details");
    try {
      const response = await fetch(
        `http://localhost:8080/GameDetails?uid=${uid}`
      );
      const json = await response.json();
      console.log(json);
      const { companies, tags, screenshots, m: metadata } = json.metadata;
      setCompanies(companies[uid]);
      setTags(tags[uid]);
      setMetadata(metadata[uid]);
      setScreenshots(screenshots[uid] || []); // Make sure it's an array
    } catch (error) {
      console.error(error);
    }
  };

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

  const unhideGame = async () => {
    try {
      console.log("Sending Get Game Details");
      const response = await fetch(
        `http://localhost:8080/unhideGame?uid=${uid}`
      );
      const json = await response.json();
      console.log(json);
    } catch (error) {
      console.error(error);
    }
    navigate("/hidden", { replace: true });
  };

  return (
    <>
      <img
        className="absolute z-0 h-full w-full rounded-2xl object-cover opacity-20 blur-md"
        src={"http://localhost:8080/screenshots/" + screenshots[0]}
      />

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
                  hidden={hidden}
                  unhideGame={unhideGame}
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
                  fetchData={fetchData}
                  coverArtPath={metadata?.CoverArtPath}
                  screenshotsArray={screenshotsArray}
                  platform={metadata?.OwnedPlatform}
                  tags={tagsArray}
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
