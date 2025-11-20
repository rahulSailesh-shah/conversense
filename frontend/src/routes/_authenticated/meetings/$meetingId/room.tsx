import { createFileRoute } from "@tanstack/react-router";

import "@livekit/components-styles";
import { MeetingRoomView } from "@/modules/meetings/ui/views/meeting-room-view";

export const Route = createFileRoute(
  "/_authenticated/meetings/$meetingId/room"
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <MeetingRoomView />;
}
