import { Settings2 } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";
import { unhideGame } from "@/lib/api/GameViewAPI";

export function SettingsDropdown({
  uid,
  setEditDialogOpen,
  hidden,
  setHideDialogOpen,
  setDeleteDialogOpen,
}: any) {
  const navigate = useNavigate();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="h-10 bg-editButton hover:bg-editButtonHover text-editButtonText">
          <Settings2 />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Edit Menu</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => setEditDialogOpen(true)}>
          Edit Metadata
        </DropdownMenuItem>
        {hidden ? (
          <DropdownMenuItem onClick={() => unhideGame(uid, navigate)}>
            Unhide Game
          </DropdownMenuItem>
        ) : (
          <DropdownMenuItem onClick={() => setHideDialogOpen(true)}>
            Hide Game
          </DropdownMenuItem>
        )}
        <DropdownMenuItem onClick={() => setDeleteDialogOpen(true)}>
          Delete Game
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
