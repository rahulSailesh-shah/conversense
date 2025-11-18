import z from "zod";

export interface Meeting {
  id: string;
  name: string;
  userId: string;
  agentId: string;
  createdAt: string;
  updatedAt: string;
}

export interface PaginatedMeetingResponse {
  meetings: Meeting[];
  currentPage: number;
  totalPages: number;
  totalCount: number;
  hasPreviousPage: boolean;
  hasNextPage: boolean;
}

export const meetingInsertSchema = z.object({
  name: z.string().min(1, "Name is required"),
  agentId: z.string().min(1, "Agent ID is required"),
});

export const meetingUpdateSchema = z.object({
  name: z.string().min(1, "Name is required"),
  agentId: z.string().min(1, "Agent ID is required"),
  id: z.string(),
});

export type MeetingData = z.infer<typeof meetingInsertSchema>;
export type MeetingUpdateData = z.infer<typeof meetingUpdateSchema>;
