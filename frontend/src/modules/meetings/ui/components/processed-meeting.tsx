import { Card, CardContent } from "@/components/ui/card";
import {
  SparklesIcon,
  FileTextIcon,
  VideoIcon as RecordingIcon,
  MessageSquareIcon,
} from "lucide-react";
import { useState } from "react";
import type { Meeting } from "../../types";
import { MeetingRecording } from "./meeting-recording";
import { MeetingTranscript } from "./meeting-transcript";
import { cn } from "@/lib/utils";
import { MeetingSummary } from "./meeting-summary";
import { MeetingChat } from "./meeting-chat";

interface ProcessedMeetingProps {
  meeting: Meeting;
}

type TabValue = "summary" | "transcript" | "recording" | "ask-ai";

interface Tab {
  value: TabValue;
  label: string;
  icon: React.ReactNode;
}

const tabs: Tab[] = [
  {
    value: "summary",
    label: "Summary",
    icon: <FileTextIcon className="size-4" />,
  },
  {
    value: "transcript",
    label: "Transcript",
    icon: <MessageSquareIcon className="size-4" />,
  },
  {
    value: "recording",
    label: "Recording",
    icon: <RecordingIcon className="size-4" />,
  },
  {
    value: "ask-ai",
    label: "Ask AI",
    icon: <SparklesIcon className="size-4" />,
  },
];

export const ProcessedMeeting = ({ meeting }: ProcessedMeetingProps) => {
  const [activeTab, setActiveTab] = useState<TabValue>("summary");

  return (
    <div className="w-full">
      {/* Custom Tabs */}
      <div className="bg-background rounded-lg border border-border shadow-sm">
        <div className="flex gap-1">
          {tabs.map((tab) => (
            <button
              key={tab.value}
              onClick={() => setActiveTab(tab.value)}
              className={cn(
                "flex items-center gap-2 py-4 px-2 mx-4 rounded-t-md border-b-2 transition-all",
                activeTab === tab.value
                  ? "bg-background border-primary text-foreground"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              )}
            >
              {tab.icon}
              <span className="text-sm font-normal">{tab.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === "summary" && (
          <Card className="max-h-[calc(100vh-300px)] overflow-auto">
            <CardContent className="p-6">
              <MeetingSummary meeting={meeting} />
            </CardContent>
          </Card>
        )}

        {activeTab === "transcript" && (
          <Card className="max-h-[calc(100vh-300px)] overflow-auto">
            <CardContent className="p-6">
              {meeting.transcriptUrl ? (
                <MeetingTranscript meetingId={meeting.id} />
              ) : (
                <p>No transcript available for this meeting</p>
              )}
            </CardContent>
          </Card>
        )}

        {activeTab === "recording" && (
          <Card className="max-h-[calc(100vh-300px)] overflow-auto">
            <CardContent className="p-6">
              {meeting.recordingUrl ? (
                <MeetingRecording meetingId={meeting.id} />
              ) : (
                <p>No recording available for this meeting</p>
              )}
            </CardContent>
          </Card>
        )}

        {activeTab === "ask-ai" && (
          <Card className="h-[calc(100vh-300px)] flex flex-col overflow-hidden">
            <CardContent className="flex-1 p-6 min-h-0">
              <MeetingChat meeting={meeting} />
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
};
