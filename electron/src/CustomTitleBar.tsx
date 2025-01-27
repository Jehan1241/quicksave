import MultipleSelector, { Option } from "./components/ui/multiple-selector";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LayoutGrid, ListCollapse, PanelLeft, Search } from "lucide-react";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { ChartNoAxesColumnIncreasing, ChartNoAxesColumnDecreasing } from "lucide-react";
import React, { useEffect, useState, type ReactNode } from "react";
import { Save } from "lucide-react";
import { BsFloppyFill } from "react-icons/bs";
import { PiBookLight, PiListHeartLight } from "react-icons/pi";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuPortal,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ArrowDownWideNarrow } from "lucide-react";
import { Filter } from "lucide-react";
import { Dices } from "lucide-react";
import { darkMode, lightMode, redMode } from "./ToggleTheme";
import { useSortContext } from "./SortContext";
import { Slider } from "@/components/ui/slider";
import { GrCatalog, GrCatalogOption } from "react-icons/gr";
import { useNavigate, useLocation } from "react-router-dom";

const closeWindow = async () => {
    window.windowFunctions.closeApp();
};

const minimizeWindow = async () => {
    window.windowFunctions.minimize();
};

const maximizeWindow = async () => {
    window.windowFunctions.maximize();
};

export default function CustomTitleBar({ children }: { children: ReactNode }) {
    const [filterDialogOpen, setFilterDialogOpen] = useState(false);
    const { viewState, setViewState, setSearchText, setRandomGameClicked } = useSortContext(); // Access context

    const location = useLocation();
    const navigate = useNavigate();

    console.log(location.pathname);
    const handleViewClick = (view: "grid" | "list" | "details") => {
        navigate("/", { replace: true });
        console.log(`${view} View Clicked`);
        setViewState(view);
    };

    // gloabal context vars
    const { tileSize, setTileSize, setSortStateUpdate } = useSortContext();

    // This one updates only on UI
    const sizeChangeHandler = (newSize: number[]) => {
        setTileSize(newSize[0]);
    };

    // This one commits to DB and triggers on release of mouse
    const sizeChangeHandlerCommit = (newSize: number[]) => {
        setTileSize(newSize[0]);
        setSortStateUpdate(true);
    };

    return (
        <>
            <div className="flex h-screen w-screen flex-row">
                <div className="flex h-full w-14 flex-col">
                    <div className="m-auto flex h-12 w-14">
                        <QuicksaveMenu />
                    </div>
                    <div className="h-full w-14">
                        <div className="my-4 flex flex-col items-center justify-start gap-4 align-middle">
                            <Button
                                variant={"ghost"}
                                onClick={() => handleViewClick("grid")}
                                className={`group h-auto hover:bg-transparent ${viewState === "grid" && location.pathname === "/" ? "rounded-none border-r-2 border-white" : ""}`}
                            >
                                <PiBookLight
                                    className={`group-hover:scale-125 ${viewState === "grid" ? "scale-150 group-hover:scale-150" : ""}`}
                                    size={22}
                                />
                            </Button>
                            <Button
                                variant={"ghost"}
                                onClick={() => handleViewClick("details")}
                                className={`group h-auto hover:bg-transparent ${viewState === "details" && location.pathname === "/" ? "rounded-none border-r-2 border-white" : ""}`}
                            >
                                <PiListHeartLight
                                    className={`group-hover:scale-125 ${viewState === "details" ? "scale-150 group-hover:scale-150" : ""}`}
                                    size={22}
                                />
                            </Button>
                            <Button
                                variant={"ghost"}
                                onClick={() => handleViewClick("list")}
                                className={`group h-auto hover:bg-transparent ${viewState === "list" && location.pathname === "/" ? "rounded-none border-r-2 border-white" : ""}`}
                            >
                                <ListCollapse
                                    className={`group-hover:scale-125 ${viewState === "list" ? "scale-150 group-hover:scale-150" : ""}`}
                                    size={22}
                                />
                            </Button>
                        </div>
                    </div>
                </div>
                <div className="flex h-full w-full flex-col">
                    <div className="bg flex flex-row">
                        <div className="flex h-10 w-full flex-row justify-between p-1">
                            <div className="flex w-full flex-row">
                                <div className="draglayer h-full flex-1"></div>
                                <div className="relative flex h-full w-[50rem] max-w-[60vw] flex-row gap-3">
                                    <Input
                                        onChange={(e) => {
                                            setSearchText(e.target.value);
                                        }}
                                        className="my-auto h-8"
                                        placeholder="Search"
                                    />
                                    <SortGames />
                                    {filterDialogOpen && (
                                        <FilterGames
                                            filterDialogOpen={filterDialogOpen}
                                            setFilterDialogOpen={setFilterDialogOpen}
                                        />
                                    )}
                                    <Button
                                        variant={"outline"}
                                        onClick={() => {
                                            setFilterDialogOpen(!filterDialogOpen);
                                        }}
                                        className="my-auto h-8 w-8"
                                    >
                                        <Filter size={18} strokeWidth={1} />
                                    </Button>
                                    <Button
                                        onClick={() => setRandomGameClicked(true)}
                                        variant={"outline"}
                                        className="my-auto h-8 w-8"
                                    >
                                        <Dices size={18} strokeWidth={1} />
                                    </Button>
                                    <Slider
                                        className="w-80"
                                        value={[tileSize]}
                                        onValueChange={sizeChangeHandler}
                                        onValueCommit={sizeChangeHandlerCommit}
                                        step={5}
                                        min={25}
                                        max={80}
                                    />
                                </div>
                                <div className="draglayer h-full flex-1"></div>
                            </div>

                            <WindowButtons />
                        </div>
                    </div>
                    <div
                        draggable={false}
                        className="relative h-full w-full rounded-tl-xl bg-black/50"
                    >
                        {children}
                    </div>
                </div>
            </div>
        </>
    );
}

