export function attachSSEListener(fetchData: () => void) {
  const eventSource = new EventSource(
    "http://localhost:8080/sse-steam-updates"
  );

  eventSource.onmessage = (event) => {
    console.log("SSE message received:", event.data);
    fetchData();
  };
  eventSource.onerror = (error) => {
    console.error("SSE Error:", error);
  };
  return () => {
    eventSource.close();
  };
}
