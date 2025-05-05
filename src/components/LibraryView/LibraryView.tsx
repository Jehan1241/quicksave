import React, { useRef, useState, useEffect, useCallback } from "react";
import useIntersectionObserver from "@/hooks/useIntersectionObserver";
import GridView from "./GridView";
import ListView from "./ListView";
import ViewHeader from "./ViewHeader";
import { getScreenshotBind } from "@/lib/screeenshots";
interface libraryViewProps {
  data: any[];
  hidden: boolean;
  viewText: string;
}

export function LibraryView({ data, hidden, viewText }: libraryViewProps) {
  const gridScrollPositionRef = useRef(0);
  const listScrollPositionRef = useRef(0);
  const listScrollRef = useRef<HTMLDivElement | null>(null);
  const [view, setView] = useState<string | null>(null);
  const { visibleItems, gridScrollRef } = useIntersectionObserver(data, view);

  const ssBind = getScreenshotBind();

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
    const savedGridScrollPos = sessionStorage.getItem(
      viewText + "GridScrollPosition"
    );
    const savedListScrollPos = sessionStorage.getItem(
      viewText + "ListScrollPosition"
    );

    if (
      view === "grid" &&
      savedGridScrollPos !== null &&
      gridScrollRef.current
    ) {
      const scrollPosition = parseInt(savedGridScrollPos, 10);
      gridScrollRef.current.scrollTop = scrollPosition;
    } else if (
      view === "list" &&
      savedListScrollPos !== null &&
      listScrollRef.current
    ) {
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
          viewText + "GridScrollPosition",
          gridScrollPositionRef.current.toString()
        );
      }
      if (view === "list") {
        sessionStorage.setItem(
          viewText + "ListScrollPosition",
          listScrollPositionRef.current.toString()
        );
      }
    };
  }, [view, viewText]); // View Text is a dependency as that determines data source

  return (
    <>
      {!data ? (
        <div className="absolute flex h-full w-full flex-col justify-center select-none">
          <ViewHeader view={view} setView={setView} text={viewText} />
          {view === "grid" && (
            <GridView
              data={data}
              scrollHandler={scrollHandler}
              gridScrollRef={gridScrollRef}
              visibleItems={visibleItems}
              hidden={hidden}
            />
          )}
          {view === "list" && (
            <ListView
              data={data}
              scrollHandler={scrollHandler}
              listScrollRef={listScrollRef}
              hidden={hidden}
            />
          )}
        </div>
      ) : (
        <div className="m-5 text-muted-foreground select-none">
          <div className="flex flex-col text-sm gap-4">
            <p className="text-lg">Welcome to quicksave!</p>
            <p>
              <strong>
                - The app is currently in beta, please report any bugs you may
                experience in the discord.
              </strong>
            </p>
            <p> - To get started add a game through the quicksave icon menu.</p>
            <p>
              - You can also confiugre library integrations to automatically
              sync.
            </p>
            <p>
              - quicksave also features a screenshot manager, bound to "{" "}
              {ssBind} " change this from the settings menu.
            </p>
            <p>
              - The app will periodically backup data in the root directory
              /quicksaveBackup folder to ensure user data is never lost.
            </p>
          </div>
        </div>
      )}
    </>
  );
}
