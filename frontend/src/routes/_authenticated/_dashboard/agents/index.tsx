import { QueryBoundary } from "@/components/query-boundary";
import { useQueryAgents } from "@/modules/agents/hooks/use-agents";
import { AgentsListView } from "@/modules/agents/ui/views/agents-view";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/_dashboard/agents/")({
  component: RouteComponent,
});

function RouteComponent() {
  const agentsQuery = useQueryAgents();
  return (
    <QueryBoundary query={agentsQuery}>
      {(data) => <AgentsListView agents={data} />}
    </QueryBoundary>
  );
}
