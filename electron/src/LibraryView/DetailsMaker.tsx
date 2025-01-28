import React, { useState } from "react";
import { useSortContext } from "@/SortContext";
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
}: DetailsMakerProps) {
    // Only does rating decimal if it has decimals to begin with
    let ratingDecimal;
    ratingDecimal = rating;
    if (rating % 1 !== 0) {
        ratingDecimal = rating.toFixed(1);
    }
    let timePlayedDecimal;
    timePlayedDecimal = timePlayed;
    if (timePlayed % 1 !== 0) {
        timePlayedDecimal = timePlayed.toFixed(1);
    }

    const navigate = useNavigate();

    const tileClickHandler = () => {
        console.log(uid);
        navigate(`gameview`, {
            state: { data: uid },
        });
    };

    return (
        <div
            className="flex gap-4 justify-between px-5 mx-10 my-1 h-8 text-sm rounded-sm bg-background hover:bg-transparent"
            onClick={tileClickHandler}
        >
            <div className="flex overflow-hidden items-center w-1/4 whitespace-nowrap text-ellipsis">
                {name}
            </div>
            <div className="flex justify-center items-center w-60 text-center">{platform}</div>
            <div className="flex justify-center items-center w-60 text-center">{ratingDecimal}</div>
            <div className="flex justify-center items-center w-60 text-center">
                {timePlayedDecimal}
            </div>
            <div className="flex justify-center items-center w-60 text-center">{releaseDate}</div>
        </div>
    );
}
