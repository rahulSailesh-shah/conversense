import { QueryBoundary } from "@/components/query-boundary";
import { useQueryAgents } from "@/modules/agents/hooks/use-agents";
import { AgentsListHeader } from "@/modules/agents/ui/components/agents-list-header";
import { AgentsListView } from "@/modules/agents/ui/views/agents-view";
import { createFileRoute } from "@tanstack/react-router";
import { EmptyState } from "@/components/empty-state";

export const Route = createFileRoute("/_authenticated/_dashboard/agents/")({
  component: RouteComponent,
});

function RouteComponent() {
  const agentsQuery = useQueryAgents();
  return (
    <>
      <AgentsListHeader />
      <QueryBoundary
        query={agentsQuery}
        emptyFallback={
          <EmptyState
            title="Create your first agent"
            description="Create an agent to join your meetings. Each agent will follow the instructions you provide and can interact with your users in real-time."
          />
        }
      >
        {(data) => <AgentsListView agents={data} />}
      </QueryBoundary>
    </>
  );
}
