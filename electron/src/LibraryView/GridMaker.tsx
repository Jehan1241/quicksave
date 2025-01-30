import React, { useState } from "react";
import { useSortContext } from "@/SortContext";
import { useNavigate } from "react-router-dom";

interface GridMakerProps {
    cleanedName: string;
    name: string;
    cover: string;
    uid: string;
    platform: string;
    hidden: boolean;
}

export default function GridMaker({
    cleanedName,
    name,
    cover,
    uid,
    platform,
    hidden,
}: GridMakerProps) {
    let { tileSize } = useSortContext(); // Access context
    let tileSizeInt = Number(tileSize);

    const [imageError, setImageError] = useState(false);
    const navigate = useNavigate();
    tileSizeInt = tileSizeInt / 30; // Calculate tile size

    const tileClickHandler = () => {
        console.log(uid);
        navigate(`gameview`, {
            state: { data: uid, hidden: hidden },
        });
    };

    const style = {
        width: `calc(11rem * ${tileSizeInt})`,
        height: `calc(16rem * ${tileSizeInt})`,
    };

    return (
        <div
            onClick={tileClickHandler}
            className="inline-flex rounded-md transition duration-300 ease-in-out hover:scale-105"
            style={style}
        >
            {!imageError ? (
                <img
                    className="rounded-lg hover:shadow-xl hover:shadow-border hover:transition-shadow"
                    src={`http://localhost:8080/cover-art/${cover}`}
                    onError={() => setImageError(true)}
                    draggable="false"
                />
            ) : (
                <div
                    className="flex items-center justify-center rounded-lg border border-border bg-background p-2 text-center text-sm hover:shadow-xl hover:shadow-border hover:transition-shadow"
                    style={style}
                >
                    {name}
                </div>
            )}
        </div>
    );
}
