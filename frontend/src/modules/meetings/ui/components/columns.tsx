import type { ColumnDef } from "@tanstack/react-table";
import { GeneratedAvatar } from "@/components/generated-avatar";
import {
  CornerDownRightIcon,
  ClockIcon,
  PlayCircleIcon,
  CheckCircleIcon,
  XCircleIcon,
  LoaderIcon,
  CalendarClockIcon,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import type { Meeting } from "../../types";

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.

const getStatusConfig = (
  status: Meeting["status"]
): {
  icon: React.ReactNode;
  label: string;
  className: string;
} => {
  switch (status) {
    case "upcoming":
      return {
        icon: <CalendarClockIcon className="size-4" />,
        label: "Upcoming",
        className: "bg-blue-50 text-blue-700 border-blue-200",
      };
    case "active":
      return {
        icon: <PlayCircleIcon className="size-4" />,
        label: "Active",
        className: "bg-green-50 text-green-700 border-green-200",
      };
    case "completed":
      return {
        icon: <CheckCircleIcon className="size-4" />,
        label: "Completed",
        className: "bg-gray-50 text-gray-700 border-gray-200",
      };
    case "cancelled":
      return {
        icon: <XCircleIcon className="size-4" />,
        label: "Cancelled",
        className: "bg-red-50 text-red-700 border-red-200",
      };
    case "processing":
      return {
        icon: <LoaderIcon className="size-4 animate-spin" />,
        label: "Processing",
        className: "bg-yellow-50 text-yellow-700 border-yellow-200",
      };
  }
};

const calculateDuration = (meeting: Meeting): string => {
  if (meeting.startTime && meeting.endTime) {
    const start = new Date(meeting.startTime);
    const end = new Date(meeting.endTime);
    const durationMs = end.getTime() - start.getTime();
    const minutes = Math.floor(durationMs / 60000);
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;

    if (hours > 0) {
      return `${hours}h ${remainingMinutes}m`;
    }
    return `${minutes}m`;
  }

  // Default values based on status
  switch (meeting.status) {
    case "upcoming":
      return "Not started";
    case "active":
      return "In progress";
    case "processing":
      return "Processing...";
    case "cancelled":
      return "—";
    case "completed":
      return "No data";
  }
};

export const columns: ColumnDef<Meeting>[] = [
  {
    accessorKey: "name",
    header: "Meeting Name",
    cell: ({ row }) => {
      return (
        <div className="flex flex-col gap-y-1">
          <div className="flex items-center gap-x-2">
            <span className="font-semibold capitalize">
              {row.original.name}
            </span>
          </div>
          <div className="flex items-center gap-x-2">
            <CornerDownRightIcon className="size-3 text-muted-foreground" />
            <GeneratedAvatar
              seed={row.original.agentDetails.name}
              variant="botttsNeutral"
              className="size-4"
            />
            <span className="text-sm text-muted-foreground">
              {row.original.agentDetails.name}
            </span>
          </div>
        </div>
      );
    },
  },

  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const config = getStatusConfig(row.original.status);
      return (
        <Badge
          variant="outline"
          className={`flex items-center gap-x-2 w-fit ${config.className}`}
        >
          {config.icon}
          <span>{config.label}</span>
        </Badge>
      );
    },
  },

  {
    accessorKey: "duration",
    header: "Duration",
    cell: ({ row }) => {
      const duration = calculateDuration(row.original);
      return (
        <div className="flex items-center gap-x-2 text-sm">
          <ClockIcon className="size-4 text-muted-foreground" />
          <span>{duration}</span>
        </div>
      );
    },
  },
];
