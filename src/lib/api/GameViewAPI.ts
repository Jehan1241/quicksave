import { CheckedState } from "@radix-ui/react-checkbox";
import { showErrorToast } from "../toastService";
import { handleApiError } from "./apiErrors";
import { useNavigate } from "react-router-dom";

interface GameData {
  companies: string[];
  tags: string[];
  metadata: {
    name: string;
    description: string;
    [key: string]: any;
  };
  screenshots: string[];
}

export const getGameDetails = async (
  uid: string,
  setCompanies: React.Dispatch<React.SetStateAction<string>>,
  setTags: React.Dispatch<React.SetStateAction<string>>,
  setMetadata: React.Dispatch<React.SetStateAction<string>>,
  setScreenshots: React.Dispatch<React.SetStateAction<string>>
) => {
  console.log("Sending Get Game Details");
  try {
    const response = await fetch(
      `http://localhost:8080/GameDetails?uid=${uid}`
    );
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    const { companies, tags, screenshots, m: metadata } = json.metadata;
    setCompanies(companies[uid]);
    setTags(tags[uid]);
    setMetadata(metadata[uid]);
    setScreenshots(screenshots[uid] || []); // Make sure it's an array
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to load game data!", String(error));
  }
};

export const doDataPreload = async (
  uid: string,
  setPreloadData: React.Dispatch<React.SetStateAction<GameData | null>>
) => {
  try {
    const response = await fetch(
      `http://localhost:8080/GameDetails?uid=${uid}`
    );
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    console.log("BBB", json);
    const { companies, tags, screenshots, m: metadata } = json.metadata;
    setPreloadData({
      companies: companies[uid],
      tags: tags[uid],
      metadata: metadata[uid],
      screenshots: screenshots[uid] || [],
    });
  } catch (error) {
    console.error("Failed to preload game data:", error);
    showErrorToast("Failed to preload game data!", String(error));
  }
};

export const hardDelete = async (uid: string, navigate: any) => {
  console.log("Sending Delete Game");
  try {
    const response = await fetch(`http://localhost:8080/DeleteGame?uid=${uid}`);
    if (!response.ok) await handleApiError(response);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to delete game!", String(error));
  }
  navigate(-1);
};

export const hideGame = async (uid: string, navigate: any) => {
  console.log("Sending Hide Game");
  try {
    const response = await fetch(`http://localhost:8080/HideGame?uid=${uid}`);
    if (!response.ok) await handleApiError(response);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to hide game!", String(error));
  }
  navigate(-1);
};

export const unhideGame = async (uid: string, navigate: any) => {
  try {
    console.log("Sending unhide game");
    const response = await fetch(`http://localhost:8080/unhideGame?uid=${uid}`);
    if (!response.ok) await handleApiError(response);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to unhide game!", String(error));
  }
  navigate("/hidden", { replace: true });
};

export const launchGame = async (
  uid: string,
  setCompanies: React.Dispatch<React.SetStateAction<string>>,
  setTags: React.Dispatch<React.SetStateAction<string>>,
  setMetadata: React.Dispatch<React.SetStateAction<string>>,
  setScreenshots: React.Dispatch<React.SetStateAction<string>>,
  setPlayingGame: React.Dispatch<React.SetStateAction<string | null>>
) => {
  setPlayingGame(uid);
  console.log("Play Game Clicked");
  try {
    const response = await fetch(`http://localhost:8080/LaunchGame?uid=${uid}`);
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    const launchStatus = json.LaunchStatus;
    if (launchStatus === "Launched") {
      getGameDetails(uid, setCompanies, setTags, setMetadata, setScreenshots);
    }
  } catch (error) {
    console.log(error);
    showErrorToast("Failed to launch game!", String(error));
  } finally {
    setPlayingGame("");
  }
};

export const sendSteamInstallReq = async (uid: string) => {
  console.log("Play Game Clicked");
  try {
    const response = await fetch(
      `http://localhost:8080/steamInstallReq?uid=${uid}`
    );
    if (!response.ok) await handleApiError(response);
  } catch (error) {
    console.log(error);
    showErrorToast("Failed to launch game!", String(error));
  }
};

export const gamePathSaveHandler = async (uid: string, gamePath: string) => {
  console.log("Saving Game Path", gamePath);
  try {
    const response = await fetch(
      `http://localhost:8080/setGamePath?uid=${uid}&path=${gamePath}`
    );
    if (!response.ok) await handleApiError(response);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to update path!", String(error));
  }
};

export const getGamePath = async (
  uid: string,
  setGamePath: React.Dispatch<React.SetStateAction<string>>
) => {
  console.log("Loading Game Path");
  try {
    const response = await fetch(
      `http://localhost:8080/getGamePath?uid=${uid}`
    );
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    setGamePath(json.path);
  } catch (error) {
    showErrorToast("Failed to get game path!", String(error));
    console.error(error);
  }
};

