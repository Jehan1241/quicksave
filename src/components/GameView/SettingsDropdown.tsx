import { Eye, EyeOff, Pencil, Settings2, Trash2 } from "lucide-react";
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
      <DropdownMenuContent className="w-44">
        <DropdownMenuItem onClick={() => setEditDialogOpen(true)}>
          <Pencil className="mr-1" /> Edit Metadata
        </DropdownMenuItem>
        {hidden ? (
          <DropdownMenuItem onClick={() => unhideGame(uid, navigate)}>
            <Eye size={16} className="mr-1" />
            Unhide Game
          </DropdownMenuItem>
        ) : (
          <DropdownMenuItem onClick={() => setHideDialogOpen(true)}>
            <EyeOff size={16} className="mr-1" />
            Hide Game
          </DropdownMenuItem>
        )}
        <DropdownMenuItem onClick={() => setDeleteDialogOpen(true)}>
          <Trash2 size={16} className="mr-1" /> Delete Game
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
