import { VideoIcon } from "lucide-react";
import type { Meeting } from "../../types";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/empty-state";
import { useMutationStartMeeting } from "../../hooks/use-meetings";
import { useNavigate } from "@tanstack/react-router";

interface UpcomingMeetingProps {
  meeting: Meeting;
}

export const UpcomingMeeting = ({ meeting }: UpcomingMeetingProps) => {
  const startMeeting = useMutationStartMeeting();
  const navigate = useNavigate();

  const handleStartMeeting = () => {
    startMeeting.mutate(meeting.id, {
      onSuccess: (data) => {
        if (!data?.token) {
          console.error("No token received from server");
          return;
        }
        sessionStorage.setItem(`meeting-token-${meeting.id}`, data.token);
        navigate({
          to: "/meetings/$meetingId/room",
          params: { meetingId: meeting.id },
        });
      },
      onError: (error) => {
        console.error("Failed to start meeting:", error);
      },
    });
  };

  return (
    <div className="flex-1 flex justify-center">
      <div className="flex flex-col items-center gap-y-6 max-w-md text-center">
        <EmptyState
          title="Meeting not started"
          description="Once you start this meeting, a summary will appear here"
        />
        <div className="flex items-center gap-x-3">
          <Button onClick={handleStartMeeting}>
            <VideoIcon className="size-4" />
            Start meeting
          </Button>
        </div>
      </div>
    </div>
  );
};
