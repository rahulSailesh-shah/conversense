import { EmptyState } from "@/components/empty-state";
import { QueryBoundary } from "@/components/query-boundary";
import { useQueryMeetings } from "@/modules/meetings/hooks/use-meetings";
import { MeetingsListView } from "@/modules/meetings/ui/views/meetings-view";
import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";

const meetingSearchSchema = z.object({
  page: z.number().int().positive().catch(1),
  limit: z.number().int().positive().catch(5),
  search: z.string().catch(""),
});

export type MeetingSearchParams = z.infer<typeof meetingSearchSchema>;

export const Route = createFileRoute("/_authenticated/_dashboard/meetings/")({
  validateSearch: meetingSearchSchema,
  component: RouteComponent,
});

function RouteComponent() {
  const meetingsQuery = useQueryMeetings();

  return (
    <>
      <QueryBoundary
        query={meetingsQuery}
        emptyFallback={
          <EmptyState
            title="Create your first meeting"
            description="Create a meeting to join your meetings. Each agent will follow the instructions you provide and can interact with your users in real-time."
          />
        }
      >
        {(data) => <MeetingsListView data={data} />}
      </QueryBoundary>
    </>
  );
}
