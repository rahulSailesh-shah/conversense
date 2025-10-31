import { createFileRoute } from "@tanstack/react-router";
import { requireNoAuth } from "@/lib/auth-utils";
// import LoginForm from "@/modules/auth/components/LoginForm";
import { SignInView } from "@/modules/auth/ui/views/sign-in-view";

export const Route = createFileRoute("/login")({
  beforeLoad: requireNoAuth,
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <SignInView />
      </div>
    </div>
  );
}
