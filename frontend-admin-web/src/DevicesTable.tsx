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
import { Select } from "@rmwc/select";
import { IconButton } from "@rmwc/icon-button";

export enum DeviceStatus {
  Online,
  Offline,
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
  devices: Device[];
  setSelectedItems: (items: string[]) => void;
}

const DevicesTable: React.FC<DevicesTableProps> = ({
  devices,
  setSelectedItems,
}) => {
  const [allSelected, setAllSelected] = useState(SelectionStatus.None);

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
    setCheckboxProps(
      Object.fromEntries(
        Object.entries(checkboxProps).map(([k, props]) => {
          return [k, { ...props, checked: newStatus === SelectionStatus.All }];
        })
      )
    );
  };

  // allSelectedProps contains the props for the 'all selected' checkbox.
  const allSelectedProps = useMemo<CheckboxProps>(() => {
    return {
      indeterminate: allSelected === SelectionStatus.Some,
      checked: allSelected === SelectionStatus.All,
    };
  }, [allSelected]);

  // checkboxProps contains the props for each checkbox (in each device row).
  const [checkboxProps, setCheckboxProps] = useState(
    // Create an initial list of checkbox props where each device is not selected by default.
    devices.reduce((accProps, _dev, i) => {
      return { ...accProps, [i]: { checked: false } };
    }, {} as { [key: number]: Parameters<typeof Checkbox>[0] })
  );

  // This effect is used to determine if all devices are selected (or some) to make the 
  // 'all selected' checkbox show the correct state (indeterminate meaning 
  // partial selection and checked meaning all on the current page selected).
  useEffect(() => {
    const selectedCheckboxes = Object.values(checkboxProps).filter(
      (props) => props.checked
    );

    const numSelected = selectedCheckboxes.length;
    const numCheckboxes = Object.keys(checkboxProps).length;

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

    // Propagate the selected items to the parent element so they can use the IDs of the
    // selected devices (the device ID is put into the tag of the checkbox).
    // TODO(rutgerbrf): this is really weird, perhaps don't use props?
    setSelectedItems(selectedCheckboxes.map((props) => props.tag));
  }, [checkboxProps, setSelectedItems]);

  // toggleSelectionStatus toggles the selection status of the checkbox with the key i.
  const toggleSelectionStatus = (i: number) => {
    setCheckboxProps({
      ...checkboxProps,
      [i]: {
        ...checkboxProps[i],
        checked: !checkboxProps[i].checked,
      },
    });
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
            {devices.map((dev, i) => {
              return (
                <DataTableRow selected={checkboxProps[i].checked} key={i}>
                  <DataTableCell
                    hasFormControl
                    style={{ whiteSpace: "nowrap" }}
                  >
                    <Checkbox
                      {...checkboxProps[i]}
                      tag={dev.id}
                      onClick={() => toggleSelectionStatus(i)}
                    />
                  </DataTableCell>

                  <DataTableCell>{dev.name}</DataTableCell>
                  <DataTableCell
                    style={{
                      display: "inline-flex",
                      alignItems: "center",
                      width: "100%"
                    }}
                  >
                    <DeviceStatusIcon status={dev.status} />
                    &nbsp; {deviceStatusString(dev.status)}
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
          options={["50", "75", "100", "125"]}
        />
        <span style={{ marginLeft: 48, marginRight: 48 }}>2 - 2 van de 2</span>
        <IconButton icon={"first_page"} disabled />
        <IconButton icon={"chevron_left"} disabled />
        <IconButton icon={"chevron_right"} disabled />
        <IconButton icon={"last_page"} disabled />
      </div>
    </div>
  );
};

export default DevicesTable;
