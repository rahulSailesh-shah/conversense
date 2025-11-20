import type { Meeting } from "../../types";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { UpcomingMeeting } from "../components/upcoming-meeting";
import { ProcessedMeeting } from "../components/processed-meeting";

export const MeetingDetailsView = ({ data }: { data: Meeting }) => {
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
            <BreadcrumbPage className="text-foreground  ">
              {data.name}
            </BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {data.status === "upcoming" ? (
        <UpcomingMeeting meeting={data} />
      ) : (
        <ProcessedMeeting meeting={data} />
      )}
    </div>
  );
};
