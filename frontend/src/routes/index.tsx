import { createFileRoute, redirect } from "@tanstack/react-router";
import { requireAuth } from "@/lib/auth-utils";

export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    await requireAuth();
    throw redirect({ to: "/users" });
  },
  component: RouteComponent,
});

function RouteComponent() {
  return <></>;
}
