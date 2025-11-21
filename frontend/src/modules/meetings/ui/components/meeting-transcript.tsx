import { useQueryMeetingRecording } from "../../hooks/use-meetings";
import { Loader2Icon, MessageSquareIcon, SearchIcon } from "lucide-react";
import { useEffect, useState } from "react";
import Highlighter from "react-highlight-words";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { GeneratedAvatar } from "@/components/generated-avatar";

interface MeetingTranscriptProps {
  meetingId: string;
}

interface TranscriptSegment {
  role: "user" | "ai";
  name: string;
  content: string;
  timestamp: string;
}

interface TranscriptData {
  segments: TranscriptSegment[];
}

export const MeetingTranscript = ({ meetingId }: MeetingTranscriptProps) => {
  const [transcript, setTranscript] = useState<TranscriptData | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const { data: transcriptUrl, isLoading } = useQueryMeetingRecording(
    meetingId,
    "transcript"
  );

  useEffect(() => {
    if (!transcriptUrl) return;
    fetch(transcriptUrl)
      .then((res) => res.json())
      .then((data) => {
        setTranscript(data);
      })
      .catch((err) => console.error("Failed to load transcript:", err));
  }, [transcriptUrl]);

  if (isLoading) {
    return (
      <div className="relative w-full max-w-4xl mx-auto bg-muted rounded-lg overflow-hidden p-8">
        <div className="flex flex-col items-center justify-center gap-4">
          <div className="relative">
            <MessageSquareIcon className="size-16 text-muted-foreground/30" />
            <Loader2Icon className="size-8 text-primary animate-spin absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
          </div>
          <p className="text-sm text-muted-foreground">Loading transcript...</p>
        </div>
      </div>
    );
  }

  if (!transcriptUrl) {
    return (
      <div className="w-full max-w-4xl mx-auto bg-muted rounded-lg flex items-center justify-center p-8">
        <div className="text-center">
          <MessageSquareIcon className="size-12 text-muted-foreground/50 mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">
            Transcript not available
          </p>
        </div>
      </div>
    );
  }

  if (!transcript || !transcript.segments || transcript.segments.length === 0) {
    return (
      <div className="w-full max-w-4xl mx-auto bg-muted rounded-lg flex items-center justify-center p-8">
        <div className="text-center">
          <MessageSquareIcon className="size-12 text-muted-foreground/50 mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">
            No transcript segments available
          </p>
        </div>
      </div>
    );
  }

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  };

  return (
    <div className="w-full max-w-4xl mx-auto space-y-4">
      <div className="relative">
        <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
        <Input
          type="text"
          placeholder="Search transcript..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10 bg-background border-border"
        />
      </div>

      <ScrollArea className="h-[600px] rounded-lg border bg-background">
        <div className="space-y-4 p-4">
          {transcript.segments.map((segment, index) => {
            const isUser = segment.role === "user";
            return (
              <div key={index} className="flex gap-3 items-start">
                <GeneratedAvatar
                  seed={segment.name}
                  variant={isUser ? "initials" : "botttsNeutral"}
                  className="size-8"
                />
                <div className="flex-1 min-w-0">
                  <div className="flex items-baseline gap-2 mb-1">
                    <span className="font-medium text-sm">{segment.name}</span>
                    <span className="text-xs text-muted-foreground">
                      {formatTimestamp(segment.timestamp)}
                    </span>
                  </div>
                  <div className="text-sm text-foreground bg-muted/50 rounded-lg px-4 py-2 inline-block">
                    <Highlighter
                      searchWords={searchQuery ? [searchQuery] : []}
                      autoEscape={true}
                      textToHighlight={segment.content}
                      highlightClassName="bg-yellow-300 dark:bg-yellow-600 font-semibold"
                    />
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </ScrollArea>
    </div>
  );
};
