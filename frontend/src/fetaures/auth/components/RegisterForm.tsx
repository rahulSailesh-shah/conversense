import { zodResolver } from "@hookform/resolvers/zod";
// import Image from "next/image";
import { Link } from "@tanstack/react-router";
import { useRouter } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { authClient } from "@/lib/auth-client";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const RegisterSchema = z
  .object({
    email: z.string().email(),
    password: z.string().min(8),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type RegisterSchema = z.infer<typeof RegisterSchema>;

export const RegisterForm = () => {
  const router = useRouter();
  const form = useForm<RegisterSchema>({
    resolver: zodResolver(RegisterSchema),
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (data: RegisterSchema) => {
    try {
      const { error } = await authClient.signUp.email(
        {
          email: data.email,
          password: data.password,
          name: data.email,
          callbackURL: "/",
        },
        {
          onSuccess: () => {
            toast.success("User registered successfully");
            router.navigate({ to: "/" });
          },
          onError: (ctx) => {
            toast.error(ctx.error.message);
          },
        }
      );
      if (error) {
        toast.error(error?.message || "Failed to register user");
      }
    } catch (error) {
      toast.error((error as Error)?.message || "Failed to register user");
    }
  };

  const isPending = form.formState.isSubmitting;

  return (
    <div className="h-screen flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle>Get started</CardTitle>
          <CardDescription>Create your account to get started</CardDescription>
          <CardContent className="mt-4">
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)}>
                <div className="grid gap-6">
                  <div className="flex flex-col gap-4">
                    <Button
                      variant={"outline"}
                      type="button"
                      disabled={isPending}
                      className="w-full"
                    >
                      <img
                        src="/logos/google.svg"
                        alt="Google"
                        className="w-5 h-5"
                      />
                      Continue with Google
                    </Button>
                    <Button
                      variant={"outline"}
                      type="button"
                      disabled={isPending}
                      className="w-full"
                    >
                      <img
                        src="/logos/github.svg"
                        alt="Github"
                        className="w-5 h-5"
                      />
                      Continue with Github
                    </Button>
                  </div>
                  <div className="grid gap-6">
                    <FormField
                      control={form.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email</FormLabel>
                          <FormControl>
                            <Input
                              type="email"
                              placeholder="Enter your email"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="password"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Password</FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder="Enter your password"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="confirmPassword"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Confirm Password</FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder="Enter your password"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <Button type="submit" disabled={isPending}>
                      Register
                    </Button>
                  </div>
                  <div className="text-center text-sm">
                    Already have an account?{" "}
                    <Link className="underline underline-offset-4" to="/login">
                      Login
                    </Link>
                  </div>
                </div>
              </form>
            </Form>
          </CardContent>
        </CardHeader>
      </Card>
    </div>
  );
};
