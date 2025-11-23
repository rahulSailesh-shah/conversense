import { GeneratedAvatar } from "@/components/generated-avatar";
import { ClockIcon, SparklesIcon } from "lucide-react";
import type { Meeting } from "../../types";
import ReactMarkdown from "react-markdown";
import { markdownComponents } from "@/components/markdown-components";

interface MeetingSummaryProps {
  meeting: Meeting;
}

export const MeetingSummary = ({ meeting }: MeetingSummaryProps) => {
  const calculateDuration = () => {
    if (meeting.startTime && meeting.endTime) {
      const start = new Date(meeting.startTime);
      const end = new Date(meeting.endTime);
      const durationMs = end.getTime() - start.getTime();
      const minutes = Math.floor(durationMs / 60000);

      if (minutes > 60) {
        const hours = Math.floor(minutes / 60);
        const remainingMinutes = minutes % 60;
        return `${hours} hr ${remainingMinutes} min`;
      }

      if (minutes == 0) {
        return "< 1 min";
      }

      return `${minutes} min`;
    }
    return null;
  };

  return (
    <div className="max-w-4xl mx-auto space-y-10 py-2">
      {/* Header Section */}
      <div className="space-y-6 border-b border-border/40 pb-8">
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground/80 uppercase tracking-wider">
            <SparklesIcon className="size-3.5" />
            <span>AI Summary</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
            {meeting.name}
          </h1>
        </div>

        <div className="flex flex-wrap items-center gap-x-6 gap-y-3 text-sm">
          {/* Agent Info */}
          <div className="flex items-center gap-2.5 bg-muted/40 px-3 py-1.5 rounded-full border border-border/40">
            <GeneratedAvatar
              seed={meeting.agentDetails.name}
              variant="botttsNeutral"
              className="size-5"
            />
            <span className="font-medium text-foreground/80">
              {meeting.agentDetails.name}
            </span>
          </div>

          {/* Date */}
          {meeting.createdAt && (
            <div className="flex items-center gap-2 text-muted-foreground">
              <span className="w-1 h-1 rounded-full bg-muted-foreground/40" />
              <span>
                {new Date(meeting.createdAt).toLocaleDateString("en-US", {
                  weekday: "long",
                  month: "long",
                  day: "numeric",
                  year: "numeric",
                })}
              </span>
            </div>
          )}

          {/* Duration */}
          {calculateDuration() && (
            <div className="flex items-center gap-2 text-muted-foreground">
              <span className="w-1 h-1 rounded-full bg-muted-foreground/40" />
              <div className="flex items-center gap-1.5">
                <ClockIcon className="size-3.5" />
                <span>{calculateDuration()}</span>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Content Section */}
      <div className="prose prose-slate dark:prose-invert max-w-none">
        {meeting.summary ? (
          <ReactMarkdown components={markdownComponents}>
            {meeting.summary}
          </ReactMarkdown>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 text-center space-y-3 bg-muted/20 rounded-lg border border-dashed border-border/60">
            <SparklesIcon className="size-8 text-muted-foreground/30" />
            <p className="text-muted-foreground">No summary generated yet.</p>
          </div>
        )}
      </div>
    </div>
  );
};
