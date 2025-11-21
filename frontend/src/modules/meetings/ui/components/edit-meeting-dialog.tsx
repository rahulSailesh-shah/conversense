import { ResponsiveDialog } from "@/components/responsive-dialog";
import { MeetingForm } from "./meeting-form";
import type { Meeting } from "../../types";

interface EditMeetingDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialValues: Meeting;
}

export const EditMeetingDialog = ({
  open,
  onOpenChange,
  initialValues,
}: EditMeetingDialogProps) => {
  return (
    <ResponsiveDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Edit Meeting"
      description="Update meeting details."
    >
      <MeetingForm
        initialValues={initialValues}
        onSuccess={() => {
          onOpenChange(false);
        }}
        onCancel={() => onOpenChange(false)}
      />
    </ResponsiveDialog>
  );
};
