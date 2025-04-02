import { useAutoUpdate } from "@/hooks/useAutoUpdate";
import { UpdateAvailableDialog } from "./UpdateAvailableDialog";
import { DownloadProgress } from "./DownloadProgress";
import { useSearchParams } from "react-router-dom";
import { useState } from "react";

export function UpdateManager() {
  const { status, version, progress, startDownload } = useAutoUpdate();
  const [showDialog, setShowDialog] = useState(true);

  if (!showDialog) return;

  return (
    <>
      {status === "available" && showDialog && (
        <UpdateAvailableDialog
          version={version!}
          onCancel={() => {
            setShowDialog(false);
          }}
          onConfirm={startDownload}
        />
      )}

      {status === "downloading" && <DownloadProgress progress={progress!} />}
    </>
  );
}
