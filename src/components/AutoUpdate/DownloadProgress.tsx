import { useEffect, useState } from "react";

export function DownloadProgress({ progress }: { progress: number }) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setIsVisible(true);
    return () => setIsVisible(false);
  }, []);

  return (
    <div
      className={`
      z-50 fixed bottom-4 right-4 bg-Sidebar rounded-lg p-4 w-64 shadow-xl border-2 border-border
      transition-all duration-300 ${isVisible ? "opacity-100 translate-y-0" : "opacity-0 translate-y-2"}
    `}
    >
      <div className="flex items-center gap-3 mb-3">
        <svg className="animate-spin h-5 w-5 text-primary" viewBox="0 0 24 24">
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
            fill="none"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
        <div className="flex-1">
          <h3 className="text-sm font-medium">Downloading Update</h3>
          <p className="text-xs text-muted">Almost there...</p>
        </div>
        <span className="text-sm font-mono">{Math.round(progress)}%</span>
      </div>

      <div className="w-full bg-gray-700 rounded-full h-2 overflow-hidden">
        <div
          className="bg-primary h-full transition-all duration-300 ease-out"
          style={{ width: `${progress}%` }}
        />
      </div>
    </div>
  );
}