function WindowButtons() {
    return (
        <div className="flex items-start">
            <button
                title="Minimize"
                type="button"
                className="p-2 hover:bg-slate-300"
                onClick={minimizeWindow}
            >
                <svg aria-hidden="true" role="img" width="12" height="12" viewBox="0 0 12 12">
                    <rect fill="currentColor" width="10" height="1" x="1" y="6"></rect>
                </svg>
            </button>
            <button
                title="Maximize"
                type="button"
                className="p-2 hover:bg-slate-300"
                onClick={maximizeWindow}
            >
                <svg aria-hidden="true" role="img" width="12" height="12" viewBox="0 0 12 12">
                    <rect
                        width="9"
                        height="9"
                        x="1.5"
                        y="1.5"
                        fill="none"
                        stroke="currentColor"
                    ></rect>
                </svg>
            </button>
            <button
                type="button"
                title="Close"
                className="p-2 hover:bg-red-300"
                onClick={closeWindow}
            >
                <svg aria-hidden="true" role="img" width="12" height="12" viewBox="0 0 12 12">
                    <polygon
                        fill="currentColor"
                        fillRule="evenodd"
                        points="11 1.576 6.583 6 11 10.424 10.424 11 6 6.583 1.576 11 1 10.424 5.417 6 1 1.576 1.576 1 6 5.417 10.424 1"
                    ></polygon>
                </svg>
            </button>
        </div>
    );
}

