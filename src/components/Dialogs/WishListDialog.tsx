import React, { useEffect, useState } from "react";

import { useToast } from "@/hooks/use-toast";
import {
  fetchTagsDevsPlatforms,
  sendGameToDB,
} from "@/lib/api/addGameManuallyAPI";
import {
  CalendarIcon,
  Globe,
  Link,
  Loader2,
  LucideArrowLeft,
  LucideArrowRight,
  Plus,
  Trash2,
} from "lucide-react";

import { format } from "date-fns";

import { cn } from "@/lib/utils";
import { useSortContext } from "@/hooks/useSortContex";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Carousel,
  CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import MultipleSelector from "@/components/ui/multiple-selector";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { ScrollArea } from "@/components/ui/scroll-area";
import { searchGame } from "@/lib/api/addGameManuallyAPI";
import { DateTimePicker } from "../ui/datetime-picker";
import ImageSearchDialog from "./ImageSearchDialog";

export default function WishlistDialog() {
  const { isWishlistAddDialogOpen, setIsWishlistAddDialogOpen } =
    useSortContext();

  return (
    <>
      {isWishlistAddDialogOpen && (
        <Dialog
          open={isWishlistAddDialogOpen}
          onOpenChange={setIsWishlistAddDialogOpen}
        >
          <DialogContent className="block h-[75vh] max-h-[75vh] max-w-[75vw]">
            <DialogHeader className="h-full max-h-full">
              <DialogTitle>Add to Wishlist</DialogTitle>
              <MetaDataView />
            </DialogHeader>
            <DialogFooter></DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

function MetaDataView() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [title, setTitle] = useState<string>("");
  const [releaseDate, setReleaseDate] = useState<any>("");
  const [rating, setRating] = useState<any>("");
  const [developers, setDevelopers] = useState<any>("");
  const [description, setDescription] = useState<any>("");
  const [tagOptions, setTagOptions] = useState([]);
  const [devOptions, setDevOptions] = useState([]);
  const [platformOptions, setPlatformOptions] = useState([]);
  const [selectedTags, setSelectedTags] = useState([]);
  const [selectedDevs, setSelectedDevs] = useState([]);
  const [selectedPlatforms, setSelectedPlatforms] = useState<any[]>([]);
  const [titleEmpty, setTitleEmpty] = useState(false);
  const [releaseDateEmpty, setReleaseDateEmpty] = useState(false);
  const [platformEmpty, setPlatformEmpty] = useState(false);
  const [addGameLoading, setAddGameLoading] = useState(false);

  useEffect(() => {
    fetchTagsDevsPlatforms(setTagOptions, setDevOptions, setPlatformOptions);
  }, []);

  const SearchGameClicked = () => {
    searchGame(title, setTitleEmpty, setLoading, setData, toast);
  };

  const [coverImage, setCoverImage] = useState<string | null>(null);
  const [ssImage, setSsImage] = useState<(string | null)[]>([null]); // Three empty image slots
  const [selectedIndex, setSelectedIndex] = useState<number | null>(0); // To track selected carousel item
  const [coverArtLinkClicked, setCoverArtLinkClicked] =
    useState<boolean>(false); // To track selected carousel item
  const [ssLinkClicked, setSSLinkClicked] = useState<number | null>(null);

  const handleCoverImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        // Set the base64 result as the image source
        setCoverImage(reader.result as string); // Cast to 'string' since 'result' can be a string or null
      };
      reader.readAsDataURL(file);
    } else {
      // If no file is selected (user cancels), clear the image
      setCoverImage(null);
    }
  };

  const handleScreenshotImageChange = (
    index: number,
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const updatedScreenshots = [...(ssImage || [])]; // Create a new array from ssImage or empty array
        updatedScreenshots[index] = reader.result as string; // Ensure the result is a string for the data URL
        console.log("Updated Screenshot Array:", updatedScreenshots); // Add a log to debug the updated state
        setSsImage(updatedScreenshots); // Update the state with the new array
      };
      reader.readAsDataURL(file);
    }
  };

  // Add a new screenshot entry to the carousel
  const addScreenshot = () => {
    // First, update the state to add a new screenshot slot
    setSsImage((prev) => {
      const updatedScreenshots = [...(prev || []), null]; // Add a new `null` entry
      return updatedScreenshots; // Return the updated screenshots array
    });

    // Use setTimeout to delay the scrolling, ensuring it happens after the state update
    setTimeout(() => {
      const lastIndex = ssImage.length; // Using the state length as the new index
      if (api) {
        api.scrollTo(lastIndex); // Scroll to the last item
      }
    }, 0); // 0ms delay ensures it happens after the state is updated
  };

  const handleDeleteScreenshot = () => {
    // Ensure selectedIndex is a valid number, falling back to 0 if it's null or undefined
    const indexToDelete = selectedIndex ?? 0;

    const updatedScreenshots = [...(ssImage || [])];

    // Only proceed with deletion if there are screenshots to delete
    if (updatedScreenshots.length > 0) {
      api?.scrollPrev();
      setTimeout(() => {
        updatedScreenshots.splice(indexToDelete, 1); // Remove the screenshot at the active index
        setSsImage(updatedScreenshots);
        if (selectedIndex) {
          setSelectedIndex(selectedIndex - 1);
        }
      }, 500); // 500ms delay to see the full scroll before slide is deleted
    }
  };
  const [api, setApi] = React.useState<CarouselApi>();

  useEffect(() => {
    if (!api) {
      return;
    }

    api.on("select", (e) => {
      setSelectedIndex(e.selectedScrollSnap());
    });
  }, [api]);

  const addGameClickHandler = () => {
    const isTitleEmpty = title.trim() === "";
    const isReleaseDateEmpty = !releaseDate;
    const isPlatformEmpty = selectedPlatforms.length === 0;
    setTitleEmpty(isTitleEmpty);
    setReleaseDateEmpty(isReleaseDateEmpty);
    setPlatformEmpty(isPlatformEmpty);

    if (isTitleEmpty || isReleaseDateEmpty || isPlatformEmpty) {
      return;
    }

    const ratingNormal =
      rating === "" || isNaN(Number(rating)) ? "0" : String(rating);

    sendGameToDB(
      title,
      releaseDate,
      selectedPlatforms,
      "0", // For 0 timePlayed
      ratingNormal,
      selectedDevs,
      selectedTags,
      description,
      coverImage,
      ssImage,
      1, // For wishlist
      setAddGameLoading,
      toast
    );
  };

  const { toast } = useToast();

  const [searchDialogOpen, setSearchDialogOpen] = useState(false);
  const [selectingForCover, setSelectingForCover] = useState(false);

  const searchCoverClickHandler = () => {
    setSelectingForCover(true);
    setSearchDialogOpen(true);
  };

  const searchScreenshotClickHandler = () => {
    setSelectingForCover(false);
    setSearchDialogOpen(true);
  };

  return (
    <>
      {searchDialogOpen && (
        <ImageSearchDialog
          searchDialogOpen={searchDialogOpen}
          setSearchDialogOpen={setSearchDialogOpen}
          title={title}
          defaultSearchSuffix={selectingForCover ? "cover" : "1080p"}
          onImageSelect={(url) => {
            if (selectingForCover) {
              setCoverImage(url); // Set as cover image
            } else if (selectedIndex !== null) {
              const updatedScreenshots = [...ssImage];
              updatedScreenshots[selectedIndex] = url;
              setSsImage(updatedScreenshots); // Set as screenshot
            }
          }}
        />
      )}
      <div className="flex overflow-hidden gap-4 w-full h-full max-h-full">
        <div className="flex flex-col w-full h-full">
          <DialogDescription
            className={`${
              titleEmpty || platformEmpty || releaseDateEmpty
                ? "text-destructive"
                : null
            }`}
          >
            {titleEmpty || platformEmpty || releaseDateEmpty
              ? "Please fill all required fields"
              : "  Enter the game metadata manually or download it from IGDB."}
          </DialogDescription>
          <div className="grid overflow-y-auto grid-cols-1 gap-4 p-2 py-4 mb-2 2xl:grid-cols-2">
            <div className="flex gap-4 items-center text-sm">
              <label
                className={`min-w-24 ${titleEmpty ? "text-destructive" : null}`}
              >
                {titleEmpty && "*"}
                Title
              </label>
              <Input
                id="title"
                value={title}
                onChange={(e) => {
                  setTitle(e.target.value);
                  setTitleEmpty(false);
                }}
                placeholder={"Enter Title"}
                className="col-span-3"
                spellCheck={false}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    SearchGameClicked();
                  }
                }}
              />
              <Button
                disabled={loading}
                type="submit"
                variant={"dialogSaveButton"}
                onClick={SearchGameClicked}
                className="h-10"
              >
                {loading && <Loader2 className="animate-spin" />}
                Download Metadata
              </Button>
            </div>

            <div className="flex gap-4 items-center text-sm">
              <label
                className={`min-w-24 ${
                  releaseDateEmpty ? "text-destructive" : null
                }`}
              >
                {titleEmpty && "*"}
                Release Date
              </label>
              <DateTimePicker
                hideTime={true}
                value={releaseDate}
                onChange={setReleaseDate}
              />
            </div>
            <div className="flex gap-4 items-center text-sm">
              <label
                className={`min-w-24 ${
                  platformEmpty ? "text-destructive" : null
                }`}
              >
                {platformEmpty && "*"}
                Platform
              </label>
              <MultipleSelector
                maxSelected={1}
                options={platformOptions}
                placeholder="Select Platforms"
                creatable
                className="overflow-y-scroll max-h-40"
                hidePlaceholderWhenSelected={true}
                value={selectedPlatforms}
                onChange={(e: any) => {
                  setSelectedPlatforms(e);
                  setPlatformEmpty(false);
                }}
                emptyIndicator={
                  <p className="text-sm text-center">no results found.</p>
                }
              />
            </div>
            <div className="flex gap-4 items-center text-sm">
              <label className="min-w-24">Rating</label>
              <Input
                type="text"
                value={rating}
                onChange={(e) => {
                  const value = e.target.value;

                  // If the value is 100, don't allow a decimal point to be added
                  if (value === "100") {
                    setRating(value); // Keep it as "100" without a decimal point
                  } else {
                    // Allow only digits with at most two decimal places and ensure the value is between 0 and 100
                    if (
                      /^\d*\.?\d{0,2}$/.test(value) &&
                      parseFloat(value) <= 100 &&
                      parseFloat(value) >= 0
                    ) {
                      setRating(value);
                    } else if (value === "") {
                      setRating(""); // Allow clearing the input
                    }
                  }
                }}
                placeholder="Rating"
                className="col-span-3"
                inputMode="decimal" // Use the numeric keypad with decimal point on mobile devices
              />
            </div>
            <div className="flex gap-4 items-center text-sm">
              <label className="min-w-24">Developers</label>
              <MultipleSelector
                options={devOptions}
                placeholder="Select Developers"
                creatable
                className="overflow-y-scroll max-h-40"
                hidePlaceholderWhenSelected={true}
                value={selectedDevs}
                onChange={(e: any) => {
                  setSelectedDevs(e);
                }}
                emptyIndicator={
                  <p className="text-sm text-center">no results found.</p>
                }
              />
            </div>
            <div className="flex gap-4 items-start w-full h-full text-sm k grow-0">
              <label className="min-w-24">Tags</label>
              <div className="flex-1 h-full max-h-24">
                <MultipleSelector
                  options={tagOptions}
                  placeholder="Select Tags"
                  creatable
                  className="overflow-y-scroll w-full h-full max-h-full"
                  hidePlaceholderWhenSelected={true}
                  onChange={(e: any) => {
                    setSelectedTags(e);
                  }}
                  value={selectedTags}
                  emptyIndicator={
                    <p className="text-sm text-center">no results found.</p>
                  }
                />
              </div>
            </div>

            <div className="flex gap-4 items-start text-sm 2xl:col-span-2">
              <label className="mt-2 min-w-24">Description</label>
              <Textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Description"
                spellCheck={false}
              ></Textarea>
            </div>

            <div className="flex gap-4 items-start text-sm">
              <label className="mt-2 min-w-24">Cover Art</label>
              <div className="flex gap-2 w-full">
                <input
                  hidden
                  id="cover-picture"
                  type="file"
                  onChange={handleCoverImageChange}
                />

                <label htmlFor="cover-picture" className="w-60 cursor-pointer">
                  <Card className="flex h-[calc(15rem*4/3)] w-60 select-none items-center justify-center overflow-hidden">
                    <CardContent className="p-0 m-0 h-full w-full flex">
                      {coverImage ? (
                        <img
                          draggable={false}
                          src={coverImage}
                          className="object-cover w-full h-full"
                        />
                      ) : (
                        <div className="text-sm m-auto text-muted-foreground">
                          Choose an Image
                        </div>
                      )}
                    </CardContent>
                  </Card>
                </label>
                <div className="flex flex-col gap-2">
                  <Button
                    variant={"outline"}
                    className="w-8 h-8 rounded-full"
                    onClick={() => {
                      setCoverImage(null);
                    }}
                  >
                    <Trash2 size={18} />
                  </Button>
                  <div className="flex gap-1">
                    <Button
                      variant={"outline"}
                      onClick={() => {
                        setCoverArtLinkClicked(!coverArtLinkClicked);
                      }}
                      className="w-8 h-8 rounded-full"
                    >
                      <Link size={18} />
                    </Button>
                    {coverArtLinkClicked && (
                      <Input
                        placeholder="Paste link and press Enter"
                        className="h-8 rounded-full"
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            const inputElement = e.target as HTMLInputElement; // Type cast to HTMLInputElement
                            setCoverImage(inputElement.value); // Access value
                          }
                        }}
                      ></Input>
                    )}
                  </div>
                  <Button
                    variant={"outline"}
                    className="w-8 h-8 rounded-full"
                    onClick={searchCoverClickHandler}
                  >
                    <Globe size={18} />
                  </Button>
                </div>
              </div>
            </div>

            <div className="flex gap-4 items-start text-sm">
              <label className="mt-2 min-w-24">Screenshots</label>
              <div className="flex flex-col w-full">
                <Carousel
                  setApi={setApi}
                  className="flex justify-center max-w-lg"
                >
                  <CarouselContent className="w-full h-full">
                    {ssImage?.map((image, index) => (
                      <CarouselItem key={index} className="w-full h-full">
                        <input
                          hidden
                          id={`ss-picture-${index}`} // Unique ID for each screenshot input
                          type="file"
                          onChange={(e) =>
                            handleScreenshotImageChange(index, e)
                          }
                        />
                        <div
                          className="w-full h-full cursor-pointer"
                          onClick={() => {
                            const fileInput = document.getElementById(
                              `ss-picture-${index}`
                            ) as HTMLInputElement | null;
                            if (fileInput) {
                              fileInput.click();
                            }
                          }}
                        >
                          <Card className="w-full h-full">
                            <CardContent className="flex justify-center items-center p-0 w-full h-full select-none">
                              <AspectRatio
                                ratio={16 / 9}
                                className="flex justify-center items-center rounded-md border-b border-border"
                              >
                                {image ? (
                                  <img
                                    draggable={false}
                                    src={image}
                                    alt="Broken Image"
                                    className="h-full rounded-md"
                                  />
                                ) : (
                                  <div className="text-muted-foreground">
                                    Click to choose image
                                  </div>
                                )}
                              </AspectRatio>
                            </CardContent>
                          </Card>
                        </div>
                      </CarouselItem>
                    ))}
                  </CarouselContent>
                  <div className="flex gap-2 justify-between mt-1">
                    <div className="flex gap-2">
                      <Button
                        variant={"outline"}
                        onClick={addScreenshot}
                        className="w-8 h-8 rounded-full"
                      >
                        <Plus size={18} />
                      </Button>
                      <Button
                        variant={"outline"}
                        onClick={handleDeleteScreenshot}
                        className="w-8 h-8 rounded-full"
                      >
                        <Trash2 size={18} />
                      </Button>
                      <Button
                        variant={"outline"}
                        onClick={searchScreenshotClickHandler}
                        className="w-8 h-8 rounded-full"
                      >
                        <Globe size={18} />
                      </Button>
                      <div className="flex gap-1">
                        <Button
                          variant={"outline"}
                          onClick={() => {
                            const scrollpoint = api?.selectedScrollSnap();
                            if (scrollpoint != null) {
                              setSSLinkClicked(scrollpoint);
                              if (scrollpoint == ssLinkClicked) {
                                setSSLinkClicked(null);
                              }
                            }
                          }}
                          className="w-8 h-8 rounded-full"
                        >
                          <Link size={18} />
                        </Button>
                        {ssLinkClicked != null ? (
                          <Input
                            className="mx-2 h-8 rounded-full"
                            placeholder="Paste link and press Enter"
                            onKeyDown={(e) => {
                              if (e.key === "Enter" && selectedIndex !== null) {
                                const link = (e.target as HTMLInputElement)
                                  .value; // Typecast to HTMLInputElement

                                // Ensure it's a valid link before updating
                                const updatedScreenshots = [...ssImage];
                                updatedScreenshots[selectedIndex] = link; // Set the link at the selected index
                                console.log(updatedScreenshots);
                                setSsImage(updatedScreenshots); // Update the state with the new array
                                setSSLinkClicked(null);
                              }
                            }}
                          />
                        ) : null}
                      </div>
                    </div>
                    <div className="flex gap-2 mr-4">
                      <Button
                        variant={"outline"}
                        onClick={() => {
                          api?.scrollPrev();
                          setSSLinkClicked(null);
                        }}
                        className="w-8 h-8 rounded-full"
                      >
                        <LucideArrowLeft size={18} />
                      </Button>
                      <Button
                        variant={"outline"}
                        onClick={() => {
                          api?.scrollNext();
                          setSSLinkClicked(null);
                        }}
                        className="w-8 h-8 rounded-full"
                      >
                        <LucideArrowRight size={18} />
                      </Button>
                    </div>
                  </div>
                </Carousel>
              </div>
            </div>
          </div>
          <div className="flex justify-end mt-auto">
            <Button
              type="submit"
              onClick={addGameClickHandler}
              variant={"dialogSaveButton"}
            >
              Add Game {addGameLoading && <Loader2 className="animate-spin" />}
            </Button>
          </div>
        </div>
        <div className="flex h-full max-h-full">
          <FoundGames
            data={data}
            setData={setData}
            title={title}
            releaseDate={releaseDate}
            rating={rating}
            developers={developers}
            description={description}
            setTitle={setTitle}
            setReleaseDate={setReleaseDate}
            setRating={setRating}
            setDevelopers={setDevelopers}
            setDescription={setDescription}
            tagOptions={tagOptions}
            setTagOptions={setTagOptions}
            selectedTags={selectedTags}
            setSelectedTags={setSelectedTags}
            selectedDevs={selectedDevs}
            setSelectedDevs={setSelectedDevs}
            setCoverImage={setCoverImage}
            setSsImage={setSsImage}
          />
        </div>
      </div>
    </>
  );
}

