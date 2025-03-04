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

export function HideDialog({ uid, hideDialogOpen, setHideDialogOpen }: any) {
  const navigate = useNavigate();

  const hide = async () => {
    console.log("Sending Hide Game");
    try {
      const response = await fetch(`http://localhost:8080/HideGame?uid=${uid}`);
      const json = await response.json();
    } catch (error) {
      console.error(error);
    }
    navigate("/", { replace: true });
  };

  const hardDelete = async () => {
    console.log("Sending Delete Game");
    try {
      const response = await fetch(
        `http://localhost:8080/DeleteGame?uid=${uid}`
      );
      const json = await response.json();
    } catch (error) {
      console.error(error);
    }
    navigate("/", { replace: true });
  };

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
          <Button variant={"secondary"} onClick={hide} className="h-12 w-32">
            Hide Game
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
