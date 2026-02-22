import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createApiKey } from "api/api";

export function useCreateApiKey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (label: string) => createApiKey(label),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["apiKeys"] });
    },
  });
}
