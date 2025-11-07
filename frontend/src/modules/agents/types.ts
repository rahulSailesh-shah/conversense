import z from "zod";

export interface Agent {
  id: string;
  name: string;
  userId: string;
  instructions: string;
  createdAt: string;
  updatedAt: string;
}

export const agentInsertSchema = z.object({
  name: z.string().min(1, "Name is required"),
  instructions: z.string().min(1, "Instructions are required"),
});

export type NewAgent = z.infer<typeof agentInsertSchema>;
