import type { PaginatedAgentResponse } from "../../types";
import { DataTable } from "../components/data-table";
import { columns } from "../components/columns";
import { useRouter, useNavigate, useSearch } from "@tanstack/react-router";
import { Route } from "@/routes/_authenticated/_dashboard/agents";
import { EmptyState } from "@/components/empty-state";
import { Pagination } from "@/components/pagination";

export const AgentsListView = ({ data }: { data: PaginatedAgentResponse }) => {
  const totalPages = data.totalPages || 1;

  const router = useRouter();

  const search = useSearch({
    from: "/_authenticated/_dashboard/agents/",
  });

  const navigate = useNavigate({
    from: Route.fullPath,
  });

  return (
    <div className="flex-1 pb-4 px-4 md:px-8 flex flex-col gap-y-4">
      {data.agents.length === 0 ? (
        <EmptyState
          title="No agents found"
          description="Create an agent to join your meetings. Each agent will follow the instructions you provide and can interact with your users in real-time."
        />
      ) : (
        <>
          <DataTable
            columns={columns}
            data={data.agents}
            onRowClick={(row) =>
              router.navigate({
                to: `/agents/${row.id}`,
              })
            }
          />
          <Pagination
            page={search.page}
            totalPages={totalPages}
            onPageChange={(page) => {
              navigate({ search: { ...search, page } });
            }}
          />
        </>
      )}
    </div>
  );
};
