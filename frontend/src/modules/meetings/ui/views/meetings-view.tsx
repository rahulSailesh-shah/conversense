import type { PaginatedMeetingResponse } from "../../types";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { Route } from "@/routes/_authenticated/_dashboard/meetings";
import { EmptyState } from "@/components/empty-state";
import { Pagination } from "@/components/pagination";

export const MeetingsListView = ({
  data,
}: {
  data: PaginatedMeetingResponse;
}) => {
  const totalPages = data.totalPages || 1;

  const search = useSearch({
    from: "/_authenticated/_dashboard/meetings/",
  });

  const navigate = useNavigate({
    from: Route.fullPath,
  });

  return (
    <div className="flex-1 pb-4 px-4 md:px-8 flex flex-col gap-y-4">
      {data.meetings.length === 0 ? (
        <EmptyState
          title="No meetings found"
          description="Create a meeting to join your meetings. Each meeting will follow the instructions you provide and can interact with your users in real-time."
        />
      ) : (
        <>
          {JSON.stringify(data)}
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
