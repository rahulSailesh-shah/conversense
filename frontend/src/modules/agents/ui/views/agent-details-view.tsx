import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MoreVertical, Pencil, Trash2, VideoIcon } from "lucide-react";
import type { Agent } from "../../types";
import { GeneratedAvatar } from "@/components/generated-avatar";
import { Badge } from "@/components/ui/badge";

import { useState } from "react";
import { DeleteDialog } from "../components/delete-dialog";
import { EditAgentDialog } from "../components/edit-agent-dialog";
import { useMutationDeleteAgent } from "../../hooks/use-agents";
import { useRouter } from "@tanstack/react-router";
import { defaultSearchParams } from "@/config/search";
import ReactMarkdown from "react-markdown";
import { markdownComponents } from "@/components/markdown-components";

export const AgentDetailsView = ({ data }: { data: Agent }) => {
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const router = useRouter();

  const deleteAgentMutation = useMutationDeleteAgent();

  const handleDelete = () => {
    console.log("Deleting agent:", data.id);
    deleteAgentMutation.mutate(data.id, {
      onSuccess: () => {
        setIsDeleteDialogOpen(false);
        router.navigate({
          to: "/agents",
          search: { ...defaultSearchParams },
        });
      },
    });
  };

  return (
    <div className="h-full flex-1 pb-4 px-4 md:px-8 flex flex-col gap-y-6 overflow-hidden">
      <Breadcrumb>
        <BreadcrumbList className="text-xl">
          <BreadcrumbItem>
            <BreadcrumbLink href="/agents">Agents</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{data.name}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <Card className="w-full flex-1 flex flex-col overflow-hidden">
        <CardContent className="flex-1 overflow-y-auto p-6 space-y-6">
          <div className="flex justify-between items-start">
            <div className="flex items-center space-x-4">
              <GeneratedAvatar
                seed={data.name}
                variant="botttsNeutral"
                className="size-20"
              />
              <h1 className="text-2xl font-semibold">{data.name}</h1>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">More options</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  onSelect={(e) => {
                    e.preventDefault();
                    setIsEditDialogOpen(true);
                  }}
                >
                  <Pencil className="mr-2 h-4 w-4" />
                  <span>Edit</span>
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive focus:text-destructive"
                  onSelect={(e) => {
                    e.preventDefault();
                    setIsDeleteDialogOpen(true);
                  }}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  <span>Delete</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <Badge
            variant="outline"
            className="flex items-center gap-x-2 [&>svg]:size-4"
          >
            <VideoIcon className="text-blue-700" />
            {data.meetingCount} Meeting
            {data.meetingCount === 1 ? "" : "s"}
          </Badge>

          <div className="mt-8">
            <div className="bg-muted/50 p-4 rounded-md">
              {data.instructions ? (
                <div className="prose prose-slate dark:prose-invert max-w-none">
                  <ReactMarkdown components={markdownComponents}>
                    {data.instructions}
                  </ReactMarkdown>
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">
                  No instructions provided
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        message={`Are you sure you want to delete the agent "${data.name}"? This will also delete all the meetings associated with it.`}
        onDelete={handleDelete}
      />
      <EditAgentDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        initialValues={data}
      />
    </div>
  );
};
