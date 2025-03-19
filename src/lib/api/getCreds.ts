export const getSteamCreds = async () => {
  console.log("Getting Steam Creds");

  try {
    const response = await fetch(`http://localhost:8080/SteamCreds`);
    const json = await response.json();

    return {
      ID: json.SteamCreds[0],
      APIKey: json.SteamCreds[1],
    };
  } catch (error) {
    console.error("Error fetching Steam credentials:", error);
    return null; // Return null if there's an error
  }
};

export const getNpsso = async () => {
  console.log("Getting Npsso");
  try {
    const response = await fetch(`http://localhost:8080/Npsso`);
    const json = await response.json();
    console.log(json.Npsso);
    return json.Npsso;
  } catch (error) {
    console.error(error);
    return null;
  }
};
