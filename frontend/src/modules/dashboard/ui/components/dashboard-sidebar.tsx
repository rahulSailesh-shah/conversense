import { Separator } from "@/components/ui/separator";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { cn } from "@/lib/utils";
import { Link, useLocation } from "@tanstack/react-router";
import { BotIcon, VideoIcon } from "lucide-react";

const firstSection = [
  {
    icon: VideoIcon,
    label: "Meetings",
    href: "/meetings",
  },
  {
    icon: BotIcon,
    label: "Agents",
    href: "/agents",
  },
];

export const DashboardSidebar = () => {
  const { pathname } = useLocation();

  return (
    <Sidebar>
      <SidebarHeader className="text-sidebar-accent-foreground">
        <Link to="/" className="flex items-center gap-2 px-2 pt-2">
          <img
            src="/logos/logo.svg"
            alt="Conversense Logo"
            className="h-8 w-8"
          />
          <p className="text-2xl font-semibold">ConverSense</p>
        </Link>
      </SidebarHeader>
      <div className="px-4 py-2">
        <Separator className="opacity-10 text-[#5D6B68]" />
      </div>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {firstSection.map((item) => (
                <SidebarMenuItem key={item.href}>
                  <SidebarMenuButton
                    asChild
                    className={cn(
                      "h-10 border  border-transparent hover:border-[#5D6B68]/10",
                      pathname === item.href &&
                        "bg-primary-foreground border-[#5D6B68]/10"
                    )}
                    isActive={pathname === item.href}
                  >
                    <Link to={item.href} className="flex items-center gap-2">
                      <item.icon className="size-5" />
                      <span className="text-sm font-medium tracking-light">
                        {item.label}
                      </span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
};
