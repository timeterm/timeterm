import { Typography } from "@rmwc/typography";
import { TextField } from "@rmwc/textfield";
import React, { forwardRef, Ref } from "react";
import { fetchAuthnd } from "../DevicesPage";
import { UserResponse } from "../AppDrawer";
import { Savable } from "../SettingsPage";
import useSetting, { SettingPageProps } from "./useSetting";

interface UserPatch {
  id?: string;
  name?: string;
}

const updateUser = (patch: UserPatch) =>
  fetchAuthnd(`/api/user/${patch.id}`, {
    method: "PATCH",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(patch),
  });

interface UserSettingsProps extends SettingPageProps {}

const UserSettings = forwardRef(
  (props: UserSettingsProps, ref: Ref<Savable | undefined>) => {
    const { patch, setPatch } = useSetting<UserResponse, UserPatch>({
      ref: ref,
      pageProps: props,
      isModified: (original, patch) => {
        return original.name !== patch.name;
      },
      fetch(): Promise<UserResponse> {
        return fetchAuthnd("/api/user/me").then((res) => res.json());
      },
      initPatch(original: UserResponse): UserPatch {
        return { id: original.id, name: original.name };
      },
      queryKey: "user",
      save: updateUser,
      settingsKey: "user",
    });

    return (
      <>
        <Typography use="headline5">Mijn account</Typography>

        <TextField
          style={{
            marginTop: 16,
            width: "25em",
          }}
          label={"Naam"}
          outlined
          value={patch?.name || ""}
          onChange={(evt) => {
            setPatch({
              ...patch,
              name: (evt.target as HTMLInputElement).value,
            });
          }}
        />
      </>
    );
  }
);

export default UserSettings;
