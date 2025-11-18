import { ResponsiveDialog } from "@/components/responsive-dialog";
import { MeetingForm } from "./meeting-form";
import { useRouter } from "@tanstack/react-router";

interface NewMeetingDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export const NewMeetingDialog = ({
  open,
  onOpenChange,
}: NewMeetingDialogProps) => {
  const router = useRouter();

  return (
    <ResponsiveDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Create New Meeting"
      description="Create a new meeting."
    >
      <MeetingForm
        onSuccess={(id) => {
          onOpenChange(false);
          router.navigate({
            to: `/meetings/${id}`,
          });
        }}
        onCancel={() => onOpenChange(false)}
      />
    </ResponsiveDialog>
  );
};
