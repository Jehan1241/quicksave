export function attachSSEListener(fetchData: () => void, setCacheBuster: any) {
  const eventSource = new EventSource(
    "http://localhost:8080/sse-steam-updates"
  );

  eventSource.onmessage = (event) => {
    console.log("SSE message received:", event.data);
    fetchData();
    //This initially solved concurrent steam and psn import imgs stuck pending but now not needed?
    //setTimeout(() => setCacheBuster(Date.now()), 1000);
  };
  eventSource.onerror = (error) => {
    console.error("SSE Error:", error);
  };
  return () => {
    eventSource.close();
  };
}
