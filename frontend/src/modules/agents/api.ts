import { apiClient, ApiError } from "@/lib/api-client";
import type { Agent, NewAgent } from "./types";

const handleApiError = (errorMsg: string, status: number) => {
  throw new ApiError(errorMsg, status);
};

export const fetchAgents = async () => {
  const { data, error, status } = await apiClient.get<Agent[]>("/agents");

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const createAgent = async (agentData: NewAgent) => {
  const { data, error, status } = await apiClient.post<Agent>(
    "/agents",
    agentData
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const getAgentById = async (agentId: string) => {
  const { data, error, status } = await apiClient.get<Agent>(
    `/agents/${agentId}`
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};
