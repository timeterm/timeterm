import { Typography } from "@rmwc/typography";
import { TextField } from "@rmwc/textfield";
import React, { forwardRef } from "react";
import { fetchAuthnd } from "../DevicesPage";
import { Savable } from "../SettingsPage";
import useSetting, { SettingPageProps } from "./useSetting";

interface OrganizationPatch {
  id?: string;
  name?: string;
}

interface OrganizationResponse {
  id: string;
  name: string;
}

const updateOrganization = (patch: OrganizationPatch) =>
  fetchAuthnd(`https://api.timeterm.nl/organization/${patch.id}`, {
    method: "PATCH",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(patch),
  });

interface OrganizationSettingProps extends SettingPageProps {}

const OrganizationSettings = forwardRef<Savable, OrganizationSettingProps>(
  (props, ref) => {
    const { patch, setPatch } = useSetting<
      OrganizationResponse,
      OrganizationPatch
    >({
      ref: ref,
      pageProps: props,
      isModified: (original, patch) => {
        return original.name !== patch.name;
      },
      fetch(): Promise<OrganizationResponse> {
        return fetchAuthnd("https://api.timeterm.nl/user/me")
          .then((res) => res.json())
          .then((user) =>
            fetchAuthnd(
              `https://api.timeterm.nl/organization/${user.organizationId}`
            ).then((res) => res.json())
          );
      },
      initPatch(original: OrganizationResponse): OrganizationPatch {
        return { id: original.id, name: original.name };
      },
      queryKey: "organization",
      save: updateOrganization,
      settingsKey: "organization",
    });

    return (
      <>
        <Typography use="headline5">Informatie</Typography>

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

export default OrganizationSettings;