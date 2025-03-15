import { useSortContext } from "@/hooks/useSortContex";
import { useState, useEffect, useRef } from "react";

const useIntersectionObserver = (data: any[], view: any) => {
  const [visibleItems, setVisibleItems] = useState<Set<string>>(new Set());
  const observer = useRef<IntersectionObserver | null>(null);
  const gridScrollRef = useRef<HTMLDivElement | null>(null);
  const { searchText } = useSortContext();
  const newVisibleItems = useRef<Set<string>>(new Set());

  // Helper function to initialize the observer
  const initObserver = () => {
    const options = {
      root: gridScrollRef.current, // Observe inside the gridScroll container
      rootMargin: "800px", // A buffer to preload items before they enter view
      threshold: 0, // Trigger when 10% of the element is visible
    };

    // Create the observer to track which items are visible
    observer.current = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        const { id } = entry.target as HTMLElement;
        if (entry.isIntersecting) {
          newVisibleItems.current.add(id);
        } else {
          newVisibleItems.current.delete(id);
        }
      });

      setVisibleItems((prev) => new Set(newVisibleItems.current));
      //setVisibleItems(newVisibleItems.current);
    }, options);
  };

  // useEffect(() => {
  //   // Check if there's a saved state in sessionStorage
  //   const savedVisibleItems = sessionStorage.getItem("visibleItems");
  //   if (savedVisibleItems) {
  //     // If saved state exists, use it to initialize the state
  //     setVisibleItems(new Set(JSON.parse(savedVisibleItems)));
  //   }

  //   // Initialize the observer only once
  //   initObserver();

  //   return () => {
  //     sessionStorage.setItem(
  //       "visibleItems",
  //       JSON.stringify(Array.from(visibleItems))
  //     );
  //     observer.current?.disconnect();
  //   };
  // }, []); // Empty dependency array to only initialize once

  useEffect(() => {
    // Reinitialize the observer when data, searchText, or view changes
    initObserver();

    // Function to observe items and trigger re-evaluation
    const observeItems = () => {
      // Remove the observer from previously observed elements (if any)
      data.forEach((_: any, index: any) => {
        const element = document.getElementById(`item-${index}`);
        if (element && observer.current) {
          observer.current.observe(element);
        }
      });
    };

    // Initialize the observer and observe items when data changes
    observeItems();

    // // Trigger recalculation of visibility after layout changes or page load
    // requestAnimationFrame(() => {
    //   observeItems();
    // });

    // Clean up observer on component unmount
    return () => {
      observer.current?.disconnect();
    };
  }, [data, searchText, view]); // Run when 'data', 'searchText', or 'view' changes

  return { visibleItems, gridScrollRef };
};

export default useIntersectionObserver;
