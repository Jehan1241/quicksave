import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Label } from "../../ui/label";
import { Input } from "../../ui/input";
import { Button } from "../../ui/button";
import { MetadataTab } from "./MetadataTab";
import { PathTab } from "./PathTab";
import { ImagesTab } from "./ImagesTab";

export function EditDialog({
  uid,
  editDialogSelectedTab,
  setEditDialogSelectedTab,
  editDialogOpen,
  setEditDialogOpen,
  fetchData,
  coverArtPath,
  screenshotsArray,
  platform,
  tags,
}: any) {
  const isPathTabDisabled = platform === "Steam" ? true : false;

  return (
    <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
      <DialogContent className="h-[75vh] max-h-[75vh] max-w-[75vw] select-none">
        <div className="h-full overflow-y-auto flex flex-col gap-2">
          <DialogTitle>Customize Game</DialogTitle>
          <DialogDescription>
            Changes made here will reflect in the rest of the app. This can be
            reverted at any time.
          </DialogDescription>
          <Tabs
            value={editDialogSelectedTab}
            onValueChange={setEditDialogSelectedTab}
            className="flex h-full flex-col"
          >
            <TabsList className="grid w-[300px] grid-cols-3 focus:outline-none">
              <TabsTrigger value="metadata">Metadata</TabsTrigger>
              <TabsTrigger value="images">Images</TabsTrigger>
              <TabsTrigger disabled={isPathTabDisabled} value="path">
                Path
              </TabsTrigger>
            </TabsList>
            <MetadataTab tags={tags} uid={uid} fetchData={fetchData} />
            <ImagesTab
              coverArtPath={coverArtPath}
              screenshotsArray={screenshotsArray}
              uid={uid}
              fetchData={fetchData}
            />
            <PathTab uid={uid} setEditDialogOpen={setEditDialogOpen} />
          </Tabs>
        </div>
      </DialogContent>
    </Dialog>
  );
}
