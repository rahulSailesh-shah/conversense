import { apiClient, ApiError } from "@/lib/api-client";
import type {
  Meeting,
  MeetingData,
  MeetingUpdateData,
  StartMeetingResponse,
} from "./types";
import type { PaginatedMeetingResponse } from "./types";
import type { MeetingSearchParams } from "@/routes/_authenticated/_dashboard/meetings";

const handleApiError = (errorMsg: string, status: number) => {
  throw new ApiError(errorMsg, status);
};

export const fetchMeetings = async (params: MeetingSearchParams) => {
  const { data, error, status } = await apiClient.get<PaginatedMeetingResponse>(
    "/meetings",
    params
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const createMeeting = async (meetingData: MeetingData) => {
  const { data, error, status } = await apiClient.post<Meeting>(
    "/meetings",
    meetingData
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const updateMeeting = async (meetingData: MeetingUpdateData) => {
  const { data, error, status } = await apiClient.put<Meeting>(
    `/meetings/${meetingData.id}`,
    meetingData
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const deleteMeeting = async (meetingId: string) => {
  const { data, error, status } = await apiClient.delete(
    `/meetings/${meetingId}`
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const getMeetingById = async (meetingId: string) => {
  const { data, error, status } = await apiClient.get<Meeting>(
    `/meetings/${meetingId}`
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const startMeeting = async (meetingId: string) => {
  const { data, error, status } = await apiClient.post<StartMeetingResponse>(
    `/meetings/${meetingId}/start`,
    {}
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const getPreSignedRecordingURL = async (
  meetingId: string,
  fileType: "recording" | "transcript"
) => {
  const { data, error, status } = await apiClient.post<string>(
    `/meetings/${meetingId}/recording-url`,
    { fileType }
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};
