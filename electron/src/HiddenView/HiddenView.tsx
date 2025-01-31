import { Button } from "@/components/ui/button";
import useIntersectionObserver from "@/hooks/useIntersectionObserver";
import DetialsMaker from "@/LibraryView/DetailsMaker";
import GridMaker from "@/LibraryView/GridMaker";
import GridView from "@/LibraryView/GridView";
import ListView from "@/LibraryView/ListView";
import ViewHeader from "@/LibraryView/ViewHeader";
import { useSortContext } from "@/SortContext";
import { ChevronDown, ChevronUp, Grid2X2, ListIcon } from "lucide-react";
import React, { useEffect, useRef, useState } from "react";

interface hiddenViewProps {
    data: any[];
}

export default function HiddenView({ data }: hiddenViewProps) {
    const gridScrollPositionRef = useRef(0);
    const listScrollPositionRef = useRef(0);
    const listScrollRef = useRef<HTMLDivElement | null>(null);
    const [view, setView] = useState<string | null>(null);
    const { visibleItems, gridScrollRef } = useIntersectionObserver(data, view);

    const scrollHandler = () => {
        if (gridScrollRef.current) {
            const currentScrollPos = gridScrollRef.current.scrollTop;
            gridScrollPositionRef.current = currentScrollPos;
        }

        if (listScrollRef.current) {
            const currentScrollPos = listScrollRef.current.scrollTop;
            listScrollPositionRef.current = currentScrollPos;
        }
    };

    useEffect(() => {
        const savedGridScrollPos = sessionStorage.getItem("hiddenGridScrollPosition");
        const savedListScrollPos = sessionStorage.getItem("hiddenListScrollPosition");

        if (view === "grid" && savedGridScrollPos !== null && gridScrollRef.current) {
            const scrollPosition = parseInt(savedGridScrollPos, 10);
            gridScrollRef.current.scrollTop = scrollPosition;
        } else if (view === "list" && savedListScrollPos !== null && listScrollRef.current) {
            const scrollPosition = parseInt(savedListScrollPos, 10);
            listScrollRef.current.scrollTop = scrollPosition;
        }
        const layout = sessionStorage.getItem("layout");
        if (layout) {
            setView(layout);
        } else {
            setView("grid");
        }

        return () => {
            if (view === "grid") {
                sessionStorage.setItem(
                    "hiddenGridScrollPosition",
                    gridScrollPositionRef.current.toString()
                );
            }
            if (view === "list") {
                sessionStorage.setItem(
                    "hiddenListScrollPosition",
                    listScrollPositionRef.current.toString()
                );
            }
        };
    }, [view]);

    return (
        <div className="absolute flex h-full w-full flex-col justify-center">
            <ViewHeader view={view} setView={setView} text={"Hidden Games"} />
            {view === "grid" && (
                <GridView
                    data={data}
                    scrollHandler={scrollHandler}
                    gridScrollRef={gridScrollRef}
                    visibleItems={visibleItems}
                />
            )}
            {view === "list" && (
                <ListView
                    onScroll={onscroll}
                    listScrollRef={listScrollRef}
                    scrollHandler={scrollHandler}
                    data={data}
                />
            )}
        </div>
    );
}
