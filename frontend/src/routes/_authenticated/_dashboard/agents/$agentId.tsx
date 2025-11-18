import { QueryBoundary } from "@/components/query-boundary";
import { useQueryAgent } from "@/modules/agents/hooks/use-agents";

import { AgentDetailsView } from "@/modules/agents/ui/views/agent-details-view";
import { createFileRoute } from "@tanstack/react-router";
import { EmptyState } from "@/components/empty-state";

export const Route = createFileRoute(
  "/_authenticated/_dashboard/agents/$agentId"
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { agentId } = Route.useParams();
  const agentQuery = useQueryAgent(agentId);
  return (
    <>
      <QueryBoundary
        query={agentQuery}
        emptyFallback={
          <EmptyState
            title="Create your first agent"
            description="Create an agent to join your meetings. Each agent will follow the instructions you provide and can interact with your users in real-time."
          />
        }
      >
        {(data) => <AgentDetailsView data={data} />}
      </QueryBoundary>
    </>
  );
}
