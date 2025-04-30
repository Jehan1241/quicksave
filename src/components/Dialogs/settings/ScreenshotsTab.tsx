import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { TabsContent } from "@/components/ui/tabs";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  getScreenshotBind,
  getScreenshotEnabled,
  setScreenshotBind,
  setScreenshotEnabled,
} from "@/lib/screeenshots";
import { QuestionMarkCircledIcon } from "@radix-ui/react-icons";
import { Ghost } from "lucide-react";
import { useEffect, useRef, useState } from "react";

export default function ScreeenshotsTab() {
  const [shortcut, setShortcut] = useState(getScreenshotBind());
  const [recording, setRecording] = useState(false);
  const keysPressed = useRef<Set<string>>(new Set());

  const handleKeyDown = (e: KeyboardEvent) => {
    e.preventDefault();
    keysPressed.current.add(e.key.length === 1 ? e.key.toUpperCase() : e.key);
  };

  const handleKeyUp = () => {
    if (recording) {
      const finalShortcut = Array.from(keysPressed.current).join("+");
      keysPressed.current.clear();
      setRecording(false);
      setShortcut(finalShortcut);
      setScreenshotBind(finalShortcut);
    }
  };

  useEffect(() => {
    if (!recording) return;

    window.addEventListener("keydown", handleKeyDown);
    window.addEventListener("keyup", handleKeyUp);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("keyup", handleKeyUp);
    };
  }, [recording]);

  const [checked, setChecked] = useState(getScreenshotEnabled());

  return (
    <TabsContent value="screenshots" className="text-sm flex flex-col gap-4">
      <div className="flex justify-between gap-8">
        <div className="gap-4 flex flex-col">
          <div className="flex gap-4 items-center">
            <label className="w-36">Enable Screenshots</label>
            <div className="w-40 flex justify-center">
              <Checkbox
                checked={checked}
                onCheckedChange={(value: boolean) => {
                  setChecked(value);
                  setScreenshotEnabled(value);
                }}
              />
            </div>
          </div>
          <div className="flex gap-4 items-center">
            <label className="w-36">Screenshot Keybind</label>
            <Button
              variant={"outline"}
              className="w-40"
              onClick={() => setRecording(true)}
            >
              {recording ? "Recording Keypress" : shortcut}
            </Button>
          </div>
        </div>

        <div className="w-full items-center h-full flex">
          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger className="text-sm flex gap-1 items-center">
                <QuestionMarkCircledIcon />
                Screenshots not working?
              </TooltipTrigger>
              <TooltipContent side="bottom">
                <p>Try using a 3 button combination</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </div>
    </TabsContent>
  );
}
