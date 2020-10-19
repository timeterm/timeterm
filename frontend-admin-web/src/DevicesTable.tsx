import React, { useMemo } from "react";
import { Icon } from "@rmwc/icon";
import { useMutation } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import "@rmwc/linear-progress/styles";
import { Column, IdType } from "react-table";
import { queryCache } from "./App";
import GeneralTable from "./GeneralTable";
import EditableCell from "./EditableCell";

export enum PrimaryDeviceStatus {
  Online = "Online",
  Offline = "Offline",
}

interface PrimaryDeviceStatusIconProps {
  status: PrimaryDeviceStatus;
}

const PrimaryDeviceStatusIcon: React.FC<PrimaryDeviceStatusIconProps> = ({
  status,
}) => {
  return (
    <Icon
      icon={status === PrimaryDeviceStatus.Online ? "check_circle" : "warning"}
      style={{
        color: status === PrimaryDeviceStatus.Online ? "#4ecd6a" : "#ffab00",
      }}
    />
  );
};

function primaryDeviceStatusString(s: PrimaryDeviceStatus): string {
  switch (s) {
    case PrimaryDeviceStatus.Online:
      return "Online";
    case PrimaryDeviceStatus.Offline:
      return "Offline";
  }
}

export interface Device {
  primaryStatus: PrimaryDeviceStatus;
  name: string;
  id: string;
}

interface DevicesTableProps {
  setSelectedItems: (items: Device[]) => void;
}

interface DevicePatch {
  id: string;
  name?: string;
}

const updateDevice = (patch: DevicePatch) =>
  fetchAuthnd(`https://api.timeterm.nl/device/${patch.id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });

const DevicesTable: React.FC<DevicesTableProps> = ({ setSelectedItems }) => {
  const [updateDeviceMut] = useMutation(updateDevice, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("devices");
    },
  });

  const columns = useMemo<Array<Column<Device>>>(
    () => [
      {
        id: "name",
        Header: "Naam",
        accessor: (dev) => dev.name,
        Cell: EditableCell,
      },
      {
        id: "status",
        Header: "Status",
        accessor: (dev) => (
          <div style={{ display: "flex", alignItems: "center" }}>
            <PrimaryDeviceStatusIcon status={dev.primaryStatus} />
            &nbsp;
            {primaryDeviceStatusString(dev.primaryStatus)}
          </div>
        ),
      },
    ],
    []
  );

  const updateData = async (
    columnId: IdType<Device>,
    item: Device,
    value: string,
    setSkipPageReset: (skipPageReset: boolean) => void
  ) => {
    if (columnId === "name") {
      setSkipPageReset(true);

      updateDeviceMut({
        id: item.id,
        name: value,
      })
        .then()
        .catch();
    }
  };

  const fetchData = (page: number, pageSize: number) => {
    return fetchAuthnd(
      `https://api.timeterm.nl/device?offset=${page * pageSize}&maxAmount=${pageSize}`
    ).then((res) => res.json());
  };

  return (
    <GeneralTable
      setSelectedItems={setSelectedItems}
      columns={columns}
      fetchData={fetchData}
      queryKey={"devices"}
      updateData={updateData}
    />
  );
};

export default DevicesTable;
