// import { useRouter } from "@tanstack/react-router";
import { agentInsertSchema } from "../../types";
import type { Agent } from "../../types";
import { useForm } from "react-hook-form";
import type z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutationCreateAgent } from "../../hooks/use-agents";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { GeneratedAvatar } from "@/components/generated-avatar";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

interface AgentFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
  initialValues?: Agent;
}

export const AgentForm = ({
  onSuccess,
  onCancel,
  initialValues,
}: AgentFormProps) => {
  //   const router = useRouter();

  const createAgentMutation = useMutationCreateAgent();

  const form = useForm<z.infer<typeof agentInsertSchema>>({
    resolver: zodResolver(agentInsertSchema),
    defaultValues: {
      name: initialValues?.name ?? "",
      instructions: initialValues?.instructions ?? "",
    },
  });

  const isEdit = !!initialValues?.id;
  const isPending = createAgentMutation.isPending;

  const onSubmit = (data: z.infer<typeof agentInsertSchema>) => {
    if (isEdit) {
      console.log("Edit agent", data);
    } else {
      createAgentMutation.mutate(data, {
        onSuccess: () => {
          onSuccess?.();
        },
        onError: (error) => {
          toast.error(error.message);
          onCancel?.();
        },
      });
    }
  };

  return (
    <Form {...form}>
      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
        <GeneratedAvatar
          seed={form.watch("name")}
          variant="botttsNeutral"
          className="border size-16"
        />

        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="e.g. Marketing Assistant" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="instructions"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Instructions</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="You are a helpful assistant that can answer questions about marketing."
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex justify-between gap-x-2">
          {onCancel && (
            <Button
              variant="ghost"
              type="button"
              onClick={onCancel}
              disabled={isPending}
            >
              Cancel
            </Button>
          )}
          <Button type="submit" disabled={isPending}>
            {isEdit ? "Update" : "Create"}
          </Button>
        </div>
      </form>
    </Form>
  );
};
