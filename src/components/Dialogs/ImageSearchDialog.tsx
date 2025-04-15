import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { Input } from "../ui/input";
import { Loader2 } from "lucide-react";

interface GoogleImage {
  ImageUrl: string;
  ThumbUrl: string;
  Width: number;
  Height: number;
}

interface ImageSearchDialogProps {
  searchDialogOpen: boolean;
  setSearchDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
  title: string;
}

export default function ImageSearchDialog({
  searchDialogOpen,
  setSearchDialogOpen,
  title,
}: ImageSearchDialogProps) {
  const [images, setImages] = useState<GoogleImage[]>([]);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1); // Track the page for infinite scroll
  const [imageErrors, setImageErrors] = useState<Set<string>>(new Set()); // Track failed images
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await window.electron.imageSearch(query);
      if (result && result.length > 0) {
        setImages(result);
      } else {
        setError("No images found.");
      }
    } catch (err) {
      console.error("Error during image search:", err);
      setError("An error occurred during the search.");
    } finally {
      setLoading(false);
    }
  };

  const loadMoreImages = async () => {
    if (loading) return; // Prevent loading more images if already loading

    setLoading(true);
    console.log("Loading more images...");

    try {
      const imageData: GoogleImage[] = await window.electron.imageSearch(
        query + `&page=${page + 1}`
      );
      setImages((prevImages) => [...prevImages, ...imageData]);
      setPage((prevPage) => prevPage + 1); // Increment page for next load
    } catch (err) {
      console.error("Error loading more images:", err);
    }

    setLoading(false);
  };

  const handleImageClick = (img: string) => {
    console.log("Image selected:", img);
    setSearchDialogOpen(false);
  };

  const handleImageError = (
    e: React.SyntheticEvent<HTMLImageElement, Event>
  ) => {
    const imageUrl = e.currentTarget.src;

    // Prevent infinite error logging for the same broken image or the fallback image
    if (imageUrl.includes("via.placeholder.com")) return; // Ignore errors for fallback image

    if (!imageErrors.has(imageUrl)) {
      console.error(`Image failed to load: ${imageUrl}`);
      setImageErrors((prevErrors) => new Set(prevErrors).add(imageUrl)); // Mark this image as having failed
    }

    // Set fallback image
    e.currentTarget.src = "https://via.placeholder.com/150";
  };

  const handleScroll = (e: React.UIEvent<HTMLDivElement, UIEvent>) => {
    const scrollHeight = e.currentTarget.scrollHeight;
    const scrollTop = e.currentTarget.scrollTop;
    const clientHeight = e.currentTarget.clientHeight;

    console.log(
      "Scroll Event - ScrollHeight:",
      scrollHeight,
      "ScrollTop:",
      scrollTop,
      "ClientHeight:",
      clientHeight
    );

    // Check if user has scrolled to the bottom (with a small tolerance)
    const nearBottom = scrollHeight - scrollTop - clientHeight < 10;

    if (nearBottom && !loading) {
      console.log("Reached bottom, loading more...");
      loadMoreImages(); // Load more images when scrolled to bottom
    }
  };

  useEffect(() => {
    console.log(images); // For debugging
  }, [images]);

  return (
    <Dialog open={searchDialogOpen} onOpenChange={setSearchDialogOpen}>
      <DialogContent className="h-2/3 max-w-none w-2/3 flex flex-col justify-between">
        <DialogHeader>
          <DialogTitle>Search Images</DialogTitle>
          <DialogDescription>
            Left click on an image to set it.
          </DialogDescription>
        </DialogHeader>
        <div
          className="h-full w-full flex flex-wrap overflow-auto justify-center"
          onScroll={handleScroll} // Add scroll event listener here
        >
          {loading && images.length === 0 ? (
            <div className="flex justify-center items-center">
              <Loader2 size={50} className="animate-spin" />
            </div> // Loading indicator for initial load
          ) : (
            images.map((img, index) => (
              <div className="flex w-[calc(16*1.5rem)] h-[calc(9*1.5rem)] m-2 border-2 p-1 border-muted overflow-hidden hover:bg-topBarButtonsHover rounded-md">
                <img
                  key={index}
                  src={img.ImageUrl}
                  alt={`Image ${index}`}
                  className="cursor-pointer object-contain w-full rounded-md"
                  onClick={() => handleImageClick(img.ImageUrl)}
                  onError={handleImageError} // Fallback to placeholder image
                />
              </div>
            ))
          )}
        </div>
        {loading && images.length > 0 && (
          <div className="m-4 flex w-full justify-center items-center">
            <Loader2 className="animate-spin" size={30} />
          </div>
        )}
        <DialogFooter>
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for images"
          />
          <Button onClick={handleSearch}>Search</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
