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

const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

type LoginSchema = z.infer<typeof loginSchema>;

const LoginForm = () => {
  const router = useRouter();
  const form = useForm<LoginSchema>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (data: LoginSchema) => {
    try {
      const { error } = await authClient.signIn.email(
        {
          email: data.email,
          password: data.password,
          callbackURL: "/",
        },
        {
          onSuccess: () => {
            toast.success("User logged in successfully");
            router.navigate({ to: "/" });
          },
          onError: (ctx) => {
            toast.error(ctx.error.message);
          },
        }
      );
      if (error) {
        toast.error(error.message);
      }
    } catch {
      toast.error("Something went wrong: ");
    }
  };

  const isPending = form.formState.isSubmitting;

  return (
    <div className="h-screen flex flex-col gap-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle>Welcome back</CardTitle>
          <CardDescription>Enter your email below to login</CardDescription>
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
                    <Button type="submit" disabled={isPending}>
                      Login
                    </Button>
                  </div>
                  <div className="text-center text-sm">
                    Don't have an account?{" "}
                    <Link className="underline underline-offset-4" to="/signup">
                      Register
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

export default LoginForm;
