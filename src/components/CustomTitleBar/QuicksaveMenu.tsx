import { Button } from "@/components/ui/button";
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
import { useSortContext } from "@/hooks/useSortContex";
import { useNavigate } from "react-router-dom";
import { FolderDown, Plus, Settings } from "lucide-react";
import image from "@/../assets/image.svg";
import {
  DiscordLogoIcon,
  GitHubLogoIcon,
  QuestionMarkCircledIcon,
} from "@radix-ui/react-icons";

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
            <img
              src={image}
              className="group-hover:scale-110 text-leftbarIcons"
              width={32}
              height={32}
            />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          <DropdownMenuGroup>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>
                <Plus /> Add a Game
              </DropdownMenuSubTrigger>
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
              <FolderDown /> Integrate Libraries
            </DropdownMenuItem>
            <DropdownMenuSub>
              <DropdownMenuSubTrigger className="pl-8">
                View
              </DropdownMenuSubTrigger>
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
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>
                <QuestionMarkCircledIcon />
                Help
              </DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuItem>
                    <GitHubLogoIcon />
                    <a href="https://github.com/Jehan1241/quicksave">Github</a>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <DiscordLogoIcon />
                    <a href="https://discord.gg/jYTmYM6YmG">Discord</a>
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>

            <DropdownMenuItem onClick={() => setSettingsDialogOpen(true)}>
              <Settings /> Settings
            </DropdownMenuItem>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
