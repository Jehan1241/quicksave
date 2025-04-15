import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { TabsContent } from "@/components/ui/tabs";
import { useNavigate } from "react-router-dom";
import {
  Globe,
  Link,
  Loader2,
  LucideArrowLeft,
  LucideArrowRight,
  Plus,
  Trash2,
} from "lucide-react";
import { useEffect, useState } from "react";
import { session } from "electron";
import {
  Carousel,
  CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { CgSpinner } from "react-icons/cg";
import { useSortContext } from "@/hooks/useSortContex";
import { saveCustomImage } from "@/lib/api/GameViewAPI";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import ImageSearchDialog from "@/components/Dialogs/ImageSearchDialog";

export function ImagesTab({
  coverArtPath,
  uid,
  screenshotsArray,
  fetchData,
  title,
}: any) {
  const { cacheBuster, setCacheBuster } = useSortContext();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [coverArtLinkClicked, setCoverArtLinkClicked] = useState(false);
  const [currentCover, setCurrentCover] = useState<string | null>(
    `./backend/coverArt${coverArtPath}?t=${cacheBuster}`
  );

  const handleCoverImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        // Set the base64 result as the image source
        setCurrentCover(reader.result as string); // Cast to 'string' since 'result' can be a string or null
      };
      reader.readAsDataURL(file);
    } else {
      // If no file is selected (user cancels), clear the image
      setCurrentCover(null);
    }
  };

  const handleDeleteCoverArt = () => {
    setCurrentCover(null);
  };

  const saveClickHandler = () => {
    saveCustomImage(
      uid,
      setLoading,
      currentCover,
      ssImage,
      navigate,
      setCacheBuster,
      fetchData
    );
  };

  const [ssImage, setSsImage] = useState<(string | null)[]>(screenshotsArray);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(0);
  const [ssLinkClicked, setSSLinkClicked] = useState<number | null>(null);

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
  const [api, setApi] = useState<CarouselApi>();

  useEffect(() => {
    if (!api) {
      return;
    }

    api.on("select", (e) => {
      setSelectedIndex(e.selectedScrollSnap());
    });
  }, [api]);

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
              setCurrentCover(url); // Set as cover image
            } else if (selectedIndex !== null) {
              const updatedScreenshots = [...ssImage];
              updatedScreenshots[selectedIndex] = url;
              setSsImage(updatedScreenshots); // Set as screenshot
            }
          }}
        />
      )}
      {/* //Extra div cause TabsContent cannot be a flex */}
      <TabsContent value="images" className=" h-full w-full">
        <div className="flex h-full flex-col w-full justify-between p-2 px-4 focus:outline-none ">
          <div className="flex justify-between w-full">
            <div className="flex flex-col gap-2 w-96">
              <input
                hidden
                id="cover-picture"
                type="file"
                onChange={handleCoverImageChange}
              />
              <label htmlFor="cover-picture" className="xl:w-60 cursor-pointer">
                <Card className="flex w-28 xl:h-[calc(15rem*4/3)] xl:w-60 select-none items-center justify-center overflow-hidden">
                  <CardContent className="p-0 m-0 h-full w-full flex">
                    {currentCover ? (
                      <img
                        draggable={false}
                        src={currentCover}
                        className="object-cover h-full w-full"
                      />
                    ) : (
                      <p className="text-sm m-auto text-muted-foreground">
                        Choose an Image
                      </p>
                    )}
                  </CardContent>
                </Card>
              </label>
              <div className="flex gap-2">
                <Button
                  variant={"outline"}
                  className="w-8 h-8 rounded-full"
                  onClick={handleDeleteCoverArt}
                >
                  <Trash2 size={18} />
                </Button>
                <Button
                  variant={"outline"}
                  onClick={searchCoverClickHandler}
                  className="w-8 h-8 rounded-full"
                >
                  <Globe size={18} />
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
                          setCurrentCover(inputElement.value); // Access value
                        }
                      }}
                    ></Input>
                  )}
                </div>
              </div>
            </div>
            <div className="w-72 xl:w-full flex justify-end">
              <Carousel setApi={setApi} className="max-w-3xl w-full">
                <CarouselContent className="w-full h-full">
                  {ssImage?.map((image, index) => (
                    <CarouselItem key={index} className="w-full h-full">
                      <input
                        hidden
                        id={`ss-picture-${index}`} // Unique ID for each screenshot input
                        type="file"
                        onChange={(e) => handleScreenshotImageChange(index, e)}
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
                              const link = (e.target as HTMLInputElement).value; // Typecast to HTMLInputElement

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
          <div className="flex justify-end">
            <Button variant={"dialogSaveButton"} onClick={saveClickHandler}>
              {loading && <Loader2 className="animate-spin" />}
              Save
            </Button>
          </div>
        </div>
      </TabsContent>
    </>
  );
}
