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
  onImageSelect: (url: string) => void;
  isCoverImage?: boolean;
  defaultSearchSuffix: string;
}

const DISPLAY_BATCH_SIZE = 20;

export default function ImageSearchDialog({
  searchDialogOpen,
  setSearchDialogOpen,
  title,
  onImageSelect,
  isCoverImage,
  defaultSearchSuffix,
}: ImageSearchDialogProps) {
  const [allImages, setAllImages] = useState<GoogleImage[]>([]);
  const [displayImages, setDisplayImages] = useState<GoogleImage[]>([]);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState("");
  const [nextPageOffset, setNextPageOffset] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [currentSearchQuery, setCurrentSearchQuery] = useState("");

  useEffect(() => {
    if (!searchDialogOpen) return;

    const initialQuery = String(title + " " + defaultSearchSuffix);
    setQuery(initialQuery);
    setCurrentSearchQuery(initialQuery);
    handleSearch(initialQuery);
  }, [searchDialogOpen]);

  const getProxiedImageUrl = (originalUrl: string) => {
    if (!originalUrl) return "";
    try {
      const decoded = decodeURIComponent(originalUrl);
      return `http://localhost:8080/image-proxy?url=${encodeURIComponent(decoded)}`;
    } catch {
      return `http://localhost:8080/image-proxy?url=${encodeURIComponent(originalUrl)}`;
    }
  };

  // Reset when dialog closes
  useEffect(() => {
    if (!searchDialogOpen) {
      setAllImages([]);
      setDisplayImages([]);
      setNextPageOffset(0);
      setHasMore(true);
      setCurrentSearchQuery("");
    }
  }, [searchDialogOpen]);

  const handleSearch = async (searchQuery?: string) => {
    const finalQuery = searchQuery || query;
    if (!finalQuery.trim()) return;

    setLoading(true);
    setError(null);
    setAllImages([]);
    setDisplayImages([]);
    setNextPageOffset(0);
    setHasMore(true);
    setCurrentSearchQuery(finalQuery);

    try {
      const result = await window.electron.imageSearch(finalQuery, 0);
      if (result?.length > 0) {
        setAllImages(result);
      } else {
        setError("No images found.");
      }
    } catch (err) {
      console.error("Search error:", err);
      setError("Search failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const loadMoreImages = async () => {
    if (loading || !hasMore || query !== currentSearchQuery) return;

    setLoading(true);
    try {
      const result = await window.electron.imageSearch(query, nextPageOffset);
      if (result?.length > 0) {
        setAllImages((prev) => [...prev, ...result]);
        setNextPageOffset((prev) => prev + 10);
      } else {
        setHasMore(false);
      }
    } catch (err) {
      console.error("Load more error:", err);
    } finally {
      setLoading(false);
    }
  };

  // Update display images when allImages changes
  useEffect(() => {
    if (allImages.length > 0 && displayImages.length < allImages.length) {
      const nextBatch = allImages.slice(
        displayImages.length,
        displayImages.length + DISPLAY_BATCH_SIZE
      );
      setDisplayImages((prev) => [...prev, ...nextBatch]);
    }
  }, [allImages, displayImages.length]);

  // Handle scroll for infinite loading
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollHeight, scrollTop, clientHeight } = e.currentTarget;
    const nearBottom = scrollHeight - scrollTop - clientHeight < 100;

    if (nearBottom) {
      if (displayImages.length < allImages.length) {
        const nextBatch = allImages.slice(
          displayImages.length,
          displayImages.length + DISPLAY_BATCH_SIZE
        );
        setDisplayImages((prev) => [...prev, ...nextBatch]);
      } else if (hasMore) {
        loadMoreImages();
      }
    }
  };

  const ImageWithFallback = ({
    image,
    index,
  }: {
    image: GoogleImage;
    index: number;
  }) => {
    const [currentUrl, setCurrentUrl] = useState(image.ImageUrl);
    const [hasError, setHasError] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    const handleError = () => {
      if (!hasError) {
        setCurrentUrl(getProxiedImageUrl(image.ImageUrl));
        setHasError(true);
        setIsLoading(true); // Retry loading with proxy URL
      } else {
        setIsLoading(false);
      }
    };

    return (
      <div
        key={`${image.ImageUrl}-${index}`}
        className="flex w-[calc(16*1.5rem)] h-[calc(9*1.5rem)] m-2 border-2 p-1 border-muted overflow-hidden hover:bg-topBarButtonsHover rounded-md relative"
      >
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center">
            <Loader2 className="animate-spin" size={20} />
          </div>
        )}
        <img
          src={currentUrl}
          alt={`Image ${index}`}
          className={`cursor-pointer object-contain w-full rounded-md ${isLoading ? "opacity-0" : "opacity-100"}`}
          onClick={() => {
            onImageSelect(image.ImageUrl);
            setSearchDialogOpen(false);
          }}
          onLoad={() => setIsLoading(false)}
          onError={handleError}
        />
      </div>
    );
  };

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
          onScroll={handleScroll}
        >
          {loading && displayImages.length === 0 ? (
            <div className="flex justify-center items-center">
              <Loader2 size={50} className="animate-spin" />
            </div>
          ) : (
            displayImages.map((img, index) => (
              <ImageWithFallback
                key={`${img.ImageUrl}-${index}`}
                image={img}
                index={index}
              />
            ))
          )}
        </div>

        {loading && displayImages.length > 0 && (
          <div className="m-4 flex w-full justify-center items-center">
            <Loader2 className="animate-spin" size={30} />
          </div>
        )}

        {error && (
          <div className="text-destructive text-center p-2">{error}</div>
        )}

        <DialogFooter>
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for images"
            onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          />
          <Button onClick={() => handleSearch()}>Search</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
