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

export function SettingsDropdown({
  setEditDialogOpen,
  hidden,
  setHideDialogOpen,
  setDeleteDialogOpen,
  unhideGame,
}: any) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
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
          <DropdownMenuItem onClick={unhideGame}>Unhide Game</DropdownMenuItem>
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
