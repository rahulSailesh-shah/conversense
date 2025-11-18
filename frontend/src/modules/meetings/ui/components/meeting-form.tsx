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
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useQuery } from "@tanstack/react-query";
import { fetchAgents } from "@/modules/agents/api";
import { useState } from "react";
import { CommandSelect } from "@/components/command-select";
import { GeneratedAvatar } from "@/components/generated-avatar";
import { NewAgentDialog } from "@/modules/agents/ui/components/new-agent-dialog";

interface MeetingFormProps {
  onSuccess?: (id?: string) => void;
  onCancel?: () => void;
  initialValues?: Meeting;
}

export const MeetingForm = ({
  onSuccess,
  onCancel,
  initialValues,
}: MeetingFormProps) => {
  const [agentSearch, setAgentSearch] = useState("");
  const [openNewAgentDialog, setOpenNewAgentDialog] = useState(false);

  const agents = useQuery({
    queryKey: ["search_agents", agentSearch],
    queryFn: () =>
      fetchAgents({
        search: agentSearch,
        page: 1,
        limit: 100,
      }),
  });
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
          onSuccess: (data) => {
            onSuccess?.(data?.id);
          },
          onError: (error) => {
            toast.error(error.message);
            onCancel?.();
          },
        }
      );
    } else {
      createMeetingMutation.mutate(data, {
        onSuccess: (data) => {
          onSuccess?.(data?.id);
        },
        onError: (error) => {
          toast.error(error.message);
          onCancel?.();
        },
      });
    }
  };

  return (
    <>
      <NewAgentDialog
        open={openNewAgentDialog}
        onOpenChange={setOpenNewAgentDialog}
      />
      <Form {...form}>
        <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input
                    placeholder="e.g. Meeting with Math Tutor"
                    {...field}
                  />
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
                <FormLabel>Agent</FormLabel>
                <FormControl>
                  <CommandSelect
                    options={(agents.data?.agents ?? []).map((agent) => ({
                      id: agent.id,
                      value: agent.id,
                      children: (
                        <div className="flex items-center gap-x-2">
                          <GeneratedAvatar
                            seed={agent.name}
                            variant="botttsNeutral"
                            className="border size-6"
                          />
                          <span>{agent.name}</span>
                        </div>
                      ),
                    }))}
                    onSelect={field.onChange}
                    onSearch={setAgentSearch}
                    value={field.value}
                    placeholder="Select an agent"
                  />
                </FormControl>
                <FormDescription>
                  Not found what you're looking for?{" "}
                  <button
                    type="button"
                    onClick={() => setOpenNewAgentDialog(true)}
                    className="text-primary hover:underline"
                  >
                    Create a new agent
                  </button>
                </FormDescription>
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
    </>
  );
};
