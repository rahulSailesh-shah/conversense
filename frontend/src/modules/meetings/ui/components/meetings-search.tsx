import { useNavigate, useSearch } from "@tanstack/react-router";
import { Route } from "@/routes/_authenticated/_dashboard/meetings";
import { useEntitySearch } from "@/hooks/use-search";
import { SearchIcon } from "lucide-react";
import { Input } from "@/components/ui/input";

export const MeetingsSearch = () => {
  const search = useSearch({
    from: "/_authenticated/_dashboard/meetings/",
  });

  const navigate = useNavigate({
    from: Route.fullPath,
  });

  const { searchValue, onSearchChange } = useEntitySearch({
    params: search,
    setParams: (params) => navigate({ search: params }),
  });

  return (
    <div className="relative mr-auto">
      <SearchIcon className="size-3.5 absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground" />
      <Input
        className="max-w-[200px] bg-background shadow-none border-border pl-8"
        type="text"
        value={searchValue}
        onChange={(e) => onSearchChange(e.target.value)}
        placeholder="Search meetings"
      />
    </div>
  );
};
