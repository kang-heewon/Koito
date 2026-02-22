import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteApiKey } from "api/api";

export function useDeleteApiKey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => deleteApiKey(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["apiKeys"] });
    },
  });
}
