import React, { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useRef } from "react";
import { Badge } from "@/components/ui/badge";
import { ArrowLeft, ArrowRight, CalendarDays, Clock, Settings2, Star } from "lucide-react";
import { Button } from "@/components/ui/button";
import { FaPlay } from "react-icons/fa";

import { Carousel, CarouselContent, CarouselItem } from "@/components/ui/carousel";
import { CgScrollH } from "react-icons/cg";
import { type CarouselApi } from "@/components/ui/carousel";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Checkbox } from "@/components/ui/checkbox";
import { CheckedState } from "@radix-ui/react-checkbox";
import { format } from "date-fns";
import { CalendarIcon } from "lucide-react";

import { cn } from "@/lib/utils";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function GameView() {
    const location = useLocation();
    const uid = location.state.data;
    const preloadData = location.state.preloadData;
    const hidden = location.state.hidden;
    const [companies, setCompanies] = useState("");
    const [customizeClicked, setCustomizeClicked] = useState(false);
    const [tags, setTags] = useState("");
    const [screenshots, setScreenshots] = useState("");
    const [metadata, setMetadata] = useState<any>();
    const [editDialogOpen, setEditDialogOpen] = useState<boolean>(false);
    const [hideDialogOpen, setHideDialogOpen] = useState<boolean>(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);
    const navigate = useNavigate();

    // Its on UID change to accomodate randomGamesClicked
    useEffect(() => {
        console.log("dataa", preloadData);
        setCompanies(preloadData.companies);
        setTags(preloadData.tags);
        setMetadata(preloadData.metadata);
        setScreenshots(preloadData.screenshots);
        //fetchData();
    }, [uid]);

    const fetchData = async () => {
        // try {
        //     console.log("Sending Get Game Details");
        //     const response = await fetch(`http://localhost:8080/GameDetails?uid=${uid}`);
        //     const json = await response.json();
        //     console.log(json);
        //     const { companies, tags, screenshots, m: metadata } = json.metadata;
        //     setCompanies(companies[uid]);
        //     setTags(tags[uid]);
        //     setMetadata(metadata[uid]);
        //     setScreenshots(screenshots[uid] || []); // Make sure it's an array
        // } catch (error) {
        //     console.error(error);
        // }
    };

    const tagsArray = Object.values(tags);
    const companiesArray = Object.values(companies);
    const screenshotsArray = Object.values(screenshots);
    const timePlayed = metadata?.TimePlayed?.toFixed(1);
    const isWishlist = metadata?.isDLC;
    const rating = metadata?.AggregatedRating?.toFixed(1);
    const releaseDate = metadata?.ReleaseDate;
    const [carouselIndex, setCarouselIndex] = useState(0);
    const [autoPlayOn, setAutoPlayOn] = useState(false);
    const [mainApi, setMainApi] = React.useState<CarouselApi>();
    const [thumbApi, setThumbApi] = React.useState<CarouselApi>();
    const autoPlayInterval = useRef<NodeJS.Timeout | null>(null);

    // To maintain sync B/W thumb and main carousel
    const onSelect = useCallback(() => {
        if (mainApi) {
            const selectedIndex = mainApi.selectedScrollSnap();
            setCarouselIndex(selectedIndex);
            thumbApi?.scrollTo(selectedIndex);
        }
    }, [mainApi, thumbApi]);
    useEffect(() => {
        if (mainApi) {
            mainApi.on("select", onSelect);
        }
    }, [mainApi, onSelect]);

    const handleAutoPlay = () => {
        setAutoPlayOn(!autoPlayOn);
        if (autoPlayOn) {
            if (autoPlayInterval.current) {
                clearInterval(autoPlayInterval.current);
            }
        } else {
            autoPlayInterval.current = setInterval(() => {
                mainApi?.scrollNext();
            }, 2000);
        }
    };

    const unhideGame = async () => {
        try {
            console.log("Sending Get Game Details");
            const response = await fetch(`http://localhost:8080/unhideGame?uid=${uid}`);
            const json = await response.json();
            console.log(json);
        } catch (error) {
            console.error(error);
        }
        navigate("/", { replace: true });
    };

    return (
        <>
            <img
                className="absolute z-0 h-full w-full rounded-2xl object-cover opacity-20 blur-md"
                src={"http://localhost:8080/screenshots/" + screenshots[0]}
            />

            <div className="absolute z-10 flex h-full w-full flex-col overflow-y-hidden px-6 py-8 text-center">
                <div className="mx-8 mb-2 text-left text-3xl font-semibold">{metadata?.Name}</div>
                <div className="mx-8 mb-4 flex h-full flex-row gap-10 overflow-hidden">
                    <div className="flex h-full w-1/3 flex-col overflow-y-auto">
                        <div className="mt-2 flex w-full flex-row items-center gap-4 text-base font-normal xl:flex-row">
                            <div className="flex gap-2">
                                <Button
                                    disabled={isWishlist === 0 ? false : true}
                                    className="h-10 lg:w-20 xl:w-40 2xl:w-48"
                                >
                                    <FaPlay /> Play
                                </Button>

                                <DropdownMenu>
                                    <DropdownMenuTrigger>
                                        <Button className="h-10">
                                            <Settings2 />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent>
                                        <DropdownMenuLabel>Edit Menu</DropdownMenuLabel>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem onClick={() => setEditDialogOpen(true)}>
                                            Edit Metadata
                                        </DropdownMenuItem>
                                        {hidden ? (
                                            <DropdownMenuItem onClick={unhideGame}>
                                                Unhide Game
                                            </DropdownMenuItem>
                                        ) : (
                                            <DropdownMenuItem
                                                onClick={() => setHideDialogOpen(true)}
                                            >
                                                Hide Game
                                            </DropdownMenuItem>
                                        )}
                                        <DropdownMenuItem onClick={() => setDeleteDialogOpen(true)}>
                                            Delete Game
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>

                                <EditDialog
                                    uid={uid}
                                    editDialogOpen={editDialogOpen}
                                    setEditDialogOpen={setEditDialogOpen}
                                    fetchData={fetchData}
                                />
                                <HideDialog
                                    uid={uid}
                                    hideDialogOpen={hideDialogOpen}
                                    setHideDialogOpen={setHideDialogOpen}
                                />
                                <DeleteDialog
                                    uid={uid}
                                    deleteDialogOpen={deleteDialogOpen}
                                    setDeleteDialogOpen={setDeleteDialogOpen}
                                />
                            </div>

                            <div className="flex flex-col items-center text-sm xl:ml-auto">
                                <div>
                                    <CalendarDays size={18} className="mb-1 inline" /> {releaseDate}
                                </div>
                                <div>
                                    <Star size={18} className="mb-1 inline" /> {rating}
                                    {isWishlist === 0 && (
                                        <span>
                                            <Clock className="mb-1 ml-2 inline" size={18} />
                                            {timePlayed}
                                        </span>
                                    )}
                                </div>
                            </div>
                        </div>

                        <DisplayInfo data={metadata} tags={tagsArray} companies={companiesArray} />
                    </div>
                    <div className="flex h-full w-2/3 flex-col items-center overflow-x-hidden">
                        <Carousel
                            opts={{ loop: true, dragFree: true }}
                            setApi={setMainApi}
                            onWheel={(e) => {
                                if (e.deltaY > 0) {
                                    mainApi?.scrollNext();
                                } else {
                                    mainApi?.scrollPrev();
                                }
                            }}
                            className="aspect-video max-h-[75%] overflow-hidden"
                        >
                            <CarouselContent className="max-h-full">
                                {screenshotsArray.map((_, index) => (
                                    <CarouselItem
                                        key={index}
                                        className="flex max-h-full justify-center"
                                    >
                                        <img
                                            src={
                                                "http://localhost:8080/screenshots" +
                                                screenshotsArray[index]
                                            }
                                            className="max-h-full max-w-full cursor-pointer rounded-xl object-contain"
                                        />
                                    </CarouselItem>
                                ))}
                            </CarouselContent>
                        </Carousel>
                        <div className="flex w-full flex-col items-center justify-center overflow-y-clip">
                            <Carousel
                                opts={{ dragFree: true }}
                                setApi={setThumbApi}
                                className="items-center"
                                onWheel={(e) => {
                                    if (e.deltaY > 0) {
                                        thumbApi?.scrollNext();
                                    } else {
                                        thumbApi?.scrollPrev();
                                    }
                                }}
                            >
                                <CarouselContent className="my-2 px-1">
                                    {screenshotsArray.map((_, index) => (
                                        <CarouselItem
                                            onClick={() => {
                                                mainApi?.scrollTo(index);
                                                setCarouselIndex(index);
                                            }}
                                            key={index}
                                            className="basis-32"
                                        >
                                            <img
                                                src={
                                                    "http://localhost:8080/screenshots" +
                                                    screenshotsArray[index]
                                                }
                                                className={`cursor-pointer rounded-md ${index === carouselIndex ? "ring-2 ring-border" : null}`}
                                            />
                                        </CarouselItem>
                                    ))}
                                </CarouselContent>
                            </Carousel>
                            <div className="flex gap-2">
                                <Button
                                    className="h-8 w-8 rounded-full"
                                    onClick={() => mainApi?.scrollPrev()}
                                    variant={"outline"}
                                >
                                    <ArrowLeft size={18} />
                                </Button>
                                <Button
                                    onClick={handleAutoPlay}
                                    variant={"outline"}
                                    className={`h-8 w-8 rounded-full ${autoPlayOn ? "animate-spin duration-1000" : null}`}
                                >
                                    <CgScrollH size={22} />
                                </Button>
                                <Button
                                    onClick={() => mainApi?.scrollNext()}
                                    className="h-8 w-8 rounded-full"
                                    variant={"outline"}
                                >
                                    <ArrowRight size={18} />
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}

function EditDialog({ uid, editDialogOpen, setEditDialogOpen, fetchData }: any) {
    const [customTitleChecked, setCustomTitleChecked] = useState<CheckedState | undefined>(false);
    const [customTitle, setCustomTitle] = useState("");
    const [customTime, setCustomTime] = useState("");
    const [customTimeOffset, setCustomTimeOffset] = useState("");
    const [customTimeChecked, setCustomTimeChecked] = useState<CheckedState | undefined>(false);
    const [customTimeOffsetChecked, setCustomTimeOffsetChecked] = useState<
        CheckedState | undefined
    >(false);
    const [customReleaseDate, setCustomReleaseDate] = useState<Date>();
    const [customReleaseDateChecked, setCustomReleaseDateChecked] = useState<
        CheckedState | undefined
    >(false);
    const [customRating, setCustomRating] = useState("");
    const [customRatingChecked, setCustomRatingChecked] = useState<CheckedState | undefined>(false);

    const loadPreferences = async () => {
        try {
            const response = await fetch(`http://localhost:8080/LoadPreferences?uid=${uid}`);
            const json = await response.json();
            console.log(json);
            setCustomTime("0");
            setCustomTimeOffset("0");
            setCustomRating("0");
            setCustomTitle(json.preferences.title.value);
            if (json.preferences.time.value) {
                setCustomTime(json.preferences.time.value);
            }
            if (json.preferences.timeOffset.value) {
                setCustomTimeOffset(json.preferences.timeOffset.value);
            }
            if (json.preferences.rating.value) {
                console.log("A");
                setCustomRating(json.preferences.rating.value);
            }
            setCustomReleaseDate(json.preferences.releaseDate.value);

            if (json.preferences.title.checked == "1") {
                setCustomTitleChecked(true);
            }
            if (json.preferences.time.checked == "1") {
                setCustomTimeChecked(true);
            }
            if (json.preferences.timeOffset.checked == "1") {
                setCustomTimeOffsetChecked(true);
            }
            if (json.preferences.releaseDate.checked == "1") {
                setCustomReleaseDateChecked(true);
            }
            if (json.preferences.rating.checked == "1") {
                setCustomRatingChecked(true);
            }
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        loadPreferences();
    }, []);

    const handleCheckboxChange = (checked: CheckedState | undefined, title: string) => {
        switch (title) {
            case "title":
                setCustomTitleChecked(checked);
                break;
            case "time":
                if (checked == false) {
                    setCustomTimeChecked(false);
                } else {
                    setCustomTimeChecked(true);
                    setCustomTimeOffsetChecked(false);
                }
                break;

            case "timeOffset":
                if (checked == false) {
                    setCustomTimeOffsetChecked(false);
                } else {
                    setCustomTimeOffsetChecked(true);
                    setCustomTimeChecked(false);
                }
                break;
            case "releaseDate":
                setCustomReleaseDateChecked(checked);
                break;
            case "rating":
                setCustomRatingChecked(checked);
                break;
            default:
                break;
        }
    };

    const saveClickHandler = () => {
        const postData = {
            customTitleChecked: customTitleChecked,
            customTitle: customTitle.trim(),
            customTimeChecked: customTimeChecked,
            customTime: customTime.trim() === "" ? "0" : customTime.trim(),
            customTimeOffsetChecked: customTimeOffsetChecked,
            customTimeOffset: customTimeOffset.trim() === "" ? "0" : customTimeOffset.trim(),
            customRatingChecked: customRatingChecked,
            customRating: customRating.trim() === "" ? "0" : customRating.trim(),
            customReleaseDateChecked: customReleaseDateChecked,
            customReleaseDate: customReleaseDate ? format(customReleaseDate, "yyyy-MM-dd") : "", // Empty string if undefined
            UID: uid,
        };
        savePreferences(postData);
    };

    const savePreferences = async (postData: any) => {
        try {
            const response = await fetch(`http://localhost:8080/SavePreferences`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postData),
            });
            const json = await response.json();
        } catch (error) {
            console.error(error);
        }
        fetchData();
    };

    return (
        <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
            <DialogContent className="h-[75vh] max-h-[75vh] max-w-[75vw]">
                <DialogHeader className="h-full">
                    <DialogTitle>Customize Game</DialogTitle>
                    <DialogDescription>
                        Changes made here will reflect in the rest of the app. This can be reverted
                        at any time.
                    </DialogDescription>
                    <Tabs defaultValue="metadata" className="flex h-full flex-col">
                        <TabsList className="grid w-[300px] grid-cols-3 focus:outline-none">
                            <TabsTrigger value="metadata">Metadata</TabsTrigger>
                            <TabsTrigger value="images">Images</TabsTrigger>
                            <TabsTrigger value="delete">Delete</TabsTrigger>
                        </TabsList>
                        <TabsContent value="metadata" className="h-full focus:ring-0">
                            <div className="flex h-full flex-col justify-between p-2 px-4 focus:outline-none">
                                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                    <div className="flex items-center">
                                        <Label className="w-60">Use Custom Title</Label>
                                        <Checkbox
                                            checked={customTitleChecked}
                                            onCheckedChange={(checked) =>
                                                handleCheckboxChange(checked, "title")
                                            }
                                            className="mr-10"
                                        />
                                        <Input
                                            disabled={!customTitleChecked}
                                            id="customTitle"
                                            value={customTitle}
                                            onChange={(e) => setCustomTitle(e.target.value)}
                                        ></Input>
                                    </div>
                                    <div className="flex items-center">
                                        <Label className="w-60">Custom Release Date</Label>
                                        <Checkbox
                                            checked={customReleaseDateChecked}
                                            onCheckedChange={(checked) =>
                                                handleCheckboxChange(checked, "releaseDate")
                                            }
                                            className="mr-10"
                                        />
                                        <Popover>
                                            <PopoverTrigger asChild>
                                                <Button
                                                    variant={"outline"}
                                                    disabled={!customReleaseDateChecked}
                                                    className={cn(
                                                        "w-full justify-start text-left font-normal",
                                                        !customReleaseDate &&
                                                            "text-muted-foreground"
                                                    )}
                                                >
                                                    <CalendarIcon size={18} />
                                                    {customReleaseDate ? (
                                                        format(customReleaseDate, "yyyy-MM-dd")
                                                    ) : (
                                                        <span>Pick a date</span>
                                                    )}
                                                </Button>
                                            </PopoverTrigger>
                                            <PopoverContent className="w-auto p-0">
                                                <Calendar
                                                    disabled={!customReleaseDateChecked}
                                                    mode="single"
                                                    selected={customReleaseDate}
                                                    onSelect={setCustomReleaseDate}
                                                    initialFocus
                                                />
                                            </PopoverContent>
                                        </Popover>
                                    </div>
                                    <div className="flex items-center">
                                        <Label className="w-60">Custom Time Played</Label>
                                        <Checkbox
                                            checked={customTimeChecked}
                                            onCheckedChange={(checked) =>
                                                handleCheckboxChange(checked, "time")
                                            }
                                            className="mr-10"
                                        />
                                        <Input
                                            disabled={!customTimeChecked}
                                            id="customTime"
                                            value={customTime}
                                            onChange={(e) => {
                                                const value = e.target.value;
                                                // Restrict input to only numbers and up to 2 decimal places
                                                if (/^\d*\.?\d{0,2}$/.test(value)) {
                                                    setCustomTime(value);
                                                }
                                            }}
                                        />
                                    </div>
                                    <div className="flex items-center">
                                        <Label className="w-60">Custom Time Offset</Label>
                                        <Checkbox
                                            checked={customTimeOffsetChecked}
                                            onCheckedChange={(checked) =>
                                                handleCheckboxChange(checked, "timeOffset")
                                            }
                                            className="mr-10"
                                        />
                                        <Input
                                            disabled={!customTimeOffsetChecked}
                                            id="customTimeOffset"
                                            value={customTimeOffset}
                                            onChange={(e) => {
                                                const value = e.target.value;
                                                if (/^-?\d*\.?\d{0,2}$/.test(value)) {
                                                    setCustomTimeOffset(value);
                                                }
                                            }}
                                        />
                                    </div>
                                    <div className="flex items-center">
                                        <Label className="w-60">Custom Rating</Label>
                                        <Checkbox
                                            checked={customRatingChecked}
                                            onCheckedChange={(checked) =>
                                                handleCheckboxChange(checked, "rating")
                                            }
                                            className="mr-10"
                                        />
                                        <Input
                                            disabled={!customRatingChecked}
                                            id="customRating"
                                            value={customRating}
                                            onChange={(e) => {
                                                const value = e.target.value;
                                                // Restrict input to only numbers and up to 2 decimal places
                                                if (/^\d*\.?\d{0,2}$/.test(value)) {
                                                    setCustomRating(value);
                                                }
                                            }}
                                        />
                                    </div>
                                </div>
                                <div className="self-end">
                                    <Button
                                        onClick={saveClickHandler}
                                        variant={"secondary"}
                                        className="h-12 w-60"
                                    >
                                        Save
                                    </Button>
                                </div>
                            </div>
                        </TabsContent>
                        <TabsContent value="images">
                            <div>xyz</div>
                        </TabsContent>
                    </Tabs>
                </DialogHeader>
            </DialogContent>
        </Dialog>
    );
}

function HideDialog({ uid, hideDialogOpen, setHideDialogOpen }: any) {
    const navigate = useNavigate();

    const hide = async () => {
        console.log("Sending Hide Game");
        try {
            const response = await fetch(`http://localhost:8080/HideGame?uid=${uid}`);
            const json = await response.json();
        } catch (error) {
            console.error(error);
        }
        navigate("/", { replace: true });
    };

    const hardDelete = async () => {
        console.log("Sending Delete Game");
        try {
            const response = await fetch(`http://localhost:8080/DeleteGame?uid=${uid}`);
            const json = await response.json();
        } catch (error) {
            console.error(error);
        }
        navigate("/", { replace: true });
    };

    return (
        <Dialog open={hideDialogOpen} onOpenChange={setHideDialogOpen}>
            <DialogContent className="h-[600px] max-h-[300px] max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Hide Game</DialogTitle>
                    <DialogDescription>
                        Hidden games can be viewed and reverted at any time. Custom metadata is
                        saved and these games will not be re-imported on a library integration
                        update.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter className="flex items-end">
                    <Button variant={"secondary"} onClick={hide} className="h-12 w-32">
                        Hide Game
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

function DeleteDialog({ uid, deleteDialogOpen, setDeleteDialogOpen }: any) {
    const navigate = useNavigate();

    const hardDelete = async () => {
        console.log("Sending Delete Game");
        try {
            const response = await fetch(`http://localhost:8080/DeleteGame?uid=${uid}`);
            const json = await response.json();
        } catch (error) {
            console.error(error);
        }
        navigate("/", { replace: true });
    };

    return (
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
            <DialogContent className="h-[600px] max-h-[300px] max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Delete Game</DialogTitle>
                    <DialogDescription>
                        The game will be permanently deleted from the app. The game will be
                        re-imported on library synchronization, if you wish to stop this behaviour,
                        hide the game instead.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter className="flex items-end">
                    <Button variant={"destructive"} onClick={hardDelete} className="h-12 w-32">
                        Delete Game
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

function DisplayInfo({ data, tags, companies }: any) {
    return (
        <div className="mt-2 flex h-full flex-col gap-4 overflow-y-auto pr-1 text-left">
            <div className="flex flex-row items-center justify-start">
                <p className="flex flex-col items-start gap-2 text-base font-medium">
                    Platform
                    <Button className="h-6 rounded-full">{data?.OwnedPlatform}</Button>
                </p>
            </div>
            <div className="flex flex-col gap-2 text-base">
                <p className="text-left text-base">Tags</p>
                <div className="flex flex-wrap gap-2 rounded-md text-center">
                    {tags.map((items: any, index: any) => (
                        <Badge key={index}>{items}</Badge>
                    ))}
                </div>
            </div>

            <div className="flex flex-col items-start justify-center gap-2 text-base">
                <p>Developers And Publishers</p>
                <div className="flex flex-wrap gap-2 rounded-md text-center">
                    {companies.map((items: any, index: any) => (
                        <Badge draggable={false} key={index}>
                            {items}
                        </Badge>
                    ))}
                </div>
            </div>

            <div className="flex flex-col items-start justify-start gap-2">
                Description
                <div className="flex h-full flex-col">
                    <p
                        dangerouslySetInnerHTML={{ __html: data?.Description }}
                        className="text-sm"
                    ></p>
                </div>
            </div>
        </div>
    );
}
