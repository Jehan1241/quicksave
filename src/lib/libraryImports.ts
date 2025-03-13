import { useSortContext } from "@/hooks/useSortContex";

export const importSteamLibrary = async (
  steamID: string,
  apiKey: string,
  setSteamLoading: (loading: "true" | "false" | "error") => void,
  setIntegrationLoadCount: (fn: (prev: number) => number) => void
) => {
  if (!steamID || !apiKey) {
    return;
  }

  setSteamLoading("true");
  setIntegrationLoadCount((prev: number) => prev + 1);
  try {
    const response = await fetch("http://localhost:8080/SteamImport", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        SteamID: steamID.trim(),
        APIkey: apiKey.trim(),
      }),
    });

    const resp = await response.json();
    console.log("Steam Errors:", resp.error);
    setSteamLoading("error");
  } catch (error) {
    console.error("Error:", error);
    setSteamLoading("error");
  }
  setSteamLoading("false");
  setIntegrationLoadCount((prev: number) => prev - 1);
};

export const importPlaystationLibrary = async (
  npsso: string,
  setPsnLoading: (loading: "true" | "false" | "error") => void,
  setPsnGamesNotMatched: (games: string[]) => void,
  setIntegrationLoadCount: (fn: (prev: number) => number) => void
) => {
  if (!npsso) {
    console.error("NPSSO token is required.");
    return;
  }

  console.log("Sending PlayStation Import Req", npsso);
  setIntegrationLoadCount((prev: number) => prev + 1);
  setPsnLoading("true");

  try {
    const response = await fetch("http://localhost:8080/PlayStationImport", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        npsso: npsso,
      }),
    });

    const resp = await response.json();
    console.log("PSN Error:", resp.error);
    console.log("PSN Games Not Matched:", resp.gamesNotMatched);

    setPsnLoading(resp.error);
    if (!resp.error) {
      setPsnGamesNotMatched(resp.gamesNotMatched);
    }
  } catch (error) {
    console.error("Error:", error);
    setPsnLoading("error");
  }

  setPsnLoading("false");
  setIntegrationLoadCount((prev: number) => prev - 1);
};
