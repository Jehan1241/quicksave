import {
  Carousel,
  CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import { Button } from "../ui/button";
import { ArrowLeft, ArrowRight, Cross, FolderSymlink, X } from "lucide-react";
import { CgScrollH } from "react-icons/cg";
import { useCallback, useEffect, useRef, useState } from "react";
import React from "react";

export function CarouselSection({ uid, screenshotsArray }: any) {
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

  const [fullscreen, setFullscreen] = useState(false);

  const openFileLocation = () => {
    window.electron.openFolder(`/backend/screenshots/${uid}`);
  };

  return (
    <>
      {fullscreen && (
        <div
          className="z-20 fixed bg-black/90 top-0 left-0 flex h-full w-full flex-col items-center overflow-x-hidden py-10"
          onClick={(e) => {
            setFullscreen(false);
          }}
        >
          <button
            onClick={() => setFullscreen(false)}
            className="fixed top-2 right-2"
          >
            <X />
          </button>
          <Carousel
            onClick={(e) => e.stopPropagation()}
            opts={{ loop: true, dragFree: true }}
            setApi={setMainApi}
            onWheel={(e) => {
              if (e.deltaY > 0) {
                mainApi?.scrollNext();
              } else {
                mainApi?.scrollPrev();
              }
            }}
            className="aspect-video overflow-hidden"
          >
            <CarouselContent className="max-h-full">
              {screenshotsArray.map((_: any, index: any) => (
                <CarouselItem
                  key={index}
                  className="flex max-h-full justify-center"
                  onClick={() => setFullscreen(true)}
                >
                  <img
                    src={screenshotsArray[index]}
                    className="max-h-full max-w-full cursor-pointer rounded-xl object-contain"
                  />
                </CarouselItem>
              ))}
            </CarouselContent>
          </Carousel>
          <div
            onClick={(e) => e.stopPropagation()}
            className="flex w-full flex-col items-center justify-center overflow-y-clip"
          >
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
                {screenshotsArray.map((_: any, index: any) => (
                  <CarouselItem
                    onClick={() => {
                      mainApi?.scrollTo(index);
                      setCarouselIndex(index);
                    }}
                    key={index}
                    className="basis-36  w-full max-h-16"
                  >
                    <img
                      src={screenshotsArray[index]}
                      className={`cursor-pointer rounded-md ${
                        index === carouselIndex ? "ring-2 ring-border" : null
                      }`}
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
                className={`h-8 w-8 rounded-full ${
                  autoPlayOn ? "animate-spin duration-1000" : null
                }`}
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
      )}
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
            {screenshotsArray.map((_: any, index: any) => (
              <CarouselItem
                key={index}
                className="flex max-h-full justify-center"
                onClick={() => setFullscreen(true)}
              >
                <img
                  src={screenshotsArray[index]}
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
              {screenshotsArray.map((_: any, index: any) => (
                <CarouselItem
                  onClick={() => {
                    mainApi?.scrollTo(index);
                    setCarouselIndex(index);
                  }}
                  key={index}
                  className="basis-36  w-full max-h-16"
                >
                  <img
                    src={screenshotsArray[index]}
                    className={`cursor-pointer rounded-md ${
                      index === carouselIndex ? "ring-2 ring-border" : null
                    }`}
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
              className={`h-8 w-8 rounded-full ${
                autoPlayOn ? "animate-spin duration-1000" : null
              }`}
            >
              <CgScrollH size={22} />
            </Button>
            <Button
              onClick={openFileLocation}
              variant={"outline"}
              className={`h-8 w-8 rounded-full ${
                autoPlayOn ? "animate-spin duration-1000" : null
              }`}
            >
              <FolderSymlink size={18} />
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
    </>
  );
}
