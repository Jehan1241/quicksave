import { Checkbox } from "@/components/ui/checkbox";
import { TabsContent } from "@/components/ui/tabs";
import { getMinimizedToTray, setMinimizeToTray } from "@/lib/generalSettings";
import { useState } from "react";

export function GeneralTab() {
  const [checked, setChecked] = useState(getMinimizedToTray());

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
    </TabsContent>
  );
}
