import { useEffect, useState } from "react";

export function useAutoUpdate() {
  const [update, setUpdate] = useState<{
    status: "idle" | "available" | "downloading" | "ready";
    version?: string;
    progress?: number;
  }>({ status: "idle" });

  useEffect(() => {
    const availableHandler = (version: string) => {
      setUpdate({ status: "available", version });
    };

    const progressHandler = (progress: number) => {
      setUpdate((prev) => ({ ...prev, status: "downloading", progress }));
    };

    const readyHandler = () => {
      setUpdate((prev) => ({ ...prev, status: "ready", progress: 100 }));
    };

    // Attach event listeners
    window.electron.onUpdateAvailable(availableHandler);
    window.electron.onProgress(progressHandler);
    window.electron.onReady(readyHandler);

    // Cleanup listeners on unmount
    return () => {
      window.electron.onUpdateAvailable(() => {});
      window.electron.onProgress(() => {});
      window.electron.onReady(() => {});
    };
  }, []);

  useEffect(() => {
    if (update.status === "ready") {
      window.electron.restartNow();
    }
  }, [update.status]);

  const startDownload = () => {
    window.electron.startDownload();
    setUpdate((prev) => ({ ...prev, status: "downloading", progress: 0 }));
  };

  return { ...update, startDownload };
}
