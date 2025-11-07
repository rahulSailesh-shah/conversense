import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { createAgent, fetchAgents, getAgentById } from "../api";
import type { NewAgent } from "../types";

export const useQueryAgents = () => {
  return useQuery({
    queryKey: ["agents"],
    queryFn: fetchAgents,
    placeholderData: keepPreviousData,
  });
};

export const useQueryAgent = (agentId: string) => {
  return useQuery({
    queryKey: ["agent", agentId],
    queryFn: () => getAgentById(agentId),
  });
};

export const useMutationCreateAgent = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: NewAgent) => createAgent(data),
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
