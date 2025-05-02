import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

interface DetailsMakerProps {
  cleanedName: string;
  name: string;
  cover: string;
  uid: string;
  platform: string;
  timePlayed: number;
  rating: number;
  releaseDate: string;
  hidden: boolean;
}
interface GameData {
  companies: string[];
  tags: string[];
  metadata: {
    name: string;
    description: string;
    [key: string]: any; // You can define the structure more precisely if needed
  };
  screenshots: string[];
}

export default function DetialsMaker({
  cleanedName,
  name,
  cover,
  uid,
  platform,
  timePlayed,
  rating,
  releaseDate,
  hidden,
}: DetailsMakerProps) {
  const ratingDecimal = rating < 0 ? "0.00" : rating.toFixed(2);
  const timePlayedDecimal = timePlayed < 0 ? "0.00" : timePlayed.toFixed(2);
  const [year, month, day] = releaseDate.split("-");
  const releaseDateNormalized = `${day}-${month}-${year}`;

  const navigate = useNavigate();

  const tileClickHandler = () => {
    console.log(uid);
    navigate(`/gameview`, {
      state: { data: uid, hidden: hidden, preloadData: preloadData },
    });
  };

  const [preloadData, setPreloadData] = useState<GameData | null>();

  const doPreload = async () => {
    console.log(uid);
    try {
      console.log("Sending Get Game Details");
      const response = await fetch(
        `http://localhost:50001/GameDetails?uid=${uid}`
      );
      const json = await response.json();
      console.log(json);
      const { companies, tags, screenshots, m: metadata } = json.metadata;
      setPreloadData({
        companies: companies[uid],
        tags: tags[uid],
        metadata: metadata[uid],
        screenshots: screenshots[uid] || [], // Ensure it's an array
      });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div
      className="mx-10 my-1 flex h-8 justify-between gap-4 rounded-sm bg-background px-5 text-sm hover:bg-transparent"
      onClick={tileClickHandler}
      onMouseEnter={doPreload}
    >
      <div className="flex w-1/4 items-center overflow-hidden text-ellipsis whitespace-nowrap">
        {name}
      </div>
      <div className="flex w-60 items-center justify-center text-center">
        {platform}
      </div>
      <div className="flex w-60 items-center justify-center text-center">
        {ratingDecimal}
      </div>
      <div className="flex w-60 items-center justify-center text-center">
        {timePlayedDecimal}
      </div>
      <div className="flex w-60 items-center justify-center text-center">
        {releaseDateNormalized}
      </div>
    </div>
  );
}