function QuicksaveMenu() {
    const { setIsAddGameDialogOpen, setIsIntegrationsDialogOpen } = useSortContext();

    return (
        <div className="flex w-16 items-center justify-center">
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant={"ghost"} className="group hover:bg-transparent">
                        <BsFloppyFill size={25} className="group-hover:scale-110" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56">
                    <DropdownMenuGroup>
                        <DropdownMenuItem onClick={() => setIsAddGameDialogOpen(true)}>
                            Add a Game
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setIsIntegrationsDialogOpen(true)}>
                            Integrate Libraries
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                            Settings
                            <DropdownMenuShortcut>F4</DropdownMenuShortcut>
                        </DropdownMenuItem>
                        <DropdownMenuItem>View</DropdownMenuItem>
                    </DropdownMenuGroup>
                    <DropdownMenuSeparator />
                    <DropdownMenuGroup>
                        <DropdownMenuItem>Check For Updates</DropdownMenuItem>
                        <DropdownMenuSub>
                            <DropdownMenuSubTrigger>Theme</DropdownMenuSubTrigger>
                            <DropdownMenuPortal>
                                <DropdownMenuSubContent>
                                    <DropdownMenuItem onClick={() => darkMode()}>
                                        Dark
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => lightMode()}>
                                        Light
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => redMode()}>
                                        Red
                                    </DropdownMenuItem>
                                </DropdownMenuSubContent>
                            </DropdownMenuPortal>
                        </DropdownMenuSub>
                        <DropdownMenuItem>
                            New Team
                            <DropdownMenuShortcut>⌘+T</DropdownMenuShortcut>
                        </DropdownMenuItem>
                    </DropdownMenuGroup>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>GitHub</DropdownMenuItem>
                    <DropdownMenuItem>Support</DropdownMenuItem>
                    <DropdownMenuItem disabled>API</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                        Quit
                        <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    );
}

function SortGames() {
    const { sortOrder, setSortOrder, sortType, setSortType, setSortStateUpdate } = useSortContext(); // Access context

    const sortTypeClicked = (type: string) => {
        console.log("Sort Order Type To", type);
        setSortType(type);
        setSortStateUpdate(true);
    };

    const sortOrderClicked = (order: "ASC" | "DESC") => {
        console.log("Sort Order Set To", order);
        setSortOrder(order);
        setSortStateUpdate(true);
    };

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant={"outline"} className="f my-auto h-8 w-8">
                    <ArrowDownWideNarrow size={18} strokeWidth={1} />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56">
                <DropdownMenuLabel className="select-none">Sort Games</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                    <div className="flex select-none flex-col text-sm">
                        <div className="flex flex-col text-sm">
                            <Button
                                className={`${sortOrder === "ASC" ? "underline underline-offset-4" : ""} h-8 justify-start border-none bg-popover pl-2`}
                                variant={"outline"}
                                onClick={() => sortOrderClicked("ASC")}
                            >
                                <ChartNoAxesColumnIncreasing size={18} /> Ascending
                            </Button>
                            <Button
                                className={`${sortOrder === "DESC" ? "underline underline-offset-4" : ""} h-8 justify-start border-none bg-popover pl-2`}
                                variant={"outline"}
                                onClick={() => sortOrderClicked("DESC")}
                            >
                                <ChartNoAxesColumnDecreasing size={18} /> Descending
                            </Button>
                        </div>
                    </div>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomTitle")}
                        className={`${sortType === "CustomTitle" ? "underline underline-offset-4" : ""}`}
                    >
                        <span>Title</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomTimePlayed")}
                        className={`${sortType === "CustomTimePlayed" ? "underline underline-offset-4" : ""}`}
                    >
                        <span>Hours Played</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomRating")}
                        className={`${sortType === "CustomRating" ? "font-extrabold" : ""}`}
                    >
                        <span>Rating</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomReleaseDate")}
                        className={`${sortType === "CustomReleaseDate" ? "font-extrabold" : ""}`}
                    >
                        <span>Release Date</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomLastPlayed")}
                        className={`${sortType === "CustomLastPlayed" ? "font-extrabold" : ""}`}
                    >
                        <span>Last Played</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomCompletionStatus")}
                        className={`${sortType === "CustomCompletionStatus" ? "font-extrabold" : ""}`}
                    >
                        <span>Completion Status</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomDateAdded")}
                        className={`${sortType === "CustomDateAdded" ? "font-extrabold" : ""}`}
                    >
                        <span>Date Added</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomDeveloper")}
                        className={`${sortType === "CustomDeveloper" ? "font-extrabold" : ""}`}
                    >
                        <span>Developer</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("CustomInstallSize")}
                        className={`${sortType === "CustomInstallSize" ? "font-extrabold" : ""}`}
                    >
                        <span>Install Size</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => sortTypeClicked("OwnedPlatform")}
                        className={`${sortType === "OwnedPlatform" ? "font-extrabold" : ""}`}
                    >
                        <span>Platform</span>
                    </DropdownMenuItem>
                </DropdownMenuGroup>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

