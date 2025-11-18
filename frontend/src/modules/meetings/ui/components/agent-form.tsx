import { meetingInsertSchema } from "../../types";
import type { Meeting } from "../../types";
import { useForm } from "react-hook-form";
import type z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  useMutationCreateMeeting,
  useMutationUpdateMeeting,
} from "../../hooks/use-meetings";
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
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

interface AgentFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
  initialValues?: Meeting;
}

export const AgentForm = ({
  onSuccess,
  onCancel,
  initialValues,
}: AgentFormProps) => {
  const createMeetingMutation = useMutationCreateMeeting();
  const updateMeetingMutation = useMutationUpdateMeeting();

  const isEdit = !!initialValues?.id;
  const isPending =
    createMeetingMutation.isPending || updateMeetingMutation.isPending;

  const form = useForm<z.infer<typeof meetingInsertSchema>>({
    resolver: zodResolver(meetingInsertSchema),
    defaultValues: {
      name: initialValues?.name ?? "",
      agentId: initialValues?.agentId ?? "",
    },
  });

  const onSubmit = (data: z.infer<typeof meetingInsertSchema>) => {
    if (isEdit) {
      updateMeetingMutation.mutate(
        {
          id: initialValues?.id,
          ...data,
        },
        {
          onSuccess: () => {
            onSuccess?.();
          },
          onError: (error) => {
            toast.error(error.message);
            onCancel?.();
          },
        }
      );
    } else {
      createMeetingMutation.mutate(data, {
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
          name="agentId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Agent ID</FormLabel>
              <FormControl>
                <Input placeholder="e.g. Marketing Assistant" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex justify-between gap-x-2">
          {onCancel && (
            <Button
              variant="secondary"
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
