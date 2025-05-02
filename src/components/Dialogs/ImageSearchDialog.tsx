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
import { showErrorToast } from "@/lib/toastService";
import React from "react";

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

export default function ImageSearchDialog({
  searchDialogOpen,
  setSearchDialogOpen,
  title,
  onImageSelect,
  defaultSearchSuffix,
}: ImageSearchDialogProps) {
  const [allImages, setAllImages] = useState<GoogleImage[]>([]);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState("");
  const [nextPageOffset, setNextPageOffset] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    if (!searchDialogOpen) {
      setAllImages([]);
      setNextPageOffset(0);
      setHasMore(true);
      setQuery("");
    } else {
      if (title != "") {
        const initialQuery = String(title + " " + defaultSearchSuffix);
        setQuery(initialQuery);
        handleSearch(initialQuery);
      }
    }
  }, [searchDialogOpen]);

  const getProxiedImageUrl = (originalUrl: string) => {
    if (!originalUrl) return "";
    try {
      const decoded = decodeURIComponent(originalUrl);
      return `http://localhost:50001/image-proxy?url=${encodeURIComponent(decoded)}`;
    } catch {
      return `http://localhost:50001/image-proxy?url=${encodeURIComponent(originalUrl)}`;
    }
  };

  const handleSearch = async (searchQuery?: string) => {
    const finalQuery = searchQuery || query;
    if (!finalQuery.trim()) return;

    setLoading(true);
    setError(null);
    setAllImages([]);
    setNextPageOffset(0);
    setHasMore(true);
    setQuery(finalQuery);

    try {
      const result = await window.electron.imageSearch(finalQuery, 0);
      if (result?.length > 0) {
        setAllImages(result);
      } else {
        setError("No images found.");
      }
    } catch (err) {
      console.error("Search error:", err);
      showErrorToast("Search Error", "Search failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const loadMoreImages = async () => {
    if (loading || !hasMore) return;

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
      showErrorToast("Error loading more images", String(err));
    } finally {
      setLoading(false);
    }
  };

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollHeight, scrollTop, clientHeight } = e.currentTarget;
    const nearBottom = scrollHeight - scrollTop - clientHeight < 100;
    if (nearBottom) loadMoreImages();
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
          {loading && allImages.length === 0 ? (
            <div className="flex justify-center items-center">
              <Loader2 size={50} className="animate-spin" />
            </div>
          ) : (
            allImages.map((img, index) => (
              <ImageWithFallback
                key={`${img.ImageUrl}-${index}`}
                image={img}
                index={index}
                onImageSelect={onImageSelect}
                setSearchDialogOpen={setSearchDialogOpen}
                getProxiedImageUrl={getProxiedImageUrl}
              />
            ))
          )}
        </div>

        {loading && allImages.length > 0 && (
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

const ImageWithFallback = React.memo(
  ({
    image,
    index,
    onImageSelect,
    setSearchDialogOpen,
    getProxiedImageUrl,
  }: {
    image: GoogleImage;
    index: number;
    onImageSelect: (url: string) => void;
    setSearchDialogOpen: React.Dispatch<React.SetStateAction<boolean>>;
    getProxiedImageUrl: (url: string) => string;
  }) => {
    const [url, setUrl] = useState(image.ImageUrl);
    const [hasError, setHasError] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [attemptedProxy, setAttemptedProxy] = useState(false);

    const handleError = () => {
      if (!attemptedProxy) {
        const proxiedURL = getProxiedImageUrl(image.ImageUrl);
        setUrl(proxiedURL);
        setAttemptedProxy(true);
      } else {
        setHasError(true);
        setIsLoading(false);
      }
    };

    return (
      <div
        key={`${image.ImageUrl}-${index}`}
        className="flex w-[calc(16*1.5rem)] h-[calc(9*1.5rem)] m-2 border-2 p-1 border-muted overflow-hidden hover:bg-topBarButtonsHover rounded-md relative justify-center items-center"
      >
        {hasError ? (
          <div className="text-sm">Error Loading Image</div>
        ) : (
          <>
            {isLoading && (
              <div className="absolute inset-0 flex items-center justify-center">
                <Loader2 className="animate-spin" size={20} />
              </div>
            )}
            <img
              src={url}
              alt={`Image ${index}`}
              className={`cursor-pointer object-contain w-full h-full rounded-md ${isLoading ? "opacity-0" : "opacity-100"}`}
              onClick={() => {
                onImageSelect(image.ImageUrl);
                setSearchDialogOpen(false);
              }}
              onLoad={() => setIsLoading(false)}
              onError={handleError}
            />
          </>
        )}
      </div>
    );
  }
);
