import type { Agent } from "../../types";
import { DataTable } from "../components/data-table";
import { columns } from "../components/columns";

export const AgentsListView = ({ agents }: { agents: Agent[] }) => {
  return (
    <div className="flex-1 pb-4 px-4 md:px-8 flex flex-col gap-y-4">
      <DataTable columns={columns} data={agents} />
    </div>
  );
};
