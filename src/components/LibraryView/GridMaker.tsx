import { useSortContext } from "@/hooks/useSortContex";
import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { doDataPreload } from "@/lib/api/GameViewAPI";

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
  const [imageSrc, setImageSrc] = useState<string | null>(null);
  const [imageLoadFailed, setImageLoadFailed] = useState<boolean>(false);

  const tileClickHandler = () => {
    console.log(uid);
    navigate(`/gameview`, {
      state: { data: uid, hidden: hidden, preloadData: preloadData },
    });
  };

  const doPreload = () => {
    doDataPreload(uid, setPreloadData);
  };

  const { cacheBuster } = useSortContext();
  const imageUrl = `http://localhost:8080/cover-art/${cover}`;

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

  return (
    <div
      onClick={tileClickHandler}
      onMouseEnter={doPreload}
      className="inline-flex rounded-md transition duration-300 ease-in-out hover:scale-105"
      style={style}
    >
      {!imageLoadFailed ? (
        <div
          className="rounded-lg hover:shadow-xl hover:shadow-border hover:transition-shadow relative"
          style={{
            ...style,
            backgroundImage: `url('${imageSrc}')`,
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
