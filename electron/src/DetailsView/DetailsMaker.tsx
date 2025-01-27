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
            className="mx-10 my-1 flex h-8 justify-between gap-4 rounded-sm bg-background px-5 text-sm hover:bg-transparent"
            onClick={tileClickHandler}
        >
            <div className="flex w-1/4 items-center overflow-hidden text-ellipsis whitespace-nowrap">
                {name}
            </div>
            <div className="flex w-60 items-center justify-center text-center">{platform}</div>
            <div className="flex w-60 items-center justify-center text-center">{ratingDecimal}</div>
            <div className="flex w-60 items-center justify-center text-center">
                {timePlayedDecimal}
            </div>
            <div className="flex w-60 items-center justify-center text-center">{releaseDate}</div>
        </div>
    );
}
