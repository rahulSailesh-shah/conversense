import { QueryBoundary } from "@/components/query-boundary";
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_authenticated/users/")({
  component: RouteComponent,
});

function RouteComponent() {
  const userData = useQuery({
    queryKey: ["repoData"],
    queryFn: async () => {
      const response = await fetch(
        "https://api.github.com/repos/TanStack/query",
      );
      return await response.json();
    },
  });
  return (
    <>
      <QueryBoundary query={userData}>
        {(data) => (
          <>
            <div>
              <h1>{data.full_name}</h1>
              <p>{data.description}</p>
              <strong>👀 {data.subscribers_count}</strong>{" "}
              <strong>✨ {data.stargazers_count}</strong>{" "}
              <strong>🍴 {data.forks_count}</strong>
            </div>
          </>
        )}
      </QueryBoundary>
    </>
  );
}
