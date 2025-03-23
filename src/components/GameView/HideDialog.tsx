import { useNavigate } from "react-router-dom";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { Button } from "../ui/button";
import { hideGame } from "@/lib/api/GameViewAPI";

export function HideDialog({ uid, hideDialogOpen, setHideDialogOpen }: any) {
  const navigate = useNavigate();

  return (
    <Dialog open={hideDialogOpen} onOpenChange={setHideDialogOpen}>
      <DialogContent className="h-[600px] max-h-[300px] max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Hide Game</DialogTitle>
          <DialogDescription>
            Hidden games can be viewed and reverted at any time. Custom metadata
            is saved and these games will not be re-imported on a library
            integration update.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex items-end">
          <Button
            variant={"dialogSaveButton"}
            onClick={() => {
              setHideDialogOpen(false);
              hideGame(uid, navigate);
            }}
            className="h-12 w-32"
          >
            Hide Game
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
