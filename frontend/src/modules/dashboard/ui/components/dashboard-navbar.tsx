import { Button } from "@/components/ui/button";
import { useSidebar } from "@/components/ui/sidebar";
import { PanelLeftCloseIcon, PanelLeftIcon } from "lucide-react";

import { DashboardUserButton } from "./dashboard-user-button";

import { ModeToggle } from "@/components/mode-toggle";

export const DashboardNavbar = () => {
  const { state, isMobile, toggleSidebar } = useSidebar();

  return (
    <>
      <nav className="flex px-4 gap-x-2 items-center justify-between py-3 bg-background">
        <Button className="size-9" variant="outline" onClick={toggleSidebar}>
          {state === "collapsed" || isMobile ? (
            <PanelLeftIcon className="size-4" />
          ) : (
            <PanelLeftCloseIcon className="size-4" />
          )}
        </Button>
        <div className="flex items-center gap-2">
          <ModeToggle />
          <DashboardUserButton />
        </div>
      </nav>
    </>
  );
};
