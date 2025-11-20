import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import "@livekit/components-styles";
import { Route } from "@/routes/_authenticated/meetings/$meetingId/room";
import { VideoConferenceClientImpl } from "../components/video-conference";

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
      <div className="flex items-center justify-center h-screen">
        <p>Loading meeting...</p>
      </div>
    );
  }

  return (
    <main data-lk-theme="default" style={{ height: "100%" }}>
      <VideoConferenceClientImpl
        liveKitUrl={SERVER_URL}
        token={token}
        codec="vp8"
      />
    </main>
  );
};
