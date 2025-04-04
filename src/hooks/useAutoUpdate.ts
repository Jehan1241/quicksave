import { useEffect, useState } from "react";

export function useAutoUpdate() {
  const [update, setUpdate] = useState<{
    status: "idle" | "available" | "downloading";
    version?: string;
    progress?: number;
  }>({ status: "idle" });

  useEffect(() => {
    const availableHandler = (update: { version: string; zipUrl: string }) => {
      setUpdate({ status: "available", version: update.version });
    };

    const progressHandler = (progress: number) => {
      console.log(progress);
      setUpdate((prev) => ({ ...prev, status: "downloading", progress }));
    };

    // Attach event listeners
    window.electron.onUpdateAvailable(availableHandler);
    window.electron.onProgress(progressHandler);

    // Cleanup listeners on unmount
    return () => {
      window.electron.onUpdateAvailable(() => {});
      window.electron.onProgress(() => {});
    };
  }, []);

  const startDownload = (toStartDownload: boolean) => {
    window.electron.sendUpdateResponse(toStartDownload);
  };

  return { ...update, startDownload };
}
