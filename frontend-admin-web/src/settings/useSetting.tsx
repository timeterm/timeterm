import React, { useCallback, useEffect, useState } from "react";
import { useMutation, useQuery } from "react-query";
import { queryCache } from "../App";
import { snackbarQueue } from "../snackbarQueue";

export interface SettingPageProps {
  setIsModified: (isModified: boolean) => void;
  setIsLoading: (isLoading: boolean) => void;
  settingsStore: SettingsStore;
  setSaveChanges: (save: () => void) => void;
}

export interface SettingsStore {
  store: { [key: string]: object | undefined };
  update: (store: { [key: string]: object | undefined }) => void;
}

export interface UseSettingProps<T, P extends object> {
  fetch: () => Promise<T>;
  save: (patch: P) => Promise<unknown>;
  initPatch: (original: T) => P;
  isModified: (original: T, patch: P) => boolean;
  queryKey: string;
  pageProps: SettingPageProps;
  settingsKey: string;
  saveInvalidatesQueries?: string[];
}

const useSetting = <T, P extends object>({
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
  const setPatch = useCallback(
    (patch: P | undefined) =>
      pageProps.settingsStore.update({
        [settingsKey]: patch,
      }),
    [pageProps.settingsStore, settingsKey]
  );
  const [saved, setSaved] = useState(false);
  const onSuccess = useCallback(
    (original: T) => {
      if (original && (!patch || saved)) {
        setPatch(initPatch(original));
        setSaved(false);
      }
    },
    [patch, saved, setPatch, initPatch]
  );

  const { isLoading, error, data: original } = useQuery<T>(
    queryKey,
    () => fetch(),
    {
      refetchInterval: false,
      refetchOnMount: true,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      onSuccess,
    }
  );

  const [saveMut] = useMutation(save, {
    onSuccess: async () => {
      setSaved(true);
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

  useEffect(() => {
    pageProps.setSaveChanges(() => () => {
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
    });
  }, [patch, saveMut]); // eslint-disable-line react-hooks/exhaustive-deps
  // For the line above we really don't want pageProps to be in the dependencies array because this causes
  // an infinite loop and imperative handles also don't seem to work too well.

  return {
    original,
    patch,
    setPatch,
  };
};

export default useSetting;
