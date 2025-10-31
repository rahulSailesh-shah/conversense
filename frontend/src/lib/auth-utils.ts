import { authClient } from "./auth-client";
import { redirect } from "@tanstack/react-router";

export async function getSession() {
  const { data: session } = await authClient.getSession();
  return session;
}

export async function requireAuth() {
  const session = await getSession();
  if (!session) {
    throw redirect({ to: "/login" });
  }
  return session;
}

export async function requireNoAuth() {
  const session = await getSession();
  if (session) {
    throw redirect({ to: "/" });
  }
  return session;
}
