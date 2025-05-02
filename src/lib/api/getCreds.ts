import { showErrorToast } from "../toastService";
import { handleApiError } from "./apiErrors";

export const getSteamCreds = async () => {
  console.log("Getting Steam Creds");

  try {
    const response = await fetch(`http://localhost:50001/SteamCreds`);
    if (!response.ok) await handleApiError(response);
    const json = await response.json();

    return {
      ID: json.SteamCreds[0],
      APIKey: json.SteamCreds[1],
    };
  } catch (error) {
    console.error("Error fetching Steam credentials:", error);
    showErrorToast("Failed to get steam credentials!", String(error));
    return null; // Return null if there's an error
  }
};

export const getNpsso = async () => {
  console.log("Getting Npsso");
  try {
    const response = await fetch(`http://localhost:50001/Npsso`);
    if (!response.ok) await handleApiError(response);
    const json = await response.json();
    console.log(json.Npsso);
    return json.Npsso;
  } catch (error) {
    console.error(error);
    showErrorToast("Failed to get npsso!", String(error));
    return null;
  }
};
