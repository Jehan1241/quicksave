import { Toast } from "@/components/ui/toast";
import { handleApiError } from "./apiErrors";
import { showErrorToast } from "../toastService";

export const sendGameToDB = async (
  title: string,
  releaseDate: any,
  selectedPlatforms: any,
  timePlayed: any,
  rating: any,
  selectedDevs: any,
  selectedTags: any,
  description: string,
  coverImage: any,
  ssImage: any,
  isWishlist: number,
  setAddGameLoading: React.Dispatch<React.SetStateAction<boolean>>,
  toast: any
) => {
  try {
    setAddGameLoading(true);
    const response = await fetch(`http://localhost:50001/addGameToDB`, {
      method: "POST",
      headers: { "Content-type": "application/json" },
      body: JSON.stringify({
        title: title,
        releaseDate: releaseDate,
        selectedPlatforms: selectedPlatforms,
        timePlayed: timePlayed,
        rating: rating,
        selectedDevs: selectedDevs,
        selectedTags: selectedTags,
        description: description,
        coverImage: coverImage,
        ssImage: ssImage,
        isWishlist: isWishlist,
      }),
    });
    if (!response.ok) await handleApiError(response);
    const resp = await response.json();
    if (resp.insertionStatus === false) {
      console.log(resp.insertionStatus);
      toast({
        variant: "destructive",
        title: "Game Insertion Error!",
        description: "This game has already been inserted.",
      });
    }
    if (resp.insertionStatus === true) {
      console.log(resp.insertionStatus);
      toast({
        variant: "default",
        title: "Game Added!",
        description: "The game has been added to the database.",
      });
    }
  } catch (error: any) {
    console.error(error);
    showErrorToast("An error occured!", String(error));
  } finally {
    setAddGameLoading(false);
  }
};

export const searchGame = async (
  title: string,
  setTitleEmpty: React.Dispatch<React.SetStateAction<boolean>>,
  setLoading: React.Dispatch<React.SetStateAction<boolean>>,
  setData: React.Dispatch<React.SetStateAction<any>>,
  toast: any
) => {
  if (!title.trim()) {
    setTitleEmpty(true);
    return;
  }

  setLoading(true);
  try {
    const response = await fetch(`http://localhost:50001/IGDBsearch`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ NameToSearch: title }),
    });

    if (!response.ok) await handleApiError(response);

    const resp = await response.json();
    setData(resp.foundGames);
  } catch (error: any) {
    console.error("Fetch Error:", error);
    showErrorToast("Failed to search IGDB!", String(error));
  } finally {
    setLoading(false);
  }
};

export const fetchTagsDevsPlatforms = async (
  setTagOptions: React.Dispatch<React.SetStateAction<any>>,
  setDevOptions: React.Dispatch<React.SetStateAction<any>>,
  setPlatformOptions: React.Dispatch<React.SetStateAction<any>>
) => {
  try {
    const response = await fetch("http://localhost:50001/getAllTags");
    if (!response.ok) await handleApiError(response);
    const resp = await response.json();

    // Transform the tags into key-value pairs
    if (resp.tags && resp.tags.length > 0) {
      const tagsAsKeyValuePairs = resp.tags.map((tag: any) => ({
        value: tag,
        label: tag,
      }));
      setTagOptions(tagsAsKeyValuePairs);
    }
  } catch (error) {
    console.error("Error fetching tags:", error);
    showErrorToast("Failed to get tags!", String(error));
  }
  try {
    const response = await fetch("http://localhost:50001/getAllDevelopers");
    if (!response.ok) await handleApiError(response);
    const resp = await response.json();
    console.log(resp);

    if (resp.devs && resp.devs.length > 0) {
      const devsAsKeyValuePairs = resp.devs.map((dev: any) => ({
        value: dev,
        label: dev,
      }));
      setDevOptions(devsAsKeyValuePairs);
    }
  } catch (error) {
    console.error("Error fetching developers:", error);
    showErrorToast("Failed to get developers!", String(error));
  }
  try {
    const response = await fetch("http://localhost:50001/getAllPlatforms");
    if (!response.ok) await handleApiError(response);
    const resp = await response.json();
    console.log(resp);

    // Transform the tags into key-value pairs
    if (resp.platforms && resp.platforms.length > 0) {
      const platsAsKeyValuePairs = resp.platforms.map((plat: any) => ({
        value: plat,
        label: plat,
      }));
      setPlatformOptions(platsAsKeyValuePairs);
    }
  } catch (error) {
    console.error("Error fetching platforms:", error);
    showErrorToast("Failed to get platforms!", String(error));
  }
};
