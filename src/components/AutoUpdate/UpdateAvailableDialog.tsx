import { Button } from "../ui/button";

// UpdateAvailableDialog.tsx
export function UpdateAvailableDialog({
  version,
  onCancel,
  onConfirm,
}: {
  version: string;
  onCancel: () => void;
  onConfirm: () => void;
}) {
  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div className="bg-Sidebar rounded-lg border border-border max-w-md w-full p-6 animate-fade-in">
        <h2 className="text-xl font-semibold mb-2">New Version Available</h2>
        <p className="text-sm mb-4">Version {version} is ready to download.</p>

        <div className="flex justify-end gap-3">
          <Button variant={"outline"} onClick={onCancel}>
            Later
          </Button>
          <Button variant={"outline"} onClick={onConfirm}>
            Download Update
          </Button>
        </div>
      </div>
    </div>
  );
}
