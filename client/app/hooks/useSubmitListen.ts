import { useMutation, useQueryClient } from "@tanstack/react-query";
import { submitListen } from "api/api";

type SubmitListenArgs = {
  trackId: string;
  ts: Date;
};

export function useSubmitListen() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ trackId, ts }: SubmitListenArgs) =>
      submitListen(trackId, ts),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["listens"] });
    },
  });
}
