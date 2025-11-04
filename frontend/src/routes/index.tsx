import { requireAuth } from "@/lib/auth-utils";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    await requireAuth();
    throw redirect({ to: "/meetings" });
  },
  component: RouteComponent,
});

function RouteComponent() {
  return <></>;
}
