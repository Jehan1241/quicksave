function AddGameSteam() {
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

  return (
    <div className="relative w-[calc(100vw-100px)] h-screen text-center text-white bg-gameView ml-1 rounded-l-3xl p-12 z-0 font-mono">
      <div className="p-6 w-1/2 rounded-2xl border-2 border-gray-700 bg-primary hover:border-gray-500">
        <p className="text-2xl text-left text-white">Import Steam Library</p>
        <div className="text-left">
          <br />
          <p className="inline m-2 text-white">Your Steam ID</p>
          <input
            id="SteamID"
            className="p-1 ml-7 w-72 rounded-xl border-2 border-gray-700 bg-gameView"
          />
          <br />
          <br />
          <p className="inline m-2 text-white">Your API Key</p>
          <input
            id="APIKey"
            className="p-1 ml-9 w-72 rounded-xl border-2 border-gray-700 bg-gameView"
          />
        </div>
        <div className="text-right">
          <button
            className="w-32 h-10 rounded-lg border-2 border-gray-700 bg-gameView hover:font-extrabold hover:border-white"
            onClick={searchClickHandler}
          >
            Search
          </button>
        </div>
      </div>
    </div>
  );
}

export default AddGameSteam;
