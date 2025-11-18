import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import {
  createMeeting,
  deleteMeeting,
  fetchMeetings,
  getMeetingById,
  updateMeeting,
} from "../api";
import type { MeetingData, MeetingUpdateData } from "../types";
import { useSearch } from "@tanstack/react-router";

export const useQueryMeetings = () => {
  const search = useSearch({
    from: "/_authenticated/_dashboard/meetings/",
  });
  return useQuery({
    queryKey: ["meetings", search],
    queryFn: () => fetchMeetings(search),
    retry: 0,
    placeholderData: keepPreviousData,
  });
};

export const useQueryMeeting = (meetingId: string) => {
  return useQuery({
    queryKey: ["meeting", meetingId],
    queryFn: () => getMeetingById(meetingId),
  });
};

export const useMutationCreateMeeting = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: MeetingData) => createMeeting(data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["agents"],
      });
      if (data?.id) {
        queryClient.invalidateQueries({
          queryKey: ["agent", data.id],
        });
      }
    },
    onError: () => {},
  });
};

export const useMutationUpdateMeeting = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: MeetingUpdateData) => updateMeeting(data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["agents"],
      });
      if (data?.id) {
        queryClient.invalidateQueries({
          queryKey: ["agent", data.id],
        });
      }
    },
    onError: () => {},
  });
};

export const useMutationDeleteMeeting = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (meetingId: string) => deleteMeeting(meetingId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["meetings"],
      });
    },
    onError: () => {},
  });
};
