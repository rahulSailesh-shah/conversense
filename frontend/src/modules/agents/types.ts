import z from "zod";

export interface Agent {
  id: string;
  name: string;
  userId: string;
  instructions: string;
  createdAt: string;
  updatedAt: string;
  meetingCount: number;
}

export interface PaginatedAgentResponse {
  agents: Agent[];
  currentPage: number;
  totalPages: number;
  totalCount: number;
  hasPreviousPage: boolean;
  hasNextPage: boolean;
}

export const agentInsertSchema = z.object({
  name: z.string().min(1, "Name is required"),
  instructions: z.string().min(1, "Instructions are required"),
});

export const agentUpdateSchema = z.object({
  name: z.string().min(1, "Name is required"),
  instructions: z.string().min(1, "Instructions are required"),
  id: z.string(),
});

export type AgentData = z.infer<typeof agentInsertSchema>;
export type AgentUpdateData = z.infer<typeof agentUpdateSchema>;
