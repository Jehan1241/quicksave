import React, {
  createContext,
  useState,
  useContext,
  ReactNode,
  useEffect,
} from "react";

interface SortContextType {
  sortType: string;
  setSortType: React.Dispatch<React.SetStateAction<string>>;
  sortOrder: "ASC" | "DESC" | "default";
  setSortOrder: React.Dispatch<
    React.SetStateAction<"ASC" | "DESC" | "default">
  >;
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
  playingGame: string | null;
  setPlayingGame: React.Dispatch<React.SetStateAction<string | null>>;
  cacheBuster: number;
  setCacheBuster: React.Dispatch<React.SetStateAction<number>>;
  integrationLoadCount: number;
  setIntegrationLoadCount: React.Dispatch<React.SetStateAction<number>>;
  filterActive: boolean;
  setFilterActive: React.Dispatch<React.SetStateAction<boolean>>;
  deleteFilterGames: boolean;
  setDeleteFilterGames: React.Dispatch<React.SetStateAction<boolean>>;
  hideFilterGames: null | "hide" | "unhide";
  setHideFilterGames: React.Dispatch<
    React.SetStateAction<null | "hide" | "unhide">
  >;
}

const SortContext = createContext<SortContextType | undefined>(undefined);

export const SortProvider = ({ children }: { children: ReactNode }) => {
  const [sortType, setSortType] = useState("default");
  const [viewState, setViewState] = useState("library");
  const [tileSize, setTileSize] = useState(-1);
  const [sortOrder, setSortOrder] = useState<"ASC" | "DESC" | "default">(
    "default"
  );
  const [metaData, setMetaData] = useState([]);
  const [sortStateUpdate, setSortStateUpdate] = useState(false);
  const [isAddGameDialogOpen, setIsAddGameDialogOpen] = useState(false);
  const [isWishlistAddDialogOpen, setIsWishlistAddDialogOpen] = useState(false);
  const [isIntegrationsDialogOpen, setIsIntegrationsDialogOpen] =
    useState(false);
  const [searchText, setSearchText] = useState<string>("");
  const [randomGameClicked, setRandomGameClicked] = useState<boolean>(false);
  const [cacheBuster, setCacheBuster] = useState<number>(Date.now());
  const [integrationLoadCount, setIntegrationLoadCount] = useState<number>(0);
  const [playingGame, setPlayingGame] = useState<string | null>(null);
  const [filterActive, setFilterActive] = useState<boolean>(false);
  const [deleteFilterGames, setDeleteFilterGames] = useState<boolean>(false);
  const [hideFilterGames, setHideFilterGames] = useState<
    null | "hide" | "unhide"
  >(null);

  return (
    <SortContext.Provider
      value={{
        sortType,
        setSortType,
        sortOrder,
        setSortOrder,
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
        cacheBuster,
        setCacheBuster,
        integrationLoadCount,
        setIntegrationLoadCount,
        playingGame,
        setPlayingGame,
        filterActive,
        setFilterActive,
        deleteFilterGames,
        setDeleteFilterGames,
        hideFilterGames,
        setHideFilterGames,
      }}
    >
      {children}
    </SortContext.Provider>
  );
};

// Custom hook to use the context
export const useSortContext = (): SortContextType => {
  const context = useContext(SortContext);
  if (!context)
    throw new Error("useSortContext must be used within a SortProvider");
  return context;
};
