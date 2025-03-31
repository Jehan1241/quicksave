import { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../ui/button";
import React from "react";
import { cn } from "@/lib/utils";

export function QuickSaveButton({
  view,
  active,
  children,
}: {
  view: "library" | "wishlist" | "hidden" | "installed";
  active: boolean;
  children: ReactNode;
}) {
  const navigate = useNavigate();
  const handleViewClick = (
    view: "library" | "wishlist" | "hidden" | "installed"
  ) => {
    navigate(`/${view}`, { replace: true });
  };
  return (
    <Button
      variant="ghost"
      onClick={() => handleViewClick(view)}
      className={cn(
        "group h-auto hover:bg-transparent",
        active && "rounded-r-none border-r border-leftbarIcons"
      )}
    >
      {React.Children.map(children, (child) =>
        React.cloneElement(child as React.ReactElement, {
          className: cn(
            (child as React.ReactElement).props.className,
            "group-hover:scale-110 transition-transform",
            active ? "scale-110" : ""
          ),
          size: 22,
        })
      )}
    </Button>
  );
}
