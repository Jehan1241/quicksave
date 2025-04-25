import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { TabsContent } from "@/components/ui/tabs";
import {
  getBackupFreq,
  getMinimizedToTray,
  setBackupFreq,
  setMinimizeToTray,
} from "@/lib/generalSettings";
import { useState } from "react";

export function GeneralTab() {
  const [checked, setChecked] = useState(getMinimizedToTray());
  const [backupTime, setBackupTime] = useState(getBackupFreq());

  const handleBackupTimeChange = (value: string) => {
    setBackupTime(value);
    setBackupFreq(value);
  };

  return (
    <TabsContent value="general" className="text-sm flex flex-col gap-4">
      <div className="flex gap-4 items-center">
        <label className="w-48">Minimize app to system tray</label>
        <div className="w-40 flex justify-center">
          <Checkbox
            checked={checked}
            onCheckedChange={(value: boolean) => {
              setChecked(value);
              setMinimizeToTray(value);
            }}
          />
        </div>
      </div>
      <div className="flex gap-4 items-center">
        <label className="w-48">Backup frequency</label>
        <div className="w-40 flex justify-center">
          <Select
            value={backupTime}
            onValueChange={(value) => handleBackupTimeChange(value)}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="on launch">on launch</SelectItem>
              <SelectItem value="every 2 days">every 2 days</SelectItem>
              <SelectItem value="every week">every week</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
    </TabsContent>
  );
}
