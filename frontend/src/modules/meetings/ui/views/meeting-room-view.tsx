import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
// import "@livekit/components-styles";
import { Route } from "@/routes/_authenticated/meetings/$meetingId/room";
import { Loader2Icon } from "lucide-react";
import { LiveKitRoom } from "@livekit/components-react";
import { VideoConference } from "../components/video-conference";

const SERVER_URL =
  import.meta.env.VITE_LIVEKIT_SERVER_URL ||
  "wss://conversense-z0ptqzuw.livekit.cloud";

export const MeetingRoomView = () => {
  const { meetingId } = Route.useParams();
  const navigate = useNavigate();
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const storedToken = sessionStorage.getItem(`meeting-token-${meetingId}`);
    console.log("Stored token:", storedToken);
    if (!storedToken) {
      navigate({
        to: "/meetings/$meetingId",
        params: { meetingId },
        replace: true,
      });
      return;
    }
    setToken(storedToken);
  }, [meetingId, navigate]);

  if (!token) {
    return (
      <div className="flex flex-col items-center justify-center h-screen flex-1">
        <Loader2Icon className="size-12 animate-spin" />
      </div>
    );
  }

  return (
    <div className="h-full w-full bg-amber-600">
      <LiveKitRoom
        className="h-full w-full"
        serverUrl={SERVER_URL}
        token={token}
        data-lk-theme="default"
        audio={true}
        video={true}
        onDisconnected={() =>
          navigate({
            to: "/meetings/$meetingId",
            params: { meetingId },
            replace: true,
          })
        }
      >
        <VideoConference />
      </LiveKitRoom>
    </div>
  );
};
