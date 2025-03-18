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

export default function QuicksaveMenu({ handleViewClick }: any) {
  const {
    setIsAddGameDialogOpen,
    setIsIntegrationsDialogOpen,
    setIsWishlistAddDialogOpen,
  } = useSortContext();

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
            <DropdownMenuItem>
              Settings
              <DropdownMenuShortcut>F4</DropdownMenuShortcut>
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
                  <DropdownMenuItem onClick={() => handleViewClick("hidden")}>
                    Hidden Games
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem>Check For Updates</DropdownMenuItem>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Theme</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuItem onClick={() => darkMode()}>
                    Dark
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => lightMode()}>
                    Light
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => redMode()}>
                    Red
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => darkPurpleMode()}>
                    Dark Purple
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem>
              New Team
              <DropdownMenuShortcut>⌘+T</DropdownMenuShortcut>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuItem>GitHub</DropdownMenuItem>
          <DropdownMenuItem>Support</DropdownMenuItem>
          <DropdownMenuItem disabled>API</DropdownMenuItem>
          <DropdownMenuItem>
            Quit
            <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
