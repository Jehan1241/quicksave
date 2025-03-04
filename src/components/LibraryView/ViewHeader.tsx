import { Button } from "@/components/ui/button";
import { Grid2X2, ListIcon } from "lucide-react";
import React from "react";
export default function ViewHeader({ view, setView, text }: any) {
  return (
    <div className="mx-5 flex items-center justify-between p-2 text-xl font-bold tracking-wide">
      <div className="flex items-center gap-2">{text}</div>
      <div className="flex gap-2">
        <Button
          className={`h-8 w-8 bg-viewButtons hover:bg-hoverViewButton ${view === "grid" ? "bg-activeViewButton" : ""}`}
          onClick={() => {
            setView("grid");
            sessionStorage.setItem("layout", "grid");
          }}
          variant={"ghost"}
        >
          <Grid2X2 strokeWidth={1.7} size={20} />
        </Button>
        <Button
          className={`h-8 w-8 bg-viewButtons hover:bg-hoverViewButton ${view === "list" ? "bg-activeViewButton" : ""}`}
          onClick={() => {
            setView("list");
            sessionStorage.setItem("layout", "list");
          }}
          variant={"ghost"}
        >
          <ListIcon size={20} />
        </Button>
      </div>
    </div>
  );
}
