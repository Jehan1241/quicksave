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
import Integrations from "../Integrations";
import { IntegrationsTab } from "./IntegrationsTab";

export default function Settings() {
  const { settingsDialogOpen, setSettingsDialogOpen } = useSortContext();

  return (
    <Dialog open={settingsDialogOpen} onOpenChange={setSettingsDialogOpen}>
      <DialogContent className="xl:w-1/2 max-w-none xl:h-1/2 flex flex-col h-2/3 w-2/3">
        <DialogHeader>
          <DialogTitle>Settings</DialogTitle>
        </DialogHeader>
        <Tabs defaultValue="theme">
          <TabsList>
            <TabsTrigger value="theme">Theme</TabsTrigger>
            <TabsTrigger value="integrations">Integrations</TabsTrigger>
            <TabsTrigger value="screenshots">Screenshots</TabsTrigger>
          </TabsList>
          <ThemeTab />
          <ScreenshotsTab />
          <IntegrationsTab />
        </Tabs>
        <DialogFooter className="mt-auto"></DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
