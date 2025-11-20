import { QueryBoundary } from "@/components/query-boundary";
import { useQueryMeeting } from "@/modules/meetings/hooks/use-meetings";

// import { MeetingDetailsView } from "@/modules/meetings/ui/views/meeting-details-view";
import { createFileRoute } from "@tanstack/react-router";
import { EmptyState } from "@/components/empty-state";
import { MeetingDetailsView } from "@/modules/meetings/ui/views/meeting-details-view";

export const Route = createFileRoute(
  "/_authenticated/_dashboard/meetings/$meetingId"
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { meetingId } = Route.useParams();
  const meetingQuery = useQueryMeeting(meetingId);
  return (
    <>
      <QueryBoundary
        query={meetingQuery}
        emptyFallback={
          <EmptyState
            title="Create your first meeting"
            description="Create a meeting to join your agents. Each meeting will follow the instructions you provide and can interact with your users in real-time."
          />
        }
      >
        {(data) => <MeetingDetailsView data={data} />}
      </QueryBoundary>
    </>
  );
}
