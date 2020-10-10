import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import DevicesTable, { Device } from "./DevicesTable";
import React, { useState } from "react";
import { queryCache, useMutation } from "react-query";
import Cookies from "universal-cookie";

const removeDevice = (dev: Device) =>
  fetchAuthnd(`/api/device/${dev.id}`, {
    method: "DELETE",
  });

const restartDevice = (dev: Device) =>
  fetchAuthnd(`/api/device/${dev.id}/restart`, {
    method: "POST",
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

  const [deleteDevice] = useMutation(removeDevice, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("organizationDevices");
    },
  });

  const onDeleteDevices = async () => {
    for (const dev of selectedItems) {
      try {
        await deleteDevice(dev);
      } catch (error) {}
    }
  };

  const [rebootDevice] = useMutation(restartDevice);

  const onRebootDevices = async () => {
    for (const dev of selectedItems) {
      try {
        await rebootDevice(dev);
      } catch (error) {}
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
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
            borderRadius: 8,
          }}
        >
          <DevicesTable setSelectedItems={setSelectedItems} />
        </Elevation>
      </Theme>
    </div>
  );
};

export default DevicesPage;
