import { Checkbox } from "@/components/ui/checkbox";
import { TabsContent } from "@/components/ui/tabs";
import {
  getIntegrateOnExitEnabled,
  getIntegrateOnLaunchEnabled,
  setIntegrateOnExitEnabled,
  setIntegrateOnLaunchEnabled,
} from "@/lib/integrationSettings";
import { useState } from "react";

export function IntegrationsTab() {
  const [checkedOnLaunch, setCheckedonLaunch] = useState(
    getIntegrateOnLaunchEnabled()
  );
  const [checkedOnExit, setCheckedonExit] = useState(
    getIntegrateOnExitEnabled()
  );

  return (
    <TabsContent value="integrations" className="text-sm flex flex-col gap-4">
      <div className="flex items-center w-full gap-4">
        <label className="w-60">Synch Library Inegrations on Launch</label>
        <Checkbox
          checked={checkedOnLaunch}
          onCheckedChange={(value: boolean) => {
            setCheckedonLaunch(value);
            setIntegrateOnLaunchEnabled(value);
          }}
        />
      </div>
      <div className="flex items-center w-full gap-4">
        <label className="w-60">Synch Library after game exit</label>
        <Checkbox
          checked={checkedOnExit}
          onCheckedChange={(value: boolean) => {
            setCheckedonExit(value);
            setIntegrateOnExitEnabled(value);
          }}
        />
      </div>
    </TabsContent>
  );
}
