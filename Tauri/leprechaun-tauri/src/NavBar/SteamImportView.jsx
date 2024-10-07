function SteamImportView() {
  const searchClickHandler = async () => {
    const SteamID = document.getElementById("SteamID").value;
    const APIkey = document.getElementById("APIKey").value;
    console.log(SteamID, APIkey);

    try {
      const response = await fetch("http://localhost:8080/SteamImport", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ SteamID: SteamID, APIkey: APIkey }),
      });
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const checkForEnterPressed = (e) => {
    if (e.key == "Enter") {
      searchClickHandler();
    }
  };

  return (
    <div className="flex flex-col p-4 mt-2 w-full h-full text-base rounded-xl">
      {/* Input Div */}
      <div className="flex flex-row gap-4">
        <div className="flex flex-col gap-4">
          <p>Steam ID</p>
          <p>Steam API Key</p>
        </div>
        <div className="flex flex-col gap-4">
          <input
            onKeyDown={checkForEnterPressed}
            id="SteamID"
            className="px-1 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
          ></input>
          <input
            onKeyDown={checkForEnterPressed}
            id="APIKey"
            className="px-1 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
          ></input>
        </div>
        <div className="flex items-end text-sm text-blue-700 underline">
          <a href="https://steamcommunity.com/dev/apikey">API Key?</a>
        </div>
      </div>
      {/* button div */}
      <div className="flex justify-end mt-auto">
        <button
          className="w-32 h-10 rounded-lg border-2 bg-primary"
          onClick={searchClickHandler}
        >
          Import
        </button>
      </div>
    </div>
  );
}

export default SteamImportView;
