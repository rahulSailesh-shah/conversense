import { createAuthClient } from "better-auth/client";
import { jwtClient } from "better-auth/client/plugins";
import { polarClient } from "@polar-sh/better-auth";

export const authClient = createAuthClient({
  baseURL: "http://localhost:3000",
  plugins: [jwtClient(), polarClient()],
});
