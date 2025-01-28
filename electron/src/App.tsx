import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { darkMode } from "./ToggleTheme";
import { BrowserRouter, data, useLocation, useNavigate } from "react-router-dom";
import { Route, Routes } from "react-router-dom";
import CustomTitleBar from "./CustomTitleBar";
import GridView from "./LibraryView/LibraryView";
import { SortProvider } from "./SortContext"; // Import the SortProvider
import { useSortContext } from "./SortContext";
import DetialsView from "./WishlistView/WishlistView";
import AddGameManuallyDialog from "./Dialogs/AddGameManually";
import { Toaster } from "./components/ui/toaster";
import Integrations from "./Dialogs/Integrations";
import GameView from "./GameView/GameView";
import LibraryView from "./LibraryView/LibraryView";
import WishlistView from "./WishlistView/WishlistView";

export default function App() {
    const {
        sortType,
        sortOrder,
        setSortType,
        setSortOrder,
        setMetaData,
        setTileSize,
        tileSize,
        sortStateUpdate,
        setSortStateUpdate,
        viewState,
        randomGameClicked,
        setRandomGameClicked,
        isAddGameDialogOpen,
        isIntegrationsDialogOpen,
    } = useSortContext();
    const [dataArray, setDataArray] = useState<any[]>([]); // Track dataArray in local state\
    const [wishlistArray, setWishlistArray] = useState<any[]>([]); // Track dataArray in local state\
    const [integrationsPreviouslyOpened, setIntegrationsPreviouslyOpened] =
        useState<boolean>(false);
    const [addGameDialogHasBeenOpened, setAddGameDialogHasBeenOpened] = useState<boolean>(false);
    const navigate = useNavigate();

    const fetchData = async () => {
        console.log("Sending Get Basic Info");
        try {
            const response = await fetch(
                `http://localhost:8080/getBasicInfo?type=${sortType}&order=${sortOrder}&size=${tileSize}`
            );
            const json = await response.json();
            console.log(json);

            const filteredLibraryGames = Object.values(json.MetaData).filter(
                (item: any) => item.isDLC === 0
            );
            setDataArray(filteredLibraryGames);
            setMetaData(json.MetaData);
            setSortOrder(json.SortOrder);
            setSortType(json.SortType);
            setTileSize(json.Size);

            const filteredWishlistGames = Object.values(json.MetaData).filter(
                (item: any) => item.isDLC === 1
            );
            setWishlistArray(filteredWishlistGames);
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        darkMode();
        fetchData();
        const eventSource = new EventSource("http://localhost:8080/sse-steam-updates");

        eventSource.onmessage = (event) => {
            console.log("SSE message received:", event.data);
            fetchData();
        };
        eventSource.onerror = (error) => {
            console.error("SSE Error:", error);
        };
        return () => {
            eventSource.close();
        };
    }, []);

    useEffect(() => {
        if (randomGameClicked && viewState == "library") {
            console.log("randomGameClicked");
            setRandomGameClicked(false);
            const len = dataArray.length;
            const randomNumber = Math.floor(Math.random() * (len + 1));
            const uid = dataArray[randomNumber].UID;
            navigate(`gameview`, {
                state: { data: uid },
            });
            setRandomGameClicked(false);
        }
        if (randomGameClicked && viewState == "wishlist") {
            console.log("randomGameClicked");
            setRandomGameClicked(false);
            const len = wishlistArray.length;
            const randomNumber = Math.floor(Math.random() * (len + 1));
            const uid = wishlistArray[randomNumber].UID;
            navigate(`gameview`, {
                state: { data: uid },
            });
            setRandomGameClicked(false);
        }
    }, [randomGameClicked]);

    useEffect(() => {
        if (sortStateUpdate === true) {
            console.log("Sort State Update");
            fetchData();
            setSortStateUpdate(false);
        }
    }, [sortStateUpdate]);

    // These 2 are used to check if dialogs have been opened atleast once
    // This is necessary as if you leave the dialog in load state it wont save otherwise
    useEffect(() => {
        if (isAddGameDialogOpen) {
            setAddGameDialogHasBeenOpened(true);
        }
    }, [isAddGameDialogOpen]);
    useEffect(() => {
        if (isIntegrationsDialogOpen) {
            setIntegrationsPreviouslyOpened(true);
        }
    }, [isIntegrationsDialogOpen]);

    return (
        <>
            <CustomTitleBar>
                {addGameDialogHasBeenOpened && <AddGameManuallyDialog />}
                {integrationsPreviouslyOpened && <Integrations />}
                <Routes>
                    {viewState === "library" && (
                        <Route element={<LibraryView data={dataArray} />} path="/" />
                    )}
                    {viewState === "wishlist" && (
                        <Route element={<WishlistView data={wishlistArray} />} path="/" />
                    )}
                    <Route element={<GameView />} path="/gameview" />
                </Routes>
                <Toaster />
            </CustomTitleBar>
        </>
    );
}

const root = createRoot(document.getElementById("app")!);
root.render(
    <SortProvider>
        <BrowserRouter>
            {/* <React.StrictMode> */}
            <App />
            {/* </React.StrictMode> */}
        </BrowserRouter>
    </SortProvider>
);
