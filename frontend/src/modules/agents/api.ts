import { apiClient, ApiError } from "@/lib/api-client";
import type { Agent } from "./types";

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
