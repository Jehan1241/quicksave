import { showErrorToast } from "../toastService";
import { handleApiError } from "./apiErrors";

export const loadFilterState = async (
  setIsLoaded: React.Dispatch<React.SetStateAction<boolean>>,
  setSelectedDevs: React.Dispatch<
    React.SetStateAction<{ value: string; label: string }[]>
  >,
  setSelectedPlatforms: React.Dispatch<
    React.SetStateAction<{ value: string; label: string }[]>
  >,
  setSelectedTags: React.Dispatch<
    React.SetStateAction<{ value: string; label: string }[]>
  >,
  setSelectedName: React.Dispatch<
    React.SetStateAction<{ value: string; label: string }[]>
  >
) => {
  setIsLoaded(false);
  try {
    console.log("Sending Load Filters");
    const response = await fetch("http://localhost:8080/LoadFilters");
    if (!response.ok) await handleApiError(response);
    const data = await response.json();
    if (data.developers) {
      const devsAsKeyValuePairs = data.developers.map((dev: any) => ({
        value: dev,
        label: dev,
      }));
      setSelectedDevs(devsAsKeyValuePairs);
    }
    if (data.platform) {
      const platsAsKeyValuePairs = data.platform.map((plat: any) => ({
        value: plat,
        label: plat,
      }));
      setSelectedPlatforms(platsAsKeyValuePairs);
    }
    if (data.tags) {
      const tagsAsKeyValuePairs = data.tags.map((tag: any) => ({
        value: tag,
        label: tag,
      }));
      setSelectedTags(tagsAsKeyValuePairs);
    }
    if (data.name) {
      const nameAsKeyValuePairs = data.name.map((name: any) => ({
        value: name,
        label: name,
      }));
      setSelectedName(nameAsKeyValuePairs);
    }
  } catch (error) {
    console.error("Error fetching filter:", error);
    showErrorToast("Failed to load filters!", String(error));
  } finally {
    setIsLoaded(true);
  }
};

interface Option {
  value: string;
  label: string;
}

export const handleFilterChange = async (
  selectedPlatforms: Option[],
  selectedTags: Option[],
  selectedName: Option[],
  selectedDevs: Option[],
  setFilterActive: any
) => {
  const platformValues = selectedPlatforms.map((platform) => platform.value);
  const tagValues = selectedTags.map((tag) => tag.value);
  const nameValues = selectedName.map((name) => name.value);
  const devValues = selectedDevs.map((dev) => dev.value);

  const filter = {
    tags: tagValues,
    name: nameValues,
    platforms: platformValues,
    devs: devValues,
  };

  const filtersActive =
    selectedPlatforms.length > 0 ||
    selectedTags.length > 0 ||
    selectedName.length > 0 ||
    selectedDevs.length > 0;
  setFilterActive(filtersActive);

  try {
    console.log("Sending Set Filter");
    const response = await fetch("http://localhost:8080/setFilter", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(filter),
    });
    if (!response.ok) await handleApiError(response);
  } catch (error: any) {
    console.error("Error fetching filter:", error);
    showErrorToast("Failed to clear filters!", String(error));
  }
};

export const clearAllFilters = async (
  setSelectedDevs: React.Dispatch<React.SetStateAction<Option[]>>,
  setSelectedName: React.Dispatch<React.SetStateAction<Option[]>>,
  setSelectedPlatforms: React.Dispatch<React.SetStateAction<Option[]>>,
  setSelectedTags: React.Dispatch<React.SetStateAction<Option[]>>
) => {
  setSelectedDevs([]);
  setSelectedName([]);
  setSelectedPlatforms([]);
  setSelectedTags([]);
  try {
    console.log("Sending Clear All Filters");
    // Send the filter as a POST request
    const response = await fetch("http://localhost:8080/clearAllFilters");

    if (!response.ok) await handleApiError(response);
  } catch (error: any) {
    console.error("Error clearing filter:", error);
    showErrorToast("Failed to clear filters!", String(error));
  }
};

export const deleteCurrentlyFiltered = async (games: any[]) => {
  const uids = games.map((game) => game.UID);
  try {
    const response = await fetch(
      `http://localhost:8080/deleteCurrentlyFiltered`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ uids }),
      }
    );

    if (!response.ok) await handleApiError(response);
  } catch (error: any) {
    console.error("error deleting games:", error);
    showErrorToast("Failed to delete games!", String(error));
  }
};

export const hideCurrentlyFiltered = async (games: any[]) => {
  const uids = games.map((game) => game.UID);
  try {
    const response = await fetch(
      `http://localhost:8080/hideCurrentlyFiltered`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ uids }), // Send UIDs as JSON
      }
    );
    if (!response.ok) await handleApiError(response);
  } catch (error: any) {
    console.error("error hiding games:", error);
    showErrorToast("Failed to hide games!", String(error));
  }
};

export const unHideCurrentlyFiltered = async (games: any[]) => {
  const uids = games.map((game) => game.UID);
  try {
    const response = await fetch(
      `http://localhost:8080/unHideCurrentlyFiltered`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ uids }), // Send UIDs as JSON
      }
    );
    if (!response.ok) await handleApiError(response);
  } catch (error: any) {
    console.error("error unhiding games:", error);
    showErrorToast("Failed to unhide games!", String(error));
  }
};
