import { useSortContext } from "@/hooks/useSortContex";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../../ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../ui/tabs";
import ThemeTab from "./ThemeTab";
import { Button } from "@/components/ui/button";
import ScreenshotsTab from "./ScreenshotsTab";

export default function Settings() {
  const { settingsDialogOpen, setSettingsDialogOpen } = useSortContext();

  return (
    <Dialog open={settingsDialogOpen} onOpenChange={setSettingsDialogOpen}>
      <DialogContent className="w-1/2 max-w-none h-1/2 flex flex-col">
        <DialogHeader>
          <DialogTitle>Settings</DialogTitle>
        </DialogHeader>
        <Tabs defaultValue="theme">
          <TabsList>
            <TabsTrigger value="theme">Theme</TabsTrigger>
            <TabsTrigger value="screenshots">Screenshots</TabsTrigger>
          </TabsList>
          <ThemeTab />
          <ScreenshotsTab />
        </Tabs>
        <DialogFooter className="mt-auto">
          <DialogTrigger>
            <Button variant={"dialogSaveButton"}>Save</Button>
          </DialogTrigger>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
