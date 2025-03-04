import React, { useRef, useState, useEffect } from "react";
import useIntersectionObserver from "@/hooks/useIntersectionObserver";
import GridView from "@/components/LibraryView/GridView";
import ListView from "@/components/LibraryView/ListView";
import ViewHeader from "@/components/LibraryView/ViewHeader";

interface wishlistViewProps {
  data: any[];
}

export default function WishlistView({ data }: wishlistViewProps) {
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
    const savedGridScrollPos = sessionStorage.getItem(
      "wishlistGridScrollPosition"
    );
    const savedListScrollPos = sessionStorage.getItem(
      "wishlistListScrollPosition"
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
          "wishlistGridScrollPosition",
          gridScrollPositionRef.current.toString()
        );
      }
      if (view === "list") {
        sessionStorage.setItem(
          "wishlistListScrollPosition",
          listScrollPositionRef.current.toString()
        );
      }
    };
  }, [view]);

  return (
    <div className="absolute flex h-full w-full flex-col justify-center select-none">
      <ViewHeader view={view} setView={setView} text={"Wishlist"} />
      {view === "grid" && (
        <GridView
          data={data}
          scrollHandler={scrollHandler}
          gridScrollRef={gridScrollRef}
          visibleItems={visibleItems}
          hidden={false}
        />
      )}
      {view === "list" && (
        <ListView
          listScrollRef={listScrollRef}
          scrollHandler={scrollHandler}
          data={data}
          hidden={false}
        />
      )}
    </div>
  );
}
