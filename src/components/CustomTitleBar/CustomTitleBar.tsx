import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import React, { useState, type ReactNode } from "react";
import { PiBookLight, PiListHeartLight } from "react-icons/pi";

import { Filter, HardDriveDownload } from "lucide-react";
import { Dices } from "lucide-react";
import { useSortContext } from "@/hooks/useSortContex";
import { Slider } from "@/components/ui/slider";
import { useNavigate, useLocation } from "react-router-dom";
import QuicksaveMenu from "./QuicksaveMenu";
import SortGames from "./SortGames";
import FilterGames from "./FilterGames";
import WindowButtons from "./WindowsButtons";
import Integrations from "../Dialogs/Integrations";
import IntegrationsLoading from "./IntegrationsLoading";
import TopBar from "./TopBar";

export default function CustomTitleBar({ children }: { children: ReactNode }) {
  const location = useLocation();
  const navigate = useNavigate();
  const page = location.pathname;
  console.log("path", page);

  console.log(location.pathname);
  const handleViewClick = (view: "" | "wishlist" | "hidden" | "installed") => {
    navigate(`/${view}`, { replace: true });
    console.log(`${view} View Clicked`);
  };

  return (
    <>
      <div className="flex h-screen w-screen flex-row">
        <div className="flex h-full w-14 flex-col bg-Sidebar">
          <div className="m-auto flex h-12 w-14">
            <QuicksaveMenu handleViewClick={handleViewClick} />
          </div>
          <div className="h-full w-14">
            <div className="my-4 flex flex-col items-center justify-start gap-4 align-middle">
              <Button
                variant={"ghost"}
                onClick={() => handleViewClick("")}
                className={`group h-auto hover:bg-transparent ${
                  page === "/"
                    ? "rounded-none border-r border-leftbarIcons"
                    : ""
                }`}
              >
                <PiBookLight
                  className={`group-hover:scale-125 text-leftbarIcons ${
                    page === "/" ? "scale-150 group-hover:scale-150" : ""
                  }`}
                  size={22}
                />
              </Button>
              <Button
                variant={"ghost"}
                onClick={() => handleViewClick("wishlist")}
                className={`group h-auto hover:bg-transparent text-leftbarIcons ${
                  page === "/wishlist"
                    ? "rounded-none border-r border-leftbarIcons"
                    : ""
                }`}
              >
                <PiListHeartLight
                  className={`group-hover:scale-125 ${
                    page === "/wishlist"
                      ? "scale-150 group-hover:scale-150"
                      : ""
                  }`}
                  size={22}
                />
              </Button>
              <Button
                variant={"ghost"}
                onClick={() => handleViewClick("installed")}
                className={`group h-auto hover:bg-transparent text-leftbarIcons ${
                  page === "/installed"
                    ? "rounded-none border-r border-leftbarIcons"
                    : ""
                }`}
              >
                <HardDriveDownload
                  className={`group-hover:scale-125 ${
                    page === "/installed"
                      ? "scale-150 group-hover:scale-150"
                      : ""
                  }`}
                  size={22}
                  strokeWidth={1.2}
                />
              </Button>
            </div>
          </div>
        </div>
        <div className="flex h-full w-full flex-col bg-Sidebar items-center">
          <TopBar />
          <div
            draggable={false}
            className="relative h-full w-full rounded-tl-xl bg-content"
          >
            {children}
          </div>
        </div>
      </div>
    </>
  );
}
