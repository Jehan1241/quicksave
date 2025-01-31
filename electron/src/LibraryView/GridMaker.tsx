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
    style: any;
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

export default function GridMaker({
    cleanedName,
    name,
    cover,
    uid,
    platform,
    hidden,
    style,
}: GridMakerProps) {
    const navigate = useNavigate();
    const [imageError, setImageError] = useState(false);
    const [preloadData, setPreloadData] = useState<GameData | null>(null); // State for game data    const navigate = useNavigate();
    const tileClickHandler = () => {
        console.log(uid);
        navigate(`gameview`, {
            state: { data: uid, hidden: hidden, preloadData: preloadData },
        });
    };

    const doPreload = async () => {
        console.log(uid);
        try {
            console.log("Sending Get Game Details");
            const response = await fetch(`http://localhost:8080/GameDetails?uid=${uid}`);
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
            onClick={tileClickHandler}
            onMouseEnter={doPreload}
            className="inline-flex rounded-md transition duration-300 ease-in-out hover:scale-105"
            style={style}
        >
            {!imageError ? (
                <img
                    className="rounded-lg hover:shadow-xl hover:shadow-border hover:transition-shadow"
                    src={`http://localhost:8080/cover-art/${cover}`}
                    onError={() => setImageError(true)}
                    draggable="false"
                    loading="eager"
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
