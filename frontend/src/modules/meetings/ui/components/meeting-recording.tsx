import { useQueryMeetingRecording } from "../../hooks/use-meetings";
import { Loader2Icon, VideoIcon } from "lucide-react";

interface MeetingRecordingProps {
  meetingId: string;
}

export const MeetingRecording = ({ meetingId }: MeetingRecordingProps) => {
  const { data: recordingUrl, isLoading } = useQueryMeetingRecording(
    meetingId,
    "recording"
  );

  if (isLoading) {
    return (
      <div className="relative w-full max-w-4xl mx-auto aspect-video bg-muted rounded-lg overflow-hidden">
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-4 z-10">
          <div className="relative">
            <VideoIcon className="size-16 text-muted-foreground/30" />
            <Loader2Icon className="size-8 text-primary animate-spin absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
          </div>
          <p className="text-sm text-muted-foreground">Loading recording...</p>
        </div>
        <div className="absolute inset-0 animate-pulse bg-muted-foreground/5" />
      </div>
    );
  }

  if (!recordingUrl) {
    return (
      <div className="w-full max-w-4xl mx-auto aspect-video bg-muted rounded-lg flex items-center justify-center">
        <div className="text-center">
          <VideoIcon className="size-12 text-muted-foreground/50 mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">
            Recording not available
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full max-w-4xl mx-auto">
      <video
        controls
        src={recordingUrl}
        className="w-full aspect-video rounded-lg shadow-md bg-black"
        controlsList="nodownload"
      >
        Your browser does not support the video tag.
      </video>
    </div>
  );
};
