import React, { useEffect, useMemo, useState } from "react";
import { Icon } from "@rmwc/icon";
import { Checkbox, CheckboxProps } from "@rmwc/checkbox";
import {
  DataTable,
  DataTableBody,
  DataTableCell,
  DataTableContent,
  DataTableHead,
  DataTableHeadCell,
  DataTableRow,
} from "@rmwc/data-table";
import { Set } from "immutable";
import { Select } from "@rmwc/select";
import { IconButton } from "@rmwc/icon-button";

export enum DeviceStatus {
  Online = "Online",
  Offline = "Offline",
}

interface DeviceStatusIconProps {
  status: DeviceStatus;
}

const DeviceStatusIcon: React.FC<DeviceStatusIconProps> = ({ status }) => {
  return (
    <Icon
      icon={status === DeviceStatus.Online ? "check_circle" : "warning"}
      style={{
        color: status === DeviceStatus.Online ? "#4ecd6a" : "#ffab00",
      }}
    />
  );
};

function deviceStatusString(s: DeviceStatus): string {
  switch (s) {
    case DeviceStatus.Online:
      return "Online";
    case DeviceStatus.Offline:
      return "Offline";
  }
}

export interface Device {
  status: DeviceStatus;
  name: string;
  id: string;
}

enum SelectionStatus {
  None,
  Some,
  All,
}

function oppositeSelectionStatus(s: SelectionStatus): SelectionStatus {
  return s === SelectionStatus.All ? SelectionStatus.None : SelectionStatus.All;
}

interface DevicesTableProps {
  devices: Paginated<Device>;
  setSelectedItems: (items: Device[]) => void;
}

export interface Paginated<T> {
  offset: number;
  maxAmount: number;
  total: number;
  data: T[];
}

