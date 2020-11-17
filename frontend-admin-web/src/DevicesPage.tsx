import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import DevicesTable, { Device } from "./DevicesTable";
import React, { useState } from "react";
import { useMutation } from "react-query";
import Cookies from "universal-cookie";
import { queryCache } from "./App";
import { saveAs } from "file-saver";

const exportConfiguration = () =>
  fetchAuthnd(`device/registrationconfig`, {
    method: "GET",
    headers: {
      Accept: "application/json",
    },
  })
    .then((response) => response.blob())
    .then((blob) => saveAs(blob, "timeterm-config.json"));

const removeDevice = (devices: Device[]) =>
  fetchAuthnd(`device`, {
    method: "DELETE",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      deviceIds: devices.map((device) => device.id),
    }),
  });

const restartDevices = (devices: Device[]) =>
  fetchAuthnd(`device/restart`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      deviceIds: devices.map((device) => device.id),
    }),
  });

export const fetchAuthnd = (path: string, init?: RequestInit) => {
  const session = new Cookies().get("ttsess");
  const headers = {
    ...init?.headers,
    "X-Api-Key": session,
  };

  return fetch(process.env.REACT_APP_API_ENDPOINT + path, {
    ...init,
    headers,
  });
};

const DevicesPage: React.FC = () => {
  const [selectedItems, setSelectedItems] = useState([] as Device[]);

  const [deleteDevices] = useMutation(removeDevice, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("devices");
    },
  });

  const onDeleteDevices = async () => {
    try {
      await deleteDevices(selectedItems);
    } catch (error) {}
  };

  const [rebootDevices] = useMutation(restartDevices);

  const onRebootDevices = async () => {
    try {
      await rebootDevices(selectedItems);
    } catch (error) {}
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
        height: "100%",
      }}
    >
      <div
        style={{
          display: "flex",
          marginLeft: 32,
          marginTop: 16,
          marginRight: 16,
          height: 40,
          justifyContent: "space-between",
        }}
      >
        <h1 style={{ marginTop: 0 }}>Apparaten</h1>
        <div>
          <Button
            icon={"download"}
            outlined
            onClick={() => exportConfiguration()}
          >
            Configuratie exporteren
          </Button>
          <Button
            icon={"delete"}
            danger
            raised
            style={{ marginLeft: 8 }}
            disabled={selectedItems.length === 0}
            onClick={() => onDeleteDevices()}
          >
            Verwijderen
          </Button>
          <Button
            icon={"power_settings_new"}
            danger
            disabled={selectedItems.length === 0}
            raised
            style={{ marginLeft: 8 }}
            onClick={() => onRebootDevices()}
          >
            Opnieuw opstarten
          </Button>
        </div>
      </div>

      <Theme use={"background"} wrap>
        <Elevation
          z={16}
          style={{
            flexGrow: 1,
            margin: 16,
            borderRadius: 4,
            height: "100%",
            overflow: "hidden",
          }}
        >
          <DevicesTable setSelectedItems={setSelectedItems} />
        </Elevation>
      </Theme>
    </div>
  );
};

export default DevicesPage;