export const saveCustomImage = async (
  uid: string,
  setLoading: React.Dispatch<React.SetStateAction<boolean>>,
  currentCover: string | null,
  ssImage: (string | null)[],
  navigate: any,
  setCacheBuster: React.Dispatch<React.SetStateAction<number>>,
  fetchData: any
) => {
  setLoading(true);
  try {
    const response = await fetch(`http://localhost:8080/setCustomImage`, {
      method: "POST",
      headers: { "Content-type": "application/json" },
      body: JSON.stringify({
        uid: uid,
        coverImage: currentCover,
        ssImage: ssImage,
      }),
    });
    if (!response.ok) await handleApiError(response);
    const resp = await response.json();
    if (resp.status === "OK") {
      fetchData();
      navigate("/library", { replace: true });
      setCacheBuster(Date.now());
      //window.windowFunctions.nukeCache();
    }
    console.log(resp);
  } catch (error) {
    showErrorToast("Failed to set image!", String(error));
    console.error(error);
  } finally {
    setLoading(false);
  }
};

export const loadPreferences = async (
  uid: string,
  setCustomTime: React.Dispatch<React.SetStateAction<string>>,
  setCustomTimeOffset: React.Dispatch<React.SetStateAction<string>>,
  setCustomRating: React.Dispatch<React.SetStateAction<string>>,
  setCustomTitle: React.Dispatch<React.SetStateAction<string>>,
  setCustomReleaseDate: React.Dispatch<React.SetStateAction<Date | undefined>>,
  setCustomTitleChecked: React.Dispatch<
    React.SetStateAction<CheckedState | undefined>
  >,
  setCustomTimeChecked: React.Dispatch<
    React.SetStateAction<CheckedState | undefined>
  >,
  setCustomTimeOffsetChecked: React.Dispatch<
    React.SetStateAction<CheckedState | undefined>
  >,
  setCustomReleaseDateChecked: React.Dispatch<
    React.SetStateAction<CheckedState | undefined>
  >,
  setCustomRatingChecked: React.Dispatch<
    React.SetStateAction<CheckedState | undefined>
  >,
  setTagOptions: React.Dispatch<
    React.SetStateAction<{ label: string; value: string }[]>
  >,
  setDevOptions: React.Dispatch<
    React.SetStateAction<{ label: string; value: string }[]>
  >
) => {
  try {
    const response = await fetch(
      `http://localhost:8080/LoadPreferences?uid=${uid}`
    );
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    console.log(json);
    setCustomTime("0");
    setCustomTimeOffset("0");
    setCustomRating("0");
    setCustomTitle(json.preferences.title.value);
    if (json.preferences.time.value) {
      setCustomTime(json.preferences.time.value);
    }
    if (json.preferences.timeOffset.value) {
      setCustomTimeOffset(json.preferences.timeOffset.value);
    }
    if (json.preferences.rating.value) {
      console.log("A");
      setCustomRating(json.preferences.rating.value);
    }
    setCustomReleaseDate(json.preferences.releaseDate.value);

    if (json.preferences.title.checked == "1") {
      setCustomTitleChecked(true);
    }
    if (json.preferences.time.checked == "1") {
      setCustomTimeChecked(true);
    }
    if (json.preferences.timeOffset.checked == "1") {
      setCustomTimeOffsetChecked(true);
    }
    if (json.preferences.releaseDate.checked == "1") {
      setCustomReleaseDateChecked(true);
    }
    if (json.preferences.rating.checked == "1") {
      setCustomRatingChecked(true);
    }
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to load preferences!", String(error));
  }
  try {
    const tagsResponse = await fetch("http://localhost:8080/getAllTags");
    if (!tagsResponse.ok) await handleApiError(tagsResponse);
    const tagsData = await tagsResponse.json();

    // Transform the tags into key-value pairs
    const tagsAsKeyValuePairs = tagsData.tags.map((tag: any) => ({
      value: tag,
      label: tag,
    }));
    setTagOptions(tagsAsKeyValuePairs);
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to load preferences!", String(error));
  }
  try {
    const devsResponse = await fetch("http://localhost:8080/getAllDevelopers");
    if (!devsResponse.ok) await handleApiError(devsResponse);
    const devsData = await devsResponse.json();
    console.log(devsData);

    // Transform the developers into key-value pairs
    const devsAsKeyValuePairs = devsData.devs.map((dev: any) => ({
      value: dev,
      label: dev,
    }));
    setDevOptions(devsAsKeyValuePairs);
  } catch (error) {
    showErrorToast("Failed to load preferences!", String(error));
    console.error(error);
  }
};
