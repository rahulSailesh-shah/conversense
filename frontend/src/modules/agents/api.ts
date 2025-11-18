import { apiClient, ApiError } from "@/lib/api-client";
import type { Agent, AgentData, AgentUpdateData } from "./types";
import type { AgentSearchParams } from "@/routes/_authenticated/_dashboard/agents";
import type { PaginatedAgentResponse } from "./types";

const handleApiError = (errorMsg: string, status: number) => {
  throw new ApiError(errorMsg, status);
};

export const fetchAgents = async (params: AgentSearchParams) => {
  const { data, error, status } = await apiClient.get<PaginatedAgentResponse>(
    "/agents",
    params
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const createAgent = async (agentData: AgentData) => {
  const { data, error, status } = await apiClient.post<Agent>(
    "/agents",
    agentData
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const updateAgent = async (agentData: AgentUpdateData) => {
  const { data, error, status } = await apiClient.put<Agent>(
    `/agents/${agentData.id}`,
    agentData
  );

  if (error) {
    handleApiError(error, status);
  }

  return data;
};

export const deleteAgent = async (agentId: string) => {
  const { data, error, status } = await apiClient.delete(`/agents/${agentId}`);

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
