import { useSortContext } from "@/hooks/useSortContex";
import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";

interface GridMakerProps {
  cleanedName: string;
  name: string;
  cover: string;
  uid: string;
  platform: string;
  style: React.CSSProperties;
  hidden: boolean;
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
  cleanedName,
  name,
  cover,
  uid,
  platform,
  style,
  hidden,
}: GridMakerProps) {
  const navigate = useNavigate();
  const [preloadData, setPreloadData] = useState<GameData | null>(null);
  const [imageLoaded, setImageLoaded] = useState(false);

  const tileClickHandler = () => {
    console.log(uid);
    navigate(`/gameview`, {
      state: { data: uid, hidden: hidden, preloadData: preloadData },
    });
  };

  const doPreload = useCallback(async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/GameDetails?uid=${uid}`
      );
      const json = await response.json();
      const { companies, tags, screenshots, m: metadata } = json.metadata;
      setPreloadData({
        companies: companies[uid],
        tags: tags[uid],
        metadata: metadata[uid],
        screenshots: screenshots[uid] || [],
      });
    } catch (error) {
      console.error("Failed to preload game data:", error);
    }
  }, [uid]);
  const { cacheBuster } = useSortContext();
  const imageUrl = `http://localhost:8080/cover-art/${cover}?t=${cacheBuster}`;

  useEffect(() => {
    const img = new Image();
    img.src = imageUrl;
    img.onload = () => setImageLoaded(true);
    img.onerror = () => setImageLoaded(false);
  }, []);

  return (
    <div
      onClick={tileClickHandler}
      onMouseEnter={doPreload}
      className="inline-flex rounded-md transition duration-300 ease-in-out hover:scale-105"
      style={style}
    >
      {imageLoaded ? (
        <div
          className="rounded-lg hover:shadow-xl hover:shadow-border hover:transition-shadow relative"
          style={{
            ...style,
            backgroundImage: `url('${imageUrl}')`,
            backgroundSize: "cover",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
            height: "100%",
            width: "100%",
          }}
        />
      ) : (
        <div
          draggable={false}
          className="flex items-center text-emptyGameTileText justify-center bg-emptyGameTile rounded-lg border border-border w-full p-2 text-center text-sm hover:shadow-xl hover:shadow-border hover:transition-shadow"
        >
          {name}
        </div>
      )}
    </div>
  );
}
