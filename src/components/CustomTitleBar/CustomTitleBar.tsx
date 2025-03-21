import { type ReactNode } from "react";
import { PiBookLight, PiListHeartLight } from "react-icons/pi";
import { HardDriveDownload } from "lucide-react";
import { useLocation } from "react-router-dom";
import QuicksaveMenu from "./QuicksaveMenu";
import TopBar from "./TopBar";
import { QuickSaveButton } from "./QuickSaveButton";

export default function CustomTitleBar({ children }: { children: ReactNode }) {
  const page = useLocation().pathname;

  return (
    <>
      <div className="flex h-screen w-screen flex-row">
        <div className="flex h-full w-14 flex-col bg-Sidebar">
          <div className="m-auto flex h-12 w-14">
            <QuicksaveMenu />
          </div>
          <div className="h-full w-14">
            <div className="my-4 flex flex-col items-center justify-start gap-4 align-middle">
              <QuickSaveButton view="" active={page === "/"}>
                <PiBookLight />
              </QuickSaveButton>
              <QuickSaveButton view="wishlist" active={page === "/wishlist"}>
                <PiListHeartLight />
              </QuickSaveButton>
              <QuickSaveButton view="installed" active={page === "/installed"}>
                <HardDriveDownload strokeWidth={1.2} />
              </QuickSaveButton>
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
