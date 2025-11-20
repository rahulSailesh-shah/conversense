import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { GeneratedAvatar } from "@/components/generated-avatar";
import {
  ClockIcon,
  SparklesIcon,
  FileTextIcon,
  VideoIcon as RecordingIcon,
  MessageSquareIcon,
} from "lucide-react";
import { useState } from "react";
import type { Meeting } from "../../types";

interface ProcessedMeetingProps {
  meeting: Meeting;
}

export const ProcessedMeeting = ({ meeting }: ProcessedMeetingProps) => {
  const [activeTab, setActiveTab] = useState("summary");

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  const calculateDuration = () => {
    if (meeting.startTime && meeting.endTime) {
      const start = new Date(meeting.startTime);
      const end = new Date(meeting.endTime);
      const durationMs = end.getTime() - start.getTime();
      const minutes = Math.floor(durationMs / 60000);
      return `${minutes} minutes`;
    }
    return null;
  };

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="bg-transparent border-b rounded-none w-full justify-start h-auto p-0 space-x-6">
        <TabsTrigger
          value="summary"
          className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-3"
        >
          <FileTextIcon className="size-4 mr-2" />
          Summary
        </TabsTrigger>
        <TabsTrigger
          value="transcript"
          className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-3"
        >
          <MessageSquareIcon className="size-4 mr-2" />
          Transcript
        </TabsTrigger>
        <TabsTrigger
          value="recording"
          className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-3"
        >
          <RecordingIcon className="size-4 mr-2" />
          Recording
        </TabsTrigger>
        <TabsTrigger
          value="ask-ai"
          className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent px-0 pb-3"
        >
          <SparklesIcon className="size-4 mr-2" />
          Ask AI
        </TabsTrigger>
      </TabsList>

      <TabsContent value="summary" className="mt-6 space-y-6">
        <div>
          <h1 className="text-3xl font-bold mb-4">{meeting.name}</h1>

          <div className="flex items-center gap-x-3 mb-4">
            <GeneratedAvatar
              seed={meeting.agentDetails.name}
              variant="botttsNeutral"
              className="size-6"
            />
            <span className="font-medium">{meeting.agentDetails.name}</span>
            <span className="text-muted-foreground">
              {meeting.createdAt && formatDate(meeting.createdAt)}
            </span>
          </div>

          <div className="flex items-center gap-x-2 mb-6">
            <SparklesIcon className="size-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              General summary
            </span>
          </div>

          {calculateDuration() && (
            <Badge variant="outline" className="gap-x-2 mb-6">
              <ClockIcon className="size-4" />
              {calculateDuration()}
            </Badge>
          )}
        </div>

        <Card>
          <CardContent className="p-6">
            <h2 className="text-lg font-semibold mb-4">Overview</h2>
            <div className="text-sm text-muted-foreground leading-relaxed">
              {meeting.summary ? (
                <p>{meeting.summary}</p>
              ) : (
                <p>
                  In this insightful exchange between John Doe and{" "}
                  {meeting.agentDetails.name}, the conversation centered around
                  what constitutes a good startup idea.{" "}
                  {meeting.agentDetails.name} emphasized the importance of a
                  startup having massive market potential, a unique value
                  proposition, and the capacity to scale significantly. The
                  coach stressed the need for startups to identify problems that
                  are ripe for disruption and develop innovative solutions.
                  Furthermore, the ability to distill this idea into a
                  compelling one-liner indicates a potentially successful
                  venture. The overarching theme of the discussion was the
                  necessity of disruption in developing a winning startup idea.
                </p>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <h2 className="text-lg font-semibold mb-4">Notes</h2>
            <div className="space-y-6">
              <div>
                <h3 className="font-medium mb-3">
                  Key Elements of a Successful Startup (00:00 - 01:00)
                </h3>
                <ul className="space-y-2 text-sm text-muted-foreground list-disc list-inside">
                  <li>
                    Massive market potential is crucial for a startup idea.
                  </li>
                  <li>The importance of having a unique value proposition.</li>
                  <li>
                    Identifying a significant business opportunity is essential.
                  </li>
                </ul>
              </div>

              <div>
                <h3 className="font-medium mb-3">
                  Scalability and Innovation (01:00 - 02:30)
                </h3>
                <ul className="space-y-2 text-sm text-muted-foreground list-disc list-inside">
                  <li>
                    Startups need to demonstrate the ability to scale
                    significantly.
                  </li>
                  <li>
                    Innovation and disruption are key factors in startup
                    success.
                  </li>
                  <li>
                    The ability to articulate the idea concisely is a good
                    indicator of clarity.
                  </li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="transcript" className="mt-6">
        <Card>
          <CardContent className="p-6">
            <div className="text-center text-muted-foreground py-12">
              {meeting.transcriptUrl ? (
                <p>Transcript content would be displayed here</p>
              ) : (
                <p>No transcript available for this meeting</p>
              )}
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="recording" className="mt-6">
        <Card>
          <CardContent className="p-6">
            <div className="text-center text-muted-foreground py-12">
              {meeting.recordingUrl ? (
                <p>Recording player would be displayed here</p>
              ) : (
                <p>No recording available for this meeting</p>
              )}
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="ask-ai" className="mt-6">
        <Card>
          <CardContent className="p-6">
            <div className="text-center text-muted-foreground py-12">
              <p>AI chat interface would be displayed here</p>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  );
};
