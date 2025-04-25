import { useSortContext } from "@/hooks/useSortContex";
import { Loader2 } from "lucide-react";

export default function IntegrationsLoading() {
  const { integrationLoadCount, backingUp } = useSortContext();

  return (
    <>
      {integrationLoadCount > 0 ? (
        <div className="flex items-center gap-3 text-sm">
          <span className="inline">Library Sync</span>
          <Loader2 className="animate-spin" size={20} />
        </div>
      ) : null}
      {backingUp && (
        <div className="flex items-center gap-3 text-sm">
          <span className="inline">Backing Up</span>
          <Loader2 className="animate-spin" size={20} />
        </div>
      )}
    </>
  );
}
