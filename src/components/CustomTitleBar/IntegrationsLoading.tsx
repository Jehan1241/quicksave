import { useSortContext } from "@/hooks/useSortContex";
import { Loader2 } from "lucide-react";

export default function IntegrationsLoading() {
  const { psnLoading, steamLoading } = useSortContext();

  return (
    <>
      {psnLoading || steamLoading ? (
        <div className="flex items-center gap-3 text-sm">
          {/* Show text only on `md` and larger screens */}
          <span className="inline">Library Sync</span>
          {/* Always show the loader */}
          <Loader2 className="animate-spin" size={20} />
        </div>
      ) : null}
    </>
  );
}
