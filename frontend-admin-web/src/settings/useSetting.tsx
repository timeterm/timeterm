import React, { Ref, useEffect, useImperativeHandle } from "react";
import { Savable } from "../SettingsPage";
import { useMutation, useQuery } from "react-query";
import { queryCache } from "../App";
import { snackbarQueue } from "../snackbarQueue";

export interface SettingPageProps {
  setIsModified: (isModified: boolean) => void;
  setIsLoading: (isLoading: boolean) => void;
  settingsStore: SettingsStore;
}

export interface SettingsStore {
  store: { [key: string]: object | undefined };
  update: (store: { [key: string]: object | undefined }) => void;
}

export interface UseSettingProps<T, P extends object> {
  ref: Ref<Savable | undefined>;
  fetch: () => Promise<T>;
  save: (patch: P) => Promise<Response>;
  initPatch: (original: T) => P;
  isModified: (original: T, patch: P) => boolean;
  queryKey: string;
  pageProps: SettingPageProps;
  settingsKey: string;
  saveInvalidatesQueries?: string[];
}

const useSetting = <T, P extends object>({
  ref,
  fetch,
  save,
  initPatch,
  isModified,
  queryKey,
  pageProps,
  settingsKey,
  saveInvalidatesQueries,
}: UseSettingProps<T, P>) => {
  const patch = pageProps.settingsStore.store[settingsKey] as P | undefined;
  const setPatch = (patch: P | undefined) =>
    pageProps.settingsStore.update({
      [settingsKey]: patch,
    });

  const { isLoading, error, data: original } = useQuery<T>(
    queryKey,
    () => fetch(),
    {
      refetchInterval: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      onSuccess(original) {
        if (original && !patch) {
          setPatch(initPatch(original));
        }
      },
    }
  );

  const [saveMut] = useMutation(save, {
    onSuccess: async () => {
      await queryCache.invalidateQueries([
        ...(saveInvalidatesQueries || []),
        queryKey,
      ]);
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
