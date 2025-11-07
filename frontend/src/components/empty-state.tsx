interface Props {
  title: string;
  description: string;
}

export const EmptyState = ({ title, description }: Props) => {
  return (
    <div className="flex flex-col items-center justify-center">
      <img src="/empty-state.svg" alt="" className="size-60" />
      <div className="flex flex-col text-center gap-y-6 max-w-md mx-auto">
        <h2 className="text-lg font-medium">{title}</h2>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
};
