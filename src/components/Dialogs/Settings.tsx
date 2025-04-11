import { useSortContext } from "@/hooks/useSortContex";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import {
  darkMode,
  darkPurpleMode,
  lightMode,
  redMode,
  setTheme,
  updateTheme,
} from "@/ToggleTheme";

export default function Settings() {
  const { settingsDialogOpen, setSettingsDialogOpen } = useSortContext();

  console.log("SEtt", settingsDialogOpen);

  const currentTheme = localStorage.getItem("theme");
  console.log(currentTheme);

  const handleThemeChange = (value: string) => {
    updateTheme(value);
  };

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
          <TabsContent
            value="theme"
            className="flex items-center gap-4 text-sm"
          >
            Application Theme
            <Select onValueChange={(value) => handleThemeChange(value)}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder={`${currentTheme}`} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="dark">dark</SelectItem>
                <SelectItem value="light">light</SelectItem>
                <SelectItem value="red">red</SelectItem>
                <SelectItem value="magenta-dark">magenta-dark</SelectItem>
              </SelectContent>
            </Select>
          </TabsContent>
          <TabsContent value="screenshots">SS</TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}
