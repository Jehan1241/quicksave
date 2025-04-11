import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { TabsContent } from "@/components/ui/tabs";
import { getTileRoundness, setTileRoundness } from "@/lib/TileRoundness";
import { updateTheme } from "@/lib/ToggleTheme";

export default function ThemeTab() {
  const currentTheme = localStorage.getItem("theme");
  const currentTileRoundness = getTileRoundness();

  const handleThemeChange = (value: string) => {
    updateTheme(value);
  };

  return (
    <TabsContent value="theme" className="flex flex-col gap-4 text-sm">
      <div className="flex gap-4 items-center">
        <label className="w-36">Application Theme</label>
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
      </div>
      <div className="flex gap-4 items-center">
        <label className="w-36">Tile Roundness</label>
        <Select onValueChange={(value) => setTileRoundness(value)}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder={`${currentTileRoundness}`} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">None</SelectItem>
            <SelectItem value="sm">Small</SelectItem>
            <SelectItem value="md">Medium</SelectItem>
            <SelectItem value="lg">Large</SelectItem>
            <SelectItem value="xl">Extra Large</SelectItem>
            <SelectItem value="2xl">2xl</SelectItem>
            <SelectItem value="3xl">3xl</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </TabsContent>
  );
}