const DevicesTable: React.FC<DevicesTableProps> = ({
  devices,
  setSelectedItems,
}) => {
  // By default no items are selected.
  const [allSelected, setAllSelected] = useState(SelectionStatus.None);
  const [selectedDevices, setSelectedDevices] = useState(Set<string>([]));

  // toggleAllSelected updates the selection status when the 'all selected' checkbox is ticked.
  // If currently no items are selected, all items are selected. If some but not all items are
  // checked, all items are selected. If all items are selected, the selection is cleared.
  const toggleAllSelected = () => {
    // Get the opposite selection status of the current.
    const newStatus = oppositeSelectionStatus(allSelected);

    // Update the 'all selected' checkbox with the new status.
    setAllSelected(newStatus);

    // If the new selection status is 'All', then select all items (tick the checkbox).
    // Otherwise, unselect all items (untick the checkbox).
    if (newStatus === SelectionStatus.All) {
      setSelectedDevices(Set(devices.data.map((dev) => dev.id)));
    } else {
      setSelectedDevices(selectedDevices.clear());
    }
  };

  // allSelectedProps contains the props for the 'all selected' checkbox.
  const allSelectedProps = useMemo<CheckboxProps>(() => {
    return {
      indeterminate: allSelected === SelectionStatus.Some,
      checked: allSelected === SelectionStatus.All,
    };
  }, [allSelected]);

  // This effect is used to determine if all devices are selected (or some) to make the
  // 'all selected' checkbox show the correct state (indeterminate meaning
  // partial selection and checked meaning all on the current page selected).
  useEffect(() => {
    const numSelected = selectedDevices.size;
    const numCheckboxes = devices.data.length;

    if (numSelected === 0) {
      // Not even a single device is selected, set the status to unchecked (empty selection).
      setAllSelected(SelectionStatus.None);
    } else if (numSelected < numCheckboxes) {
      // Not all devices are selected, set the status to indeterminate (partial selection).
      setAllSelected(SelectionStatus.Some);
    } else {
      // All devices are selected, set the status to checked (full selection).
      setAllSelected(SelectionStatus.All);
    }
  }, [devices, selectedDevices, setSelectedItems]);

  useEffect(() => {
    const newSelectedDevices = selectedDevices.intersect(
      Set(
        devices.data
          .map((dev) => dev.id)
          .filter((id) => selectedDevices.has(id))
      )
    );

    if (!newSelectedDevices.equals(selectedDevices)) {
      setSelectedDevices(newSelectedDevices);

      // Propagate the selected items to the parent element so they can use the IDs of the selected devices.
      setSelectedItems(
        devices.data.filter((dev) => selectedDevices.has(dev.id))
      );
    }
  }, [devices, selectedDevices, setSelectedItems]);

  // toggleSelectionStatus toggles the selection status of the checkbox with the key i.
  const toggleSelectionStatus = (device: Device) => {
    if (selectedDevices.has(device.id)) {
      setSelectedDevices(selectedDevices.delete(device.id));
    } else {
      setSelectedDevices(selectedDevices.add(device.id));
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        height: "100%",
      }}
    >
      <DataTable
        style={{
          width: "100%",
          height: "100%",
          borderRadius: "4px 4px 0 0",
          borderTop: 0,
          borderLeft: 0,
          borderRight: 0,
        }}
      >
        <DataTableContent>
          <DataTableHead>
            <DataTableRow>
              <DataTableHeadCell
                hasFormControl
                style={{ whiteSpace: "nowrap", width: 48 }}
              >
                <Checkbox
                  {...allSelectedProps}
                  onClick={() => toggleAllSelected()}
                />
              </DataTableHeadCell>
              <DataTableHeadCell>Naam</DataTableHeadCell>
              <DataTableHeadCell>Status</DataTableHeadCell>
              <DataTableHeadCell />
            </DataTableRow>
          </DataTableHead>
          <DataTableBody>
            {devices.data.map((device, i) => {
              return (
                <DataTableRow selected={selectedDevices.has(device.id)} key={i}>
                  <DataTableCell
                    hasFormControl
                    style={{ whiteSpace: "nowrap" }}
                  >
                    <Checkbox
                      checked={selectedDevices.has(device.id)}
                      onClick={() => toggleSelectionStatus(device)}
                    />
                  </DataTableCell>

                  <DataTableCell>{device.name}</DataTableCell>
                  <DataTableCell
                    style={{
                      display: "inline-flex",
                      alignItems: "center",
                      width: "100%",
                    }}
                  >
                    <DeviceStatusIcon status={device.status} />
                    &nbsp; {deviceStatusString(device.status)}
                  </DataTableCell>

                  <DataTableCell
                    hasFormControl
                    style={{ whiteSpace: "nowrap", width: 48 }}
                  >
                    <div
                      style={{
                        display: "inline-flex",
                        justifyContent: "flex-end",
                      }}
                    >
                      <IconButton icon={"edit"} style={{ marginRight: 16 }} />
                    </div>
                  </DataTableCell>
                </DataTableRow>
              );
            })}
          </DataTableBody>
        </DataTableContent>
      </DataTable>

      <div
        style={{
          display: "flex",
          justifyContent: "flex-end",
          alignItems: "center",
          margin: 8,
        }}
      >
        <span style={{ margin: 16 }}>Items per pagina</span>
        <Select
          outlined
          enhanced
          defaultValue={"50"}
          options={["50", "75", "100"]}
        />
        <span style={{ marginLeft: 48, marginRight: 48 }}>
          {devices.offset + 1} - {devices.data.length} van de {devices.total}
        </span>
        <IconButton icon={"first_page"} disabled={devices.offset === 0} />
        <IconButton icon={"chevron_left"} disabled={devices.offset === 0} />
        <IconButton
          icon={"chevron_right"}
          disabled={devices.offset + devices.data.length >= devices.total}
        />
        <IconButton
          icon={"last_page"}
          disabled={devices.offset + devices.data.length >= devices.total}
        />
      </div>
    </div>
  );
};

export default DevicesTable;
