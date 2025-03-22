import { useEffect } from "react";
import { useLocation } from "react-router-dom";

export function useSSEListener(fetchData: () => void) {
  const location = useLocation();

  useEffect(() => {
    const eventSource = new EventSource(
      "http://localhost:8080/sse-steam-updates"
    );

    eventSource.onmessage = (event) => {
      console.log("SSE message received:", event.data);
      console.log("Current location:", location.pathname);

      // Prevent updates if user is on /gameview
      if (location.pathname !== "/gameview") {
        fetchData();
      }
    };

    eventSource.onerror = (error) => {
      console.error("SSE Error:", error);
    };

    return () => {
      console.log("Closing SSE connection...");
      eventSource.close();
    };
  }, [location.pathname]); // Depend on location.pathname
}
