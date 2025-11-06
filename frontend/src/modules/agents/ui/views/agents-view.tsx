import type { Agent } from "../../types";

export const AgentsListView = ({ agents }: { agents: Agent[] }) => {
  console.log("Rendering AgentsListView with agents:", agents);

  return <div>{JSON.stringify(agents)}</div>;
};
