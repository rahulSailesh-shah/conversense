import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { fetchAgents } from "../api";

export const useQueryAgents = () => {
  return useQuery({
    queryKey: ["agents"],
    queryFn: fetchAgents,
  });
};
