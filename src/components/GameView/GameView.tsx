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
import {
  getGameDetails,
  launchGame,
  sendSteamInstallReq,
} from "@/lib/api/GameViewAPI";
import { useSortContext } from "@/hooks/useSortContex";
import { time } from "node:console";
import { Download, Loader2 } from "lucide-react";
import { importSteamLibrary } from "@/lib/api/libraryImports";
import { platform } from "node:os";
import { useExePathContext } from "@/hooks/useExePathContext";

export default React.memo(GameView);

function enhancedSort(a: string, b: string): number {
  // Extract the part after the last slash
  const getSortablePart = (path: string) => {
    const parts = path.split("/");
    return parts[parts.length - 1].toLowerCase(); // Make comparison case-insensitive
  };

  const partA = getSortablePart(a);
  const partB = getSortablePart(b);

  const isUserA = partA.startsWith("user-");
  const isUserB = partB.startsWith("user-");
  const isGenericA = partA.startsWith("generic-");
  const isGenericB = partB.startsWith("generic-");

  // Both are user files
  if (isUserA && isUserB) {
    const numA = parseInt(partA.match(/user-(\d+)/i)?.[1] || "0");
    const numB = parseInt(partB.match(/user-(\d+)/i)?.[1] || "0");
    return numB - numA; // Higher numbers first
  }

  // Both are generic files
  if (isGenericA && isGenericB) {
    const numA = parseInt(partA.match(/generic-(\d+)/i)?.[1] || "0");
    const numB = parseInt(partB.match(/generic-(\d+)/i)?.[1] || "0");
    return numA - numB; // Lower numbers first for generic
  }

  // One is user, one is generic
  if (isUserA && isGenericB) return -1; // user comes first
  if (isGenericA && isUserB) return 1; // user comes first

  // One is user/generic, the other is neither (shouldn't happen in your case)
  if (isUserA || isGenericA) return -1;
  if (isUserB || isGenericB) return 1;

  // Default alphabetical sort (as fallback)
  return partA.localeCompare(partB);
}

function GameView() {
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
  const [loading, setLoading] = useState(true);

  const updateDetails = () => {
    getGameDetails(uid, setCompanies, setTags, setMetadata, setScreenshots);
  };
  const { playingGame, setPlayingGame } = useSortContext();

  const playGame = async () => {
    if (installed) {
      launchGame(
        uid,
        setCompanies,
        setTags,
        setMetadata,
        setScreenshots,
        setPlayingGame
      );
    } else if (metadata?.OwnedPlatform === "Steam") {
      sendSteamInstallReq(uid);
    } else {
      setEditDialogOpen(true);
      setEditDialogSelectedTab("path");
    }
  };

  // Its on UID change to accomodate randomGamesClicked
  useEffect(() => {
    setLoading(true);
    console.log("Preload Data", preloadData);
    if (preloadData !== null && preloadData !== undefined) {
      setCompanies(preloadData.companies ?? { 0: "Unknown" });
      setTags(preloadData.tags ?? { 0: "Unknown" });
      setScreenshots(preloadData.screenshots ?? {});

      setMetadata({
        TimePlayed: preloadData.metadata?.TimePlayed ?? 0,
        isDLC: preloadData.metadata?.isDLC ?? false,
        AggregatedRating: preloadData.metadata?.AggregatedRating ?? 0,
        Description:
          preloadData.metadata?.Description ?? "No description available.",
        Name: preloadData.metadata?.Name ?? "Unknown",
        OwnedPlatform: preloadData.metadata?.OwnedPlatform ?? "Unknown",
        ReleaseDate: preloadData.metadata?.ReleaseDate ?? "?-?-?",
        CoverArtPath: preloadData.metadata?.CoverArtPath ?? "",
        InstallPath: preloadData.metadata?.InstallPath ?? "",
      });
      console.log(preloadData);
      setLoading(false);
    } else {
      console.log("else");
      getGameDetails(
        uid,
        (data) => setCompanies(data ?? { 0: "Unknown" }),
        (data) => setTags(data ?? { 0: "Unknown" }),
        (data) =>
          setMetadata(
            data ?? {
              TimePlayed: 0,
              isDLC: false,
              AggregatedRating: 0,
              Description: "No description available.",
              Name: "Unknown",
              OwnedPlatform: "Unknown",
              ReleaseDate: "?-?-?",
              CoverArtPath: "",
            }
          ),
        (data) => setScreenshots(data ?? {})
      ).finally(() => setLoading(false));
    }
  }, [uid]);

  const installed = metadata?.InstallPath === "" ? false : true;
  const tagsArray = Object.values(tags);
  const companiesArray = Object.values(companies);
  let screenshotsArray = Object.values(screenshots);
  screenshotsArray = screenshotsArray.sort(enhancedSort);
  const { cacheBuster } = useSortContext();
  const { exePath } = useExePathContext();
  screenshotsArray = screenshotsArray.map((screenshot) => {
    if (import.meta.env.MODE === "production") {
      return `${exePath}/backend/screenshots/${screenshot}?t=${cacheBuster}`;
    } else {
      return `./backend/screenshots/${screenshot}?t=${cacheBuster}`;
    }
  });
  console.log("ccc", screenshotsArray);
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
      {loading ? (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="animate-spin w-10 h-10" />
        </div>
      ) : (
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
                      disabled={hidden || playingGame || isWishlist}
                      className="h-10 lg:w-20 xl:w-40 2xl:w-48 bg-playButton hover:bg-playButtonHover text-playButtonText text-md font-semibold"
                    >
                      {playingGame ? (
                        "Launched"
                      ) : (
                        <>
                          {installed ? (
                            <>
                              <FaPlay /> Play
                            </>
                          ) : (
                            <>
                              <Download /> Install
                            </>
                          )}
                        </>
                      )}
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
                      title={metadata?.Name}
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
              <CarouselSection uid={uid} screenshotsArray={screenshotsArray} />
            </div>
          </div>
        </>
      )}
    </>
  );
}
