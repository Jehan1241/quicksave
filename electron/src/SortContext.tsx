import React, { createContext, useState, useContext, ReactNode, useEffect } from "react";

interface SortContextType {
    sortType: string;
    setSortType: React.Dispatch<React.SetStateAction<string>>;
    sortOrder: "ASC" | "DESC" | "default";
    setSortOrder: React.Dispatch<React.SetStateAction<"ASC" | "DESC" | "default">>;
    fetchData: (type: string, order: "ASC" | "DESC" | "default", size: number) => Promise<any>;
    metaData: any;
    setMetaData: React.Dispatch<React.SetStateAction<any>>;
    tileSize: number;
    setTileSize: React.Dispatch<React.SetStateAction<number>>;
    sortStateUpdate: boolean;
    setSortStateUpdate: React.Dispatch<React.SetStateAction<boolean>>;
    viewState: string;
    setViewState: React.Dispatch<React.SetStateAction<string>>;
    isAddGameDialogOpen: boolean;
    setIsAddGameDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
    isWishlistAddDialogOpen: boolean;
    setIsWishlistAddDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
    isIntegrationsDialogOpen: boolean;
    setIsIntegrationsDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
    searchText: string;
    setSearchText: React.Dispatch<React.SetStateAction<string>>;
    randomGameClicked: boolean;
    setRandomGameClicked: React.Dispatch<React.SetStateAction<boolean>>;
}

const SortContext = createContext<SortContextType | undefined>(undefined);

export const SortProvider = ({ children }: { children: ReactNode }) => {
    const [sortType, setSortType] = useState("default");
    const [viewState, setViewState] = useState("library");
    const [tileSize, setTileSize] = useState(-1);
    const [sortOrder, setSortOrder] = useState<"ASC" | "DESC" | "default">("default");
    const [metaData, setMetaData] = useState([]);
    const [sortStateUpdate, setSortStateUpdate] = useState(false);
    const [isAddGameDialogOpen, setIsAddGameDialogOpen] = useState(false);
    const [isWishlistAddDialogOpen, setIsWishlistAddDialogOpen] = useState(false);
    const [isIntegrationsDialogOpen, setIsIntegrationsDialogOpen] = useState(false);
    const [searchText, setSearchText] = useState<string>("");
    const [randomGameClicked, setRandomGameClicked] = useState<boolean>(false);

    const fetchData = async (type: string, order: "ASC" | "DESC" | "default", size: number) => {
        console.log("Sending Get Basic Info");
        let json = null;
        try {
            const response = await fetch(
                `http://localhost:8080/getBasicInfo?type=${type}&order=${order}&size=${size}`
            );
            json = await response.json();
            /* console.log(json); */
        } catch (error) {
            console.error(error);
        }
        return json;
    };

    return (
        <SortContext.Provider
            value={{
                sortType,
                setSortType,
                sortOrder,
                setSortOrder,
                fetchData,
                metaData,
                tileSize,
                setMetaData,
                setTileSize,
                sortStateUpdate,
                setSortStateUpdate,
                viewState,
                setViewState,
                isAddGameDialogOpen,
                setIsAddGameDialogOpen,
                isIntegrationsDialogOpen,
                setIsIntegrationsDialogOpen,
                searchText,
                setSearchText,
                randomGameClicked,
                setRandomGameClicked,
                isWishlistAddDialogOpen,
                setIsWishlistAddDialogOpen,
            }}
        >
            {children}
        </SortContext.Provider>
    );
};

// Custom hook to use the context
export const useSortContext = (): SortContextType => {
    const context = useContext(SortContext);
    if (!context) throw new Error("useSortContext must be used within a SortProvider");
    return context;
};