function FilterGames({ filterDialogOpen, setFilterDialogOpen }: any) {
    const [tagOptions, setTagOptions] = useState([]);
    const [platformOptions, setPlatformOptions] = useState([]);
    const [devOptions, setDevOptions] = useState([]);
    const [selectedPlatforms, setSelectedPlatforms] = useState<{ value: string; label: string }[]>(
        []
    );
    const [selectedTags, setSelectedTags] = useState<{ value: string; label: string }[]>([]);
    const [selectedDevs, setSelectedDevs] = useState<{ value: string; label: string }[]>([]);
    const [selectedName, setSelectedName] = useState<{ value: string; label: string }[]>([]);
    const [isLoaded, setIsLoaded] = useState(false); // Flag to track if data has been loaded

    const OPTIONS: Option[] = [
        { label: "A", value: "A" },
        { label: "B", value: "B" },
        { label: "C", value: "C" },
        { label: "D", value: "D" },
        { label: "E", value: "E" },
        { label: "F", value: "F" },
        { label: "G", value: "G" },
        { label: "H", value: "H" },
        { label: "I", value: "I" },
        { label: "J", value: "J" },
        { label: "K", value: "K" },
        { label: "L", value: "L" },
        { label: "M", value: "M" },
        { label: "N", value: "N" },
        { label: "O", value: "O" },
        { label: "P", value: "P" },
        { label: "Q", value: "Q" },
        { label: "R", value: "R" },
        { label: "S", value: "S" },
        { label: "T", value: "T" },
        { label: "U", value: "U" },
        { label: "V", value: "V" },
        { label: "W", value: "W" },
        { label: "X", value: "X" },
        { label: "Y", value: "Y" },
        { label: "Z", value: "Z" },
    ];

    const fetchTagsDevsPlatforms = async () => {
        try {
            const response = await fetch("http://localhost:8080/getAllTags");
            const resp = await response.json();

            // Transform the tags into key-value pairs
            const tagsAsKeyValuePairs = resp.tags.map((tag: any) => ({
                value: tag,
                label: tag,
            }));
            setTagOptions(tagsAsKeyValuePairs);
        } catch (error) {
            console.error("Error fetching tags:", error);
        }
        try {
            const response = await fetch("http://localhost:8080/getAllDevelopers");
            const resp = await response.json();
            console.log(resp);

            // Transform the tags into key-value pairs
            const devsAsKeyValuePairs = resp.devs.map((dev: any) => ({
                value: dev,
                label: dev,
            }));

            setDevOptions(devsAsKeyValuePairs);
        } catch (error) {
            console.error("Error fetching developers:", error);
        }
        try {
            const response = await fetch("http://localhost:8080/getAllPlatforms");
            const resp = await response.json();
            console.log(resp);

            // Transform the tags into key-value pairs
            const platsAsKeyValuePairs = resp.platforms.map((plat: any) => ({
                value: plat,
                label: plat,
            }));

            setPlatformOptions(platsAsKeyValuePairs);
        } catch (error) {
            console.error("Error fetching developers:", error);
        }
    };

    const handleFilterChange = async () => {
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

        try {
            // Send the filter as a POST request
            console.log("Sending Set Filter");
            const response = await fetch("http://localhost:8080/setFilter", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json", // Make sure to set the Content-Type to application/json
                },
                body: JSON.stringify(filter), // Convert the filter object to JSON
            });

            const data = await response.json();
            console.log(data); // Log the response from the server
        } catch (error) {
            console.error("Error fetching filter:", error);
        }
    };

    const clearAllFilters = async () => {
        setSelectedDevs([]);
        setSelectedName([]);
        setSelectedPlatforms([]);
        setSelectedTags([]);
        try {
            console.log("Sending Clear All Filters");
            // Send the filter as a POST request
            const response = await fetch("http://localhost:8080/clearAllFilters");
            const data = await response.json();
        } catch (error) {
            console.error("Error fetching filter:", error);
        }
    };

    const loadFilterState = async () => {
        setIsLoaded(false);
        /* setSelectedDevs([]);
        setSelectedName([]);
        setSelectedPlatforms([]);
        setSelectedTags([]); */
        try {
            console.log("Sending Load Filters");
            const response = await fetch("http://localhost:8080/LoadFilters");
            const data = await response.json();
            console.log(data);
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
            setIsLoaded(true);
        }
        setIsLoaded(true);
    };

    useEffect(() => {
        fetchTagsDevsPlatforms();
        loadFilterState();
    }, []);

    useEffect(() => {
        if (isLoaded) {
            handleFilterChange();
        }
    }, [selectedTags, selectedPlatforms, selectedDevs, selectedName]);

    return (
        <Sheet open={filterDialogOpen} onOpenChange={setFilterDialogOpen}>
            {/*             <SheetTrigger asChild>
                <Button variant={"outline"} className="my-auto w-8 h-8">
                    <Filter size={20} strokeWidth={1} />
                </Button>
            </SheetTrigger> */}
            <SheetContent>
                <SheetHeader>
                    <SheetTitle>Filter</SheetTitle>
                </SheetHeader>
                <div className="my-4 flex flex-col gap-4">
                    <div>
                        <Button className="w-full" variant={"outline"} onClick={clearAllFilters}>
                            Clear All Filter
                        </Button>
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Platform</Label>

                        <MultipleSelector
                            options={platformOptions}
                            value={selectedPlatforms}
                            onChange={(e: any) => {
                                setSelectedPlatforms(e);
                            }}
                            hidePlaceholderWhenSelected={true}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-sm">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Name</Label>

                        <MultipleSelector
                            defaultOptions={OPTIONS}
                            value={selectedName}
                            onChange={(e: any) => {
                                setSelectedName(e);
                            }}
                            hidePlaceholderWhenSelected={true}
                            placeholder="Select Name"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Tags</Label>

                        <MultipleSelector
                            options={tagOptions}
                            value={selectedTags}
                            onChange={(e: any) => {
                                setSelectedTags(e);
                            }}
                            hidePlaceholderWhenSelected={true}
                            placeholder="Select Tags"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Release Year</Label>

                        <MultipleSelector
                            defaultOptions={OPTIONS}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Developer</Label>

                        <MultipleSelector
                            options={devOptions}
                            value={selectedDevs}
                            onChange={(e: any) => {
                                setSelectedDevs(e);
                            }}
                            hidePlaceholderWhenSelected={true}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Time Played</Label>

                        <MultipleSelector
                            defaultOptions={OPTIONS}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Completion Status</Label>

                        <MultipleSelector
                            defaultOptions={OPTIONS}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="flex flex-row items-center gap-4">
                        <Label className="w-32 text-center">Installation Status</Label>

                        <MultipleSelector
                            defaultOptions={OPTIONS}
                            placeholder="Select Platforms"
                            emptyIndicator={
                                <p className="text-center text-lg leading-10">no results found.</p>
                            }
                        />
                    </div>
                    <div className="mx-2 flex flex-col gap-2">
                        <div className="flex items-center gap-2">
                            <Checkbox /> Include Hidden Games
                        </div>
                        <div className="flex items-center gap-2">
                            <Checkbox /> Favorite Only
                        </div>
                    </div>
                </div>
            </SheetContent>
        </Sheet>
    );
}
