import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";

export function DownloadProgress({ progress }: { progress: number }) {
  return (
    <div className="fixed top-0 right-0 w-full h-full bg-black/90 flex justify-center items-center">
      <div className="z-50 bg-Sidebar rounded-lg p-4 w-1/2 h-1/6 shadow-xl border-2 border-border transition-all duration-300">
        <div className="flex items-center gap-3 mb-3">
          <Loader2 className="animate-spin" />
          <div className="flex-1">
            <h3 className="text-sm font-medium">Downloading Update</h3>
            <p className="text-xs text-muted-foreground">
              The app will automatically restart once downloaded.
            </p>
          </div>
          <span className="text-sm font-mono">{Math.round(progress)}%</span>
        </div>

        <div className="w-full bg-muted-foreground rounded-full h-2 overflow-hidden">
          <div
            className="bg-primary h-full transition-all duration-300 ease-out"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>
    </div>
  );
}
