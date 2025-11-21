import type { Meeting } from "../../types";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MoreVertical, Pencil, Trash2 } from "lucide-react";
import { UpcomingMeeting } from "../components/upcoming-meeting";
import { ProcessedMeeting } from "../components/processed-meeting";
import { useState } from "react";
import { DeleteDialog } from "../components/delete-dialog";
import { EditMeetingDialog } from "../components/edit-meeting-dialog";
import { useMutationDeleteMeeting } from "../../hooks/use-meetings";
import { useRouter } from "@tanstack/react-router";
import { defaultSearchParams } from "@/config/search";

export const MeetingDetailsView = ({ data }: { data: Meeting }) => {
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const router = useRouter();
  const deleteMeetingMutation = useMutationDeleteMeeting();

  const handleDelete = () => {
    deleteMeetingMutation.mutate(data.id, {
      onSuccess: () => {
        setIsDeleteDialogOpen(false);
        router.navigate({
          to: "/meetings",
          search: { ...defaultSearchParams },
        });
      },
    });
  };

  return (
    <div className="flex-1 pb-4 px-4 md:p-8 flex flex-col gap-y-6">
      <Breadcrumb className="mb-8">
        <BreadcrumbList className="text-xl">
          <BreadcrumbItem>
            <BreadcrumbLink
              href="/meetings"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              My Meetings
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator className="text-muted-foreground" />
          <BreadcrumbItem>
            <BreadcrumbPage className="text-foreground flex items-center gap-2">
              {data.name}
            </BreadcrumbPage>
          </BreadcrumbItem>
          <div className="ml-auto">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">More options</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {data.status === "upcoming" && (
                  <DropdownMenuItem
                    onSelect={(e) => {
                      e.preventDefault();
                      setIsEditDialogOpen(true);
                    }}
                  >
                    <Pencil className="mr-2 h-4 w-4" />
                    <span>Edit</span>
                  </DropdownMenuItem>
                )}
                <DropdownMenuItem
                  variant="destructive"
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
        </BreadcrumbList>
      </Breadcrumb>

      {data.status === "upcoming" ? (
        <UpcomingMeeting meeting={data} />
      ) : (
        <ProcessedMeeting meeting={data} />
      )}

      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        message={`Are you sure you want to delete the meeting "${data.name}"?`}
        onDelete={handleDelete}
      />
      <EditMeetingDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        initialValues={data}
      />
    </div>
  );
};
