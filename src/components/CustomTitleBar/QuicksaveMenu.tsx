import { Button } from "@/components/ui/button";
import React from "react";
import { BsFloppyFill } from "react-icons/bs";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { darkMode, darkPurpleMode, lightMode, redMode } from "@/ToggleTheme";
import { useSortContext } from "@/hooks/useSortContex";
import { useNavigate } from "react-router-dom";

export default function QuicksaveMenu() {
  const {
    setIsAddGameDialogOpen,
    setIsIntegrationsDialogOpen,
    setIsWishlistAddDialogOpen,
    setSettingsDialogOpen,
  } = useSortContext();

  const navigate = useNavigate();
  const handleViewClick = (view: "" | "wishlist" | "hidden" | "installed") => {
    navigate(`/${view}`, { replace: true });
  };

  return (
    <div className="flex w-16 items-center justify-center">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant={"ghost"} className="group hover:bg-transparent px-2">
            <BsFloppyFill
              size={25}
              className="group-hover:scale-110 text-leftbarIcons"
            />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          <DropdownMenuGroup>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Add a Game</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuItem
                    onClick={() => setIsAddGameDialogOpen(true)}
                  >
                    Add to Library
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setIsWishlistAddDialogOpen(true)}
                  >
                    Add to Wishlist
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem onClick={() => setIsIntegrationsDialogOpen(true)}>
              Integrate Libraries
            </DropdownMenuItem>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>View</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuItem onClick={() => handleViewClick("")}>
                    Library
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleViewClick("wishlist")}>
                    Wishlist
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => handleViewClick("installed")}
                  >
                    Installed
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleViewClick("hidden")}>
                    Hidden
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem onClick={() => setSettingsDialogOpen(true)}>
              Settings
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuItem disabled>Quit</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
