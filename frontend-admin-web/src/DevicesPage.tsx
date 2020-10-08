import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import DevicesTable, { Device, Paginated } from "./DevicesTable";
import React, { useState } from "react";
import { queryCache, useMutation, useQuery } from "react-query";
import Cookies from "universal-cookie";

const removeDevice = (dev: Device) => {
  return fetch(`/api/device/${dev.id}`, {
    method: "DELETE",
    headers: {
      "X-Api-Key": new Cookies().get("ttsess"),
    },
  });
};

const DevicesPage: React.FC = () => {
  const [selectedItems, setSelectedItems] = useState([] as Device[]);

  const { isLoading, error, data: devices } = useQuery<Paginated<Device>>(
    "organizationDevices",
    () =>
      fetch("/api/device", {
        headers: {
          "X-Api-Key": new Cookies().get("ttsess"),
        },
      }).then((res) => res.json())
  );

  const [mutate] = useMutation(removeDevice, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("organizationDevices");
    },
  });

  const onDeleteDevices = async () => {
    for (const dev of selectedItems) {
      try {
        await mutate(dev);
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
          <DevicesTable
            devices={
              (!isLoading && !error && devices) || {
                total: 0,
                data: [],
                maxAmount: 0,
                offset: 0,
              }
            }
            setSelectedItems={setSelectedItems}
          />
        </Elevation>
      </Theme>
    </div>
  );
};

export default DevicesPage;
