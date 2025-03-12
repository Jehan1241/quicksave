import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { TabsContent } from "@/components/ui/tabs";
import { useEffect, useRef, useState } from "react";

export function PathTab({ uid, setEditDialogOpen }: any) {
  const [gamePath, setGamePath] = useState("");
  const [pathIsValid, setPathIsValid] = useState<boolean | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const loadPath = async () => {
    console.log("Loading Game Path");
    try {
      const response = await fetch(
        `http://localhost:8080/getGamePath?uid=${uid}`
      );
      const json = await response.json();
      setGamePath(json.path);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    loadPath();
  }, []);

  const gamePathCheckHandler = async () => {
    console.log("Checking Game Path Validity");
    if (gamePath == "") {
      setPathIsValid(true);
      return;
    }
    const result = await window.electron.validateGamePath(gamePath);
    if (result.isValid) {
      console.log(result.message);
      setGamePath(result.message);
      setPathIsValid(true);
    } else {
      setPathIsValid(false);
    }
  };

  const gamePathSaveHandler = async () => {
    console.log("Saving Game Path", gamePath);
    try {
      const response = await fetch(
        `http://localhost:8080/setGamePath?uid=${uid}&path=${gamePath}`
      );
      const json = await response.json();
      console.log(json);
    } catch (error) {
      console.error(error);
    }
  };

  const browseFileHandler = async () => {
    const result = await window.electron.browseFileHandler({
      title: "Select Game File",
      filters: [
        { name: "Game Files", extensions: ["exe", "lnk", "bat", "cmd"] },
      ],
      properties: ["openFile"],
    });

    if (!result.canceled && result.filePaths.length > 0) {
      console.log("Selected File Path:", result.filePaths[0]);
      setGamePath(result.filePaths[0]);
      setPathIsValid(null);
    }
  };

  const fileSelectedHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      console.log("AA", file.path);
      setGamePath(file.path);
      setPathIsValid(null);
    }
  };

  return (
    <TabsContent value="path" className="h-full focus:ring-0">
      <div className="flex h-full flex-col justify-between gap-4 p-2 px-4 focus:outline-none">
        <div className="flex  items-center gap-2">
          <Label className="w-32">Game Path</Label>
          <Input
            id="gamePath"
            className={`${pathIsValid && "text-green-500"} ${
              pathIsValid == false && "text-red-500"
            }`}
            value={gamePath}
            onChange={(e) => {
              setGamePath(e.target.value);
              setPathIsValid(null);
            }}
          ></Input>
          <Button onClick={browseFileHandler} variant="outline">
            Browse
          </Button>
          <input
            type="file"
            ref={fileInputRef}
            style={{ display: "none" }}
            accept=".exe,.lnk,.bat,.cmd"
            onChange={fileSelectedHandler}
          />
        </div>
        <div className="self-end">
          {!pathIsValid && (
            <Button onClick={gamePathCheckHandler} variant={"dialogSaveButton"}>
              Validate Path
            </Button>
          )}
          {pathIsValid === true && (
            <Button
              onClick={() => {
                gamePathSaveHandler();
                setEditDialogOpen(false);
              }}
              variant={"dialogSaveButton"}
            >
              Save Path
            </Button>
          )}
        </div>
      </div>
    </TabsContent>
  );
}
