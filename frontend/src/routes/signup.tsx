import { createFileRoute, Link } from "@tanstack/react-router";
import { requireNoAuth } from "@/lib/auth-utils";
import { RegisterForm } from "@/fetaures/auth/components/RegisterForm";

export const Route = createFileRoute("/signup")({
  beforeLoad: requireNoAuth,
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center">
      <div className="w-full max-w-sm flex flex-col gap-6">
        <Link to="/" className="flex items-center self-center font-medium">
          <img src="/logos/logo.svg" alt="Logo" className="w-24 h-24" />
        </Link>
        <RegisterForm />
      </div>
    </div>
  );
}
