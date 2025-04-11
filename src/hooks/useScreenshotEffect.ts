import { useEffect } from "react";
import { useSortContext } from "@/hooks/useSortContex";
import { getScreenshotBind, getScreenshotEnabled } from "@/lib/screeenshots";
import { getSteamCreds } from "@/lib/api/getCreds";
import { importSteamLibrary } from "@/lib/api/libraryImports";

export function useScreenshotEffect() {
  const { playingGame, setIntegrationLoadCount } = useSortContext();

  const updateSteam = async () => {
    const steamCreds = await getSteamCreds();

    await importSteamLibrary(
      steamCreds?.ID,
      steamCreds?.APIKey,
      () => {},
      setIntegrationLoadCount,
      () => {}
    );
  };

  useEffect(() => {
    if (!getScreenshotEnabled()) return;

    const screenshotBind = getScreenshotBind();

    if (playingGame === "") {
      window.windowFunctions.updatePlayingGame("", screenshotBind);
      updateSteam();
    } else if (playingGame != null) {
      window.windowFunctions.updatePlayingGame(playingGame, screenshotBind);
    }
  }, [playingGame]);
}
