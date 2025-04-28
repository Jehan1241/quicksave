import { useSortContext } from "@/hooks/useSortContex";
import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import {
  doDataPreload,
  hardDelete,
  hideGame,
  launchGame,
  unhideGame,
} from "@/lib/api/GameViewAPI";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "../ui/context-menu";
import { Button } from "../ui/button";
import { FaPlay } from "react-icons/fa";
import { Clock, Download, Eye, EyeOff, Trash2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { useExePathContext } from "@/hooks/useExePathContext";

interface GridMakerProps {
  data: any;
  style: React.CSSProperties;
  hidden: boolean;
  roundness: string;
}

interface GameData {
  companies: string[];
  tags: string[];
  metadata: {
    name: string;
    description: string;
    [key: string]: any;
  };
  screenshots: string[];
}

export default function GridMaker({
  data,
  style,
  hidden,
  roundness,
}: GridMakerProps) {
  const {
    UID,
    Name,
    CoverArtPath: cover,
    OwnedPlatform: platform,
    TimePlayed,
    InstallPath,
    isDLC: isWishlist,
  } = data;

  const installed = InstallPath === "" ? false : true;
  let playtime = TimePlayed.toFixed(2);
  if (playtime < 0) playtime = 0.0;
  const navigate = useNavigate();
  const [preloadData, setPreloadData] = useState<GameData | null>(null);
  const [imageSrc, setImageSrc] = useState<string | null>(null);
  const [imageLoadFailed, setImageLoadFailed] = useState<boolean>(false);

  const tileClickHandler = () => {
    console.log(UID);
    navigate(`/gameview`, {
      state: { data: UID, hidden: hidden, preloadData: preloadData },
    });
  };

  const doPreload = () => {
    doDataPreload(UID, setPreloadData);
  };

  const { cacheBuster } = useSortContext();

  const { exePath } = useExePathContext();

  let imageUrl = `${exePath}/backend/coverArt${cover}`;
  console.log("Image Url", imageUrl);
  if (import.meta.env.MODE !== "production") {
    imageUrl = `./backend/coverArt${cover}`;
  }

  //Check if image exists & is loadable
  const checkImageLoadable = async (url: string) => {
    try {
      const response = await fetch(url, { method: "HEAD" });
      if (response.ok) {
        const img = new Image();
        img.src = `${url}?t=${cacheBuster}`;
        img.onload = () => {
          setImageSrc(img.src);
          setImageLoadFailed(false);
        };
        img.onerror = () => {
          console.error("Image load failed:", url);
          setImageLoadFailed(true);
        };
      } else {
        setImageLoadFailed(true);
      }
    } catch (error) {
      console.error("Image check error:", error);
      setImageLoadFailed(true);
    }
  };

  useEffect(() => {
    checkImageLoadable(imageUrl);
  }, [imageUrl, cacheBuster]);

  const { setPlayingGame } = useSortContext();

  const playClickHandler = () => {
    if (installed) {
      launchGame(
        UID,
        () => {},
        () => {},
        () => {},
        () => {},
        setPlayingGame
      );
    } else {
      navigate(`/gameview`, {
        state: { data: UID, hidden: hidden, preloadData: preloadData },
      });
    }
  };

  const hideClickHandler = () => {
    if (!hidden) {
      hideGame(UID, () => {});
    } else {
      unhideGame(UID, () => {});
    }
  };

  const [dialogOpen, setDialogOpen] = useState(false);

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    setDialogOpen(true);
  };

  const confirmDelete = () => {
    setDialogOpen(false);
    hardDelete(UID, () => {});
  };

  return (
    <>
      <ContextMenu modal={true}>
        <ContextMenuTrigger>
          <div
            onClick={tileClickHandler}
            onMouseEnter={doPreload}
            className="inline-flex rounded-md transition duration-300 ease-in-out hover:scale-105"
            style={style}
          >
            {!imageLoadFailed ? (
              <div
                className={`group flex flex-col rounded-${roundness} hover:shadow-xl hover:shadow-border hover:transition-shadow overflow-hidden cursor-pointer`}
                style={{
                  ...style,
                  backgroundImage: `url('${imageSrc}')`,
                  backgroundSize: "cover",
                  backgroundPosition: "center",
                  backgroundRepeat: "no-repeat",
                  height: "100%",
                  width: "100%",
                }}
              >
                <div className="inline-flex mx-1 mt-1 gap-4 justify-between">
                  <div
                    className={`px-3 py-1 bg-emptyGameTile text-emptyGameTileText rounded-xl opacity-0 group-hover:opacity-90 transition-opacity duration-300 text-xs truncate`}
                  >
                    {platform}
                  </div>
                  <div
                    className={`px-3 py-1 bg-emptyGameTile text-emptyGameTileText rounded-xl opacity-0 group-hover:opacity-90 transition-opacity duration-300 text-xs truncate`}
                  >
                    <Clock size={16} className="inline mr-1" />
                    {playtime}
                  </div>
                </div>

                <div className="inline-flex mx-1 mb-1 mt-auto">
                  <div
                    className={`px-3 py-1 bg-emptyGameTile text-emptyGameTileText rounded-xl opacity-0 group-hover:opacity-90 transition-opacity duration-300 text-xs truncate`}
                  >
                    {Name}
                  </div>
                </div>
              </div>
            ) : (
              <div
                draggable={false}
                className={`flex items-center text-emptyGameTileText justify-center bg-emptyGameTile rounded-${roundness} border border-border w-full p-2 text-center text-sm hover:shadow-xl hover:shadow-border hover:transition-shadow`}
              >
                {Name}
              </div>
            )}
          </div>
        </ContextMenuTrigger>
        <ContextMenuContent
          onInteractOutside={(e) => {
            if (dialogOpen) e.preventDefault();
          }}
          className="text-sm min-w-52"
        >
          <div className="flex flex-col p-1 gap-2">
            <ContextMenuItem asChild>
              <Button
                disabled={hidden || isWishlist}
                onClick={playClickHandler}
                className="h-10 bg-playButton hover:!bg-playButtonHover hover:!text-playButtonText text-playButtonText text-md"
              >
                {installed ? (
                  <>
                    <FaPlay /> Play
                  </>
                ) : (
                  <>
                    <Download /> Install
                  </>
                )}
              </Button>
            </ContextMenuItem>
            <ContextMenuItem onClick={hideClickHandler}>
              {hidden ? (
                <>
                  <Eye size={16} className="mr-2" /> Unhide
                </>
              ) : (
                <>
                  <EyeOff size={16} className="mr-2" />
                  Hide
                </>
              )}
            </ContextMenuItem>
            <ContextMenuItem onClick={handleDeleteClick}>
              <Trash2 size={16} className="mr-2" />
              Delete
            </ContextMenuItem>
          </div>
        </ContextMenuContent>
      </ContextMenu>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogTrigger></DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {Name}</DialogTitle>
          </DialogHeader>
          <DialogDescription>
            This will remove the game from your library. Running a library
            inegration will re-import it
          </DialogDescription>
          <DialogFooter>
            <Button onClick={confirmDelete} variant={"destructive"}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
