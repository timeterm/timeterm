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
import { Theme } from "@rmwc/theme";

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
        color: status === DeviceStatus.Online ? "#4ECD6A" : "#ffab00",
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
}

enum SelectionStatus {
  None,
  Some,
  All,
}

function oppositeSelectionStatus(s: SelectionStatus): SelectionStatus {
  switch (s) {
    case SelectionStatus.All || SelectionStatus.Some:
      return SelectionStatus.None;
    default:
      return SelectionStatus.All;
  }
}

interface DevicesTableProps {
  devices: Device[];
}

const DevicesTable: React.FC<DevicesTableProps> = ({ devices }) => {
  const [allSelected, setAllSelected] = useState(SelectionStatus.None);
  const toggleAllSelected = () => {
    const newStatus = oppositeSelectionStatus(allSelected);

    setAllSelected(newStatus);

    switch (newStatus) {
      case SelectionStatus.None:
        setCheckboxProps(
          Object.fromEntries(
            Object.entries(checkboxProps).map(([k, props]) => {
              return [k, { ...props, checked: false }];
            })
          )
        );
        break;
      case SelectionStatus.All:
        setCheckboxProps(
          Object.fromEntries(
            Object.entries(checkboxProps).map(([k, props]) => {
              return [k, { ...props, checked: true }];
            })
          )
        );
        break;
    }
  };

  const allSelectedProps = useMemo<CheckboxProps>(() => {
    switch (allSelected) {
      case SelectionStatus.Some:
        return { indeterminate: true };
      case SelectionStatus.All:
        return { checked: true };
      case SelectionStatus.None:
        return { checked: false };
    }
  }, [allSelected]);

  const [checkboxProps, setCheckboxProps] = useState(
    Object.fromEntries(
      devices.map((dev, i) => {
        return [
          i,
          {
            checked: false,
          },
        ];
      })
    ) as { [key: number]: CheckboxProps }
  );

  useEffect(() => {
    const numSelected = Object.values(checkboxProps).reduce(
      (acc, props) => acc + (props.checked ? 1 : 0),
      0
    );

    const numCheckboxes = Object.keys(checkboxProps).length;
    if (numSelected === 0) {
      setAllSelected(SelectionStatus.None);
    } else if (numSelected < numCheckboxes) {
      setAllSelected(SelectionStatus.Some);
    } else {
      setAllSelected(SelectionStatus.All);
    }
  }, [checkboxProps]);

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
    <DataTable
      style={{
        width: "100%",
        height: "100%",
      }}
    >
      <DataTableContent>
        <DataTableHead>
          <DataTableRow>
            <DataTableHeadCell hasFormControl style={{ whiteSpace: "nowrap" }}>
              <Checkbox
                {...allSelectedProps}
                onClick={() => toggleAllSelected()}
              />
            </DataTableHeadCell>
            <DataTableHeadCell style={{ width: "54.5%" }}>
              Naam
            </DataTableHeadCell>
            <DataTableHeadCell style={{ width: "54.5%" }}>
              Status
            </DataTableHeadCell>
          </DataTableRow>
        </DataTableHead>
        <DataTableBody>
          {devices.map((dev, i) => {
            return (
              <DataTableRow selected={checkboxProps[i].checked}>
                <DataTableCell hasFormControl style={{ whiteSpace: "nowrap" }}>
                  <Checkbox
                    {...checkboxProps[i]}
                    onClick={() => toggleSelectionStatus(i)}
                  />
                </DataTableCell>

                <DataTableCell style={{ width: "54.5%" }}>
                  {dev.name}
                </DataTableCell>
                <DataTableCell
                  style={{
                    display: "inline-flex",
                    alignItems: "center",
                    width: "100%",
                  }}
                >
                  <DeviceStatusIcon status={dev.status} />
                  &nbsp; {deviceStatusString(dev.status)}
                </DataTableCell>
              </DataTableRow>
            );
          })}
        </DataTableBody>
      </DataTableContent>
    </DataTable>
  );
};

export default DevicesTable;
