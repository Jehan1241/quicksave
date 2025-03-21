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
import { hardDelete } from "@/lib/api/GameViewAPI";

export function DeleteDialog({
  uid,
  deleteDialogOpen,
  setDeleteDialogOpen,
}: any) {
  const navigate = useNavigate();

  return (
    <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
      <DialogContent className="h-[600px] max-h-[300px] max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Delete Game</DialogTitle>
          <DialogDescription>
            The game will be permanently deleted from the app. The game will be
            re-imported on library synchronization, if you wish to stop this
            behaviour, hide the game instead.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex items-end">
          <Button
            variant={"destructive"}
            onClick={() => hardDelete(uid, navigate)}
            className="h-12 w-32"
          >
            Delete Game
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
