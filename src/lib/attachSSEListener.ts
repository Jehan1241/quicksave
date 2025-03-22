export function attachSSEListener(fetchData: () => void) {
  const eventSource = new EventSource(
    "http://localhost:8080/sse-steam-updates"
  );

  eventSource.onmessage = (event) => {
    console.log("SSE message received:", event.data);
    console.log("Current location:", location.pathname);

    // Prevent updates if user is on /gameview
    fetchData();
  };

  eventSource.onerror = (error) => {
    console.error("SSE Error:", error);
  };

  return () => {
    console.log("Closing SSE connection...");
    eventSource.close();
  };
}
