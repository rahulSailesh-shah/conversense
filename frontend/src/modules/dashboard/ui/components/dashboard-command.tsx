import {
  CommandInput,
  CommandItem,
  CommandList,
  CommandResponsiveDialog,
} from "@/components/ui/command";
import type { Dispatch, SetStateAction } from "react";

interface DashboardCommandProps {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

export const DashboardCommand = ({ open, setOpen }: DashboardCommandProps) => {
  return (
    <CommandResponsiveDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Finding a meeting or a agent" />
      <CommandList>
        <CommandItem>Test</CommandItem>
      </CommandList>
    </CommandResponsiveDialog>
  );
};
