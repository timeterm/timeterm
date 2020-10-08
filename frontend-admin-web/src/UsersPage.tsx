import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import DevicesTable, { Device, DeviceStatus } from "./DevicesTable";
import React, { useState } from "react";

const UsersPage: React.FC = () => {
  const [selectedItems, setSelectedItems] = useState([] as Device[]);

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
        <h1 style={{ marginTop: 0 }}>Gebruikers</h1>
        <div>
          <Button
            icon={"delete"}
            danger
            raised
            disabled={selectedItems.length === 0}
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
            devices={{
              offset: 0,
              maxAmount: 50,
              total: 2,
              data: [
                {
                  name: "Mediatheek 1",
                  id: "2f100454-7209-4379-9fc0-8f323a14da8b",
                  status: DeviceStatus.Online,
                },
                {
                  name: "Mediatheek 2",
                  id: "b1960ab9-d8dd-4a56-837e-9ccb80292ad2",
                  status: DeviceStatus.Offline,
                },
              ],
            }}
            setSelectedItems={setSelectedItems}
          />
        </Elevation>
      </Theme>
    </div>
  );
};

export default UsersPage;
