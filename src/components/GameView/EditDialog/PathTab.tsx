import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { TabsContent } from "@/components/ui/tabs";
import { useEffect, useState } from "react";

export function PathTab({ uid, setEditDialogOpen }: any) {
  const [gamePath, setGamePath] = useState("");
  const [pathIsValid, setPathIsValid] = useState<boolean | null>(null);

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

  return (
    <TabsContent value="path" className="h-full focus:ring-0">
      <div className="flex h-full flex-col justify-between p-2 px-4 focus:outline-none">
        <div className="flex  items-center">
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
        </div>
        <div className="self-end">
          {!pathIsValid && (
            <Button onClick={gamePathCheckHandler} variant={"dialogSaveButton"}>
              Check Path
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
              Save
            </Button>
          )}
        </div>
      </div>
    </TabsContent>
  );
}
