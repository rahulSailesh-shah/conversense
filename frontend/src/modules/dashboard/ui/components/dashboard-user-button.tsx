import { GeneratedAvatar } from "@/components/generated-avatar";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { useIsMobile } from "@/hooks/use-mobile";
import { authClient } from "@/lib/auth-client";
import { getSession } from "@/lib/auth-utils";
import { queryClient } from "@/lib/query-client";
import { AvatarImage } from "@radix-ui/react-avatar";
import {
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { useRouter } from "@tanstack/react-router";
import { ChevronDownIcon, LogOutIcon } from "lucide-react";
import { useEffect, useState } from "react";

export const DashboardUserButton = () => {
  const [session, setSession] = useState<Awaited<
    ReturnType<typeof getSession>
  > | null>(null);

  const router = useRouter();
  const isMobile = useIsMobile();

  useEffect(() => {
    const fetchSession = async () => {
      const sessionData = await getSession();
      setSession(sessionData);
    };
    fetchSession();
  }, []);

  if (!session || !session?.user) {
    return null;
  }

  const onLogout = async () => {
    try {
      await authClient.signOut({
        fetchOptions: {
          onSuccess: () => {
            queryClient.removeQueries({ queryKey: ["session"] });
            router.navigate({ to: "/login", replace: true });
          },
        },
      });
    } catch (err) {
      console.error("Error signing out:", err);
    }
  };

  if (isMobile) {
    return (
      <Drawer>
        <DrawerTrigger className="rounded-lg border border-border/10 p-3 w-full flex items-center justify-between bg-white/5 hover:bg-white/10 overflow-hidden">
          {session.user.image ? (
            <Avatar className="mr-3">
              <AvatarImage
                src={session.user.image}
                alt={session.user.name}
                className="size-9"
              />
            </Avatar>
          ) : (
            <GeneratedAvatar
              seed={session.user.name}
              variant="initials"
              className="size-9 mr-3"
            />
          )}
          <div className="flex flex-col gap-0.5 text-left overflow-hidden flex-1 min-w-0">
            <p className="text-sm truncate w-full">{session.user.name}</p>
            <p className="text-xs truncate w-full">{session.user.email}</p>
          </div>
          <ChevronDownIcon className="size-4 shrink-0" />
        </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>{session.user.name}</DrawerTitle>
            <DrawerDescription>{session.user.email}</DrawerDescription>
          </DrawerHeader>
          <DrawerFooter>
            <Button variant="outline" onClick={onLogout}>
              Logout <LogOutIcon className="size-4" />
            </Button>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className="flex items-center gap-2 p-1.5 pr-3 rounded-full border border-border/40 hover:bg-muted/50 transition-colors outline-none">
        {session.user.image ? (
          <Avatar className="size-8">
            <AvatarImage
              src={session.user.image}
              alt={session.user.name}
              className="size-8"
            />
          </Avatar>
        ) : (
          <GeneratedAvatar
            seed={session.user.name}
            variant="initials"
            className="size-8"
          />
        )}
        <div className="flex flex-col text-left">
          <p className="text-sm font-medium truncate max-w-[150px]">
            {session.user.name}
          </p>
        </div>
        <ChevronDownIcon className="size-4 text-muted-foreground" />
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-72 p-2" align="end" side="bottom">
        <DropdownMenuLabel>
          <div className="flex flex-col gap-1">
            <span className="font-medium truncate">{session.user.name}</span>
            <span className="text-sm font-medium text-muted-foreground truncate">
              {session.user.email}
            </span>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        <DropdownMenuItem
          className="cursor-pointer flex items-center justify-between"
          onClick={onLogout}
        >
          Logout <LogOutIcon className="size-4" />
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
