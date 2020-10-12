import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import DevicesTable, { Device } from "./DevicesTable";
import React, { useState } from "react";
import { useMutation } from "react-query";
import Cookies from "universal-cookie";
import { queryCache } from "./App";

const removeDevice = (devices: Device[]) =>
  fetchAuthnd(`/api/device`, {
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
  fetchAuthnd(`/api/device/restart`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      deviceIds: devices.map((device) => device.id),
    }),
  });

export const fetchAuthnd = (input: RequestInfo, init?: RequestInit) => {
  const session = new Cookies().get("ttsess");
  const headers = {
    ...init?.headers,
    "X-Api-Key": session,
  };

  return fetch(input, {
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
            icon={"delete"}
            danger
            raised
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
