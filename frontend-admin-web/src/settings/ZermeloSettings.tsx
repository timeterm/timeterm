import { fetchAuthnd } from "../DevicesPage";
import useSetting, { SettingPageProps } from "./useSetting";
import { Typography } from "@rmwc/typography";
import { TextField } from "@rmwc/textfield";
import React from "react";
import { UserResponse } from "../AppDrawer";

interface ZermeloSettingsPatch {
  organizationId?: string;
  institution?: string;
  token?: string;
}

const saveZermeloSettings = async (patch: ZermeloSettingsPatch) => {
  if (patch.organizationId && patch.institution) {
    await fetchAuthnd(`/organization/${patch.organizationId}`, {
      method: "PATCH",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        zermelo: {
          institution: patch.institution,
        },
      }),
    });
  }

  if (patch.token) {
    await fetchAuthnd(`/zermelo/connect`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: patch.token,
      }),
    });
  }
};

const getOrganization = () =>
  fetchAuthnd(`/user/me`)
    .then((res) => res.json() as Promise<UserResponse>)
    .then((user) =>
      fetchAuthnd(`/organization/${user.organizationId}`).then(
        (res) => res.json() as Promise<Organization>
      )
    );

interface ZermeloSettingsProps extends SettingPageProps {}

interface Organization {
  id: string;
  name: string;
  zermelo?: OrganizationZermeloSettings;
}

interface OrganizationZermeloSettings {
  institution?: string;
}

const ZermeloSettings = (props: ZermeloSettingsProps) => {
  const { patch, setPatch } = useSetting<Organization, ZermeloSettingsPatch>({
    pageProps: props,
    isModified: (original, patch) => {
      return (
        original.zermelo?.institution !== patch.institution || !!patch.token
      );
    },
    fetch(): Promise<Organization> {
      return getOrganization();
    },
    initPatch(original: Organization): ZermeloSettingsPatch {
      return {
        organizationId: original.id,
        institution: original.zermelo?.institution,
      };
    },
    queryKey: "zermeloSettings",
    save: saveZermeloSettings,
    settingsKey: "zermeloSettings",
  });

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
      }}
    >
      <Typography use="headline5">Zermelo-koppeling</Typography>

      <TextField
        style={{
          marginTop: 16,
          width: "25em",
        }}
        label={"Zermelo-institutie"}
        outlined
        value={patch?.institution || ""}
        onInput={(evt) => {
          setPatch({
            ...patch,
            institution: (evt.target as HTMLInputElement).value,
          });
        }}
      />

      <TextField
        style={{
          marginTop: 16,
          width: "25em",
        }}
        label={"Token van Timeterm-gebruiker (maatwerkcoördinator) in Zermelo"}
        outlined
        value={patch?.token || ""}
        onInput={(evt) => {
          setPatch({
            ...patch,
            token: (evt.target as HTMLInputElement).value,
          });
        }}
      />
    </div>
  );
};

export default ZermeloSettings;