function FoundGames({
  data,
  setData,
  title,
  setTitle,
  releaseDate,
  setReleaseDate,
  rating,
  setRating,
  developers,
  setDevelopers,
  description,
  setDescription,
  tagOptions,
  setTagOptions,
  selectedTags,
  setSelectedTags,
  selectedDevs,
  setSelectedDevs,
  setCoverImage,
  setSsImage,
}: {
  data: any;
  setData: React.Dispatch<React.SetStateAction<string | null>>;
  title: string;
  setTitle: React.Dispatch<React.SetStateAction<string>>;
  releaseDate: any;
  setReleaseDate: React.Dispatch<React.SetStateAction<any>>;
  rating: any;
  setRating: React.Dispatch<React.SetStateAction<any>>;
  developers: any;
  setDevelopers: React.Dispatch<React.SetStateAction<any>>;
  description: any;
  setDescription: React.Dispatch<React.SetStateAction<any>>;
  tagOptions: { value: string; label: string }[];
  setTagOptions: React.Dispatch<React.SetStateAction<any>>;
  selectedTags: { value: string; label: string }[];
  setSelectedTags: React.Dispatch<React.SetStateAction<any>>;
  selectedDevs: { value: string; label: string }[];
  setSelectedDevs: React.Dispatch<React.SetStateAction<any>>;
  setCoverImage: React.Dispatch<React.SetStateAction<any>>;
  setSsImage: React.Dispatch<React.SetStateAction<any>>;
}) {
  const [gameInfoLoading, setGameInfoLoading] = useState(false);
  const [loadingAppId, setLoadingAppId] = useState<string | null>(null);

  const IgdbGameClicked = async (appid: any) => {
    try {
      setGameInfoLoading(true);
      setLoadingAppId(appid);
      const response = await fetch("http://localhost:50001/GetIgdbInfo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          key: appid,
        }),
      });
      const data = await response.json();
      console.log(data);

      // Overrides prev set tags and devs
      const newTags = data.metadata.tags;
      const newDevs = data.metadata.involvedCompanies;

      const selectedTags = newTags.map((tag: string) => ({
        value: tag,
        label: tag,
      }));

      const selectedDevs = newDevs.map((dev: string) => ({
        value: dev,
        label: dev,
      }));

      // Update the tag options state with the newly selected tags
      setSelectedTags(selectedTags);
      setSelectedDevs(selectedDevs);

      // Set the other metadata states
      setTitle(data.metadata.name);
      setReleaseDate(data.metadata.releaseDate);
      setRating(Number(data.metadata.aggregatedRating || 0).toFixed(2));
      setDescription(data.metadata.description);
      setCoverImage(data.metadata.cover);
      setSsImage(data.metadata.screenshots);
      setData(null);
      setGameInfoLoading(false);
      setLoadingAppId(null);
    } catch (error) {
      console.error("Error:", error);
      setGameInfoLoading(false);
      setLoadingAppId(null);
    }
  };

  console.log(data);
  let dataJSON = [];
  if (data) {
    dataJSON = JSON.parse(data);
  }
  if (data === null) {
    return;
  } else if (Object.keys(dataJSON).length === 0) {
    return (
      <div className="flex justify-center w-60 text-left bg-gameView">
        No Games Found
      </div>
    );
  } else {
    return (
      <>
        <ScrollArea className="flex mt-4">
          <div className="flex flex-col gap-2 p-5">
            {Object.values(dataJSON).map((game: any) => (
              <Button
                key={game.appid}
                variant={"outline"}
                className={`flex h-16 min-w-56 max-w-72 xl:max-w-none xl:min-w-80 overflow-hidden justify-between ${
                  loadingAppId === game.appid
                    ? "animate-pulse bg-black/10"
                    : null
                }`}
                onClick={() => IgdbGameClicked(game.appid)}
              >
                <div className="flex w-full flex-col gap-1 overflow-hidden">
                  <div className="mr-auto">{game.name}</div>
                  <div className="flex">
                    {loadingAppId === game.appid && (
                      <Loader2 className="animate-spin" />
                    )}
                    <div className="mt-auto ml-auto">
                      {new Date(game.date).getFullYear()}
                    </div>
                  </div>
                </div>
              </Button>
            ))}
          </div>
        </ScrollArea>
      </>
    );
  }
}
