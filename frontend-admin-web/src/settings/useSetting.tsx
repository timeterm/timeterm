import React, { Ref, useEffect, useImperativeHandle, useState } from "react";
import { Savable } from "../SettingsPage";
import { useMutation, useQuery } from "react-query";
import { queryCache } from "../App";
import { snackbarQueue } from "../snackbarQueue";

export interface SettingPageProps {
  setIsModified: (isModified: boolean) => void;
  setIsLoading: (isLoading: boolean) => void;
}

export interface UseSettingProps<T, P> {
  ref: Ref<Savable | undefined>;
  fetch: () => Promise<T>;
  save: (patch: P) => Promise<Response>;
  initPatch: (original: T) => P;
  isModified: (original: T, patch: P) => boolean;
  queryKey: string;
  pageProps: SettingPageProps;
}

const useSetting = <T, P>({
  ref,
  fetch,
  save,
  initPatch,
  isModified,
  queryKey,
  pageProps,
}: UseSettingProps<T, P>) => {
  const [patch, setPatch] = useState(undefined as P | undefined);

  const { isLoading, error, data: original } = useQuery<T>(
    queryKey,
    () => fetch(),
    {
      refetchInterval: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      onSuccess(original) {
        if (original) setPatch(initPatch(original));
      },
    }
  );

  const [saveMut] = useMutation(save, {
    onSuccess: async () => {
      await queryCache.invalidateQueries(queryKey);
    },
  });

  useEffect(() => {
    pageProps.setIsLoading(isLoading);
  }, [isLoading, pageProps]);

  useEffect(() => {
    if (original && patch) pageProps.setIsModified(isModified(original, patch));
  }, [isModified, original, pageProps, patch]);

  useEffect(() => {
    if (error)
      snackbarQueue.notify({
        title: <b>Er is een fout opgetreden</b>,
        body: "Kon data niet van server ophalen",
        icon: "error",
        dismissesOnAction: true,
        actions: [
          {
            title: "Sluiten",
            icon: "close",
          },
        ],
      });
  }, [error]);

  useImperativeHandle(ref, () => ({
    save() {
      saveMut(patch)
        .then()
        .catch(() =>
          snackbarQueue.notify({
            title: <b>Er is een fout opgetreden</b>,
            body: "Kon wijzigingen niet opslaan",
            icon: "error",
            dismissesOnAction: true,
            actions: [
              {
                title: "Sluiten",
                icon: "close",
              },
            ],
          })
        );
    },
  }));

  return {
    original,
    patch,
    setPatch,
  };
};

export default useSetting;
