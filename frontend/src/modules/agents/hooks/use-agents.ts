import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import {
  createAgent,
  deleteAgent,
  fetchAgents,
  getAgentById,
  updateAgent,
} from "../api";
import type { AgentData, AgentUpdateData } from "../types";
import { useSearch } from "@tanstack/react-router";

export const useQueryAgents = () => {
  const search = useSearch({
    from: "/_authenticated/_dashboard/agents/",
  });
  return useQuery({
    queryKey: ["agents", search],
    queryFn: () => fetchAgents(search),
    retry: 0,
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
    mutationFn: (data: AgentData) => createAgent(data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["agents"],
      });
      queryClient.invalidateQueries({
        queryKey: ["search_agents"],
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

export const useMutationUpdateAgent = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: AgentUpdateData) => updateAgent(data),
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

export const useMutationDeleteAgent = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (agentId: string) => deleteAgent(agentId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["agents"],
      });
    },
    onError: () => {},
  });
};
