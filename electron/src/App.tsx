import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { darkMode } from "./ToggleTheme";
import { BrowserRouter, data, useLocation, useNavigate } from "react-router-dom";
import { Route, Routes } from "react-router-dom";
import CustomTitleBar from "./CustomTitleBar";
import GridView from "./GridView/GridView";
import { SortProvider } from "./SortContext"; // Import the SortProvider
import { useSortContext } from "./SortContext";
import DetialsView from "./DetailsView/DetailsView";
import AddGameManuallyDialog from "./Dialogs/AddGameManually";
import { Toaster } from "./components/ui/toaster";
import Integrations from "./Dialogs/Integrations";
import GameView from "./GameView/GameView";

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
            setDataArray(Object.values(json.MetaData));
            setMetaData(json.MetaData);
            setSortOrder(json.SortOrder);
            setSortType(json.SortType);
            setTileSize(json.Size);
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
        if (randomGameClicked) {
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
                    {viewState === "grid" && (
                        <Route element={<GridView data={dataArray} />} path="/" />
                    )}
                    {viewState === "details" && (
                        <Route element={<DetialsView data={dataArray} />} path="/" />
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
