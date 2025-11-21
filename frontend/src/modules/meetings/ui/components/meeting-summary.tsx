import { GeneratedAvatar } from "@/components/generated-avatar";
import { ClockIcon, SparklesIcon } from "lucide-react";
import type { Meeting } from "../../types";
import ReactMarkdown from "react-markdown";

interface MeetingSummaryProps {
  meeting: Meeting;
}

// Reusable markdown components configuration with premium styling
const markdownComponents = {
  h1: ({ node, ...props }: any) => (
    <h1
      className="text-3xl font-bold tracking-tight text-foreground mt-10 mb-6 first:mt-0"
      {...props}
    />
  ),
  h2: ({ node, ...props }: any) => (
    <h2
      className="text-2xl font-semibold tracking-tight text-foreground mt-8 mb-4 first:mt-0"
      {...props}
    />
  ),
  h3: ({ node, ...props }: any) => (
    <h3
      className="text-xl font-semibold tracking-tight text-foreground mt-8 mb-4 pb-2 border-b border-border/40 first:mt-0"
      {...props}
    />
  ),
  h4: ({ node, ...props }: any) => (
    <h4
      className="text-base font-semibold tracking-tight text-foreground/90 mt-6 mb-3"
      {...props}
    />
  ),
  ul: ({ node, ...props }: any) => (
    <ul
      className="list-disc list-outside ml-5 space-y-2 mb-6 marker:text-muted-foreground/60"
      {...props}
    />
  ),
  ol: ({ node, ...props }: any) => (
    <ol
      className="list-decimal list-outside ml-5 space-y-2 mb-6 marker:text-muted-foreground/60"
      {...props}
    />
  ),
  li: ({ node, ...props }: any) => (
    <li className="pl-1 leading-7 text-muted-foreground" {...props} />
  ),
  p: ({ node, ...props }: any) => (
    <p className="mb-4 last:mb-0 leading-7 text-muted-foreground" {...props} />
  ),
  strong: ({ node, ...props }: any) => (
    <span className="font-semibold text-foreground" {...props} />
  ),
  em: ({ node, ...props }: any) => (
    <span className="italic text-foreground/80" {...props} />
  ),
  a: ({ node, ...props }: any) => (
    <a
      className="text-primary font-medium hover:underline underline-offset-4 transition-colors"
      {...props}
    />
  ),
  blockquote: ({ node, ...props }: any) => (
    <blockquote
      className="border-l-4 border-primary/20 bg-muted/30 pl-4 py-2 pr-4 rounded-r italic my-6 text-muted-foreground"
      {...props}
    />
  ),
  code: ({ node, ...props }: any) => (
    <code
      className="bg-muted/50 px-1.5 py-0.5 rounded text-sm font-mono text-foreground/80 border border-border/50"
      {...props}
    />
  ),
};

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
