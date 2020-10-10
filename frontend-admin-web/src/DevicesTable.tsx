import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Icon } from "@rmwc/icon";
import { Checkbox } from "@rmwc/checkbox";
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
import { usePaginatedQuery } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import { LinearProgress } from "@rmwc/linear-progress";
import "@rmwc/linear-progress/styles";
import {
  CellProps,
  Column,
  HeaderProps,
  Hooks,
  usePagination,
  useRowSelect,
  useTable,
} from "react-table";

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

interface DevicesTableProps {
  setSelectedItems: (items: Device[]) => void;
}

export interface Paginated<T> {
  offset: number;
  maxAmount: number;
  total: number;
  data: T[];
}

const selectionHook = <T extends {}>(hooks: Hooks<T>) => {
  hooks.visibleColumns.push((columns) => [
    {
      id: "selection",
      style: {
        width: 0,
      },
      Header: ({ getToggleAllRowsSelectedProps }: HeaderProps<T>) => (
        <div>
          <Checkbox {...getToggleAllRowsSelectedProps()} />
        </div>
      ),
      Cell: ({ row }: CellProps<T>) => (
        <div>
          <Checkbox {...row.getToggleRowSelectedProps()} />
        </div>
      ),
    },
    ...columns,
  ]);
};

const DevicesTable: React.FC<DevicesTableProps> = ({ setSelectedItems }) => {
  const [currentData, setCurrentData] = useState({
    total: 0,
    offset: 0,
    data: [],
    maxAmount: 0,
  } as Paginated<Device>);
  const [currentPageCount, setCurrentPageCount] = useState(0);

  const columns = useMemo<Array<Column<Device>>>(
    () => [
      {
        Header: "Naam",
        accessor: (dev) => dev.name,
      },
      {
        Header: "Status",
        accessor: (dev) => (
          <div style={{ display: "flex", alignItems: "center" }}>
            <DeviceStatusIcon status={dev.status} />
            &nbsp;
            {deviceStatusString(dev.status)}
          </div>
        ),
      },
    ],
    []
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    page,
    canPreviousPage,
    canNextPage,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    setPageSize,
    selectedFlatRows,
    state: { pageIndex, pageSize, selectedRowIds },
  } = useTable(
    {
      columns: columns,
      data: currentData.data,
      manualPagination: true,
      initialState: {
        pageIndex: 0,
        pageSize: 50,
      },
      pageCount: currentPageCount,
    },
    usePagination,
    useRowSelect,
    selectionHook
  );

  const fetchDevices = useCallback(async (key, page = 0, pageSize = 50) => {
    return fetchAuthnd(
      `/api/device?offset=${page * pageSize}&maxAmount=${pageSize}`
    ).then((res) => res.json());
  }, []);

  const { resolvedData, error, isFetching } = usePaginatedQuery<
    Paginated<Device>
  >(["organizationDevices", pageIndex, pageSize], fetchDevices);

  useEffect(() => {
    if (resolvedData && resolvedData.data) {
      setCurrentData(resolvedData);
      setCurrentPageCount(Math.ceil(resolvedData.maxAmount / pageSize));
    }
  }, [pageIndex, pageSize, resolvedData]);

  useEffect(() => {
    setSelectedItems(selectedFlatRows.map((row) => row.original));
  }, [selectedFlatRows, setSelectedItems]);

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
        {...getTableProps()}
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
            {headerGroups.map((headerGroup) => (
              <DataTableRow {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map((column) => (
                  <DataTableHeadCell
                    {...column.getHeaderProps()}
                    style={(column as any).style}
                  >
                    {column.render("Header")}
                  </DataTableHeadCell>
                ))}
              </DataTableRow>
            ))}
          </DataTableHead>
          <DataTableBody {...getTableBodyProps()}>
            <tr>
              <td colSpan={10000} style={{ padding: 0 }}>
                <LinearProgress
                  style={{ maxHeight: isFetching ? 4 : 0, transition: "0.5s" }}
                />
              </td>
            </tr>
            {page.map((row, i) => {
              prepareRow(row);
              return (
                <DataTableRow
                  {...row.getRowProps()}
                  selected={selectedRowIds[i]}
                >
                  {row.cells.map((cell) => {
                    return (
                      <DataTableCell {...cell.getCellProps()}>
                        {cell.render("Cell")}
                      </DataTableCell>
                    );
                  })}
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
          onChange={(e) => {
            setPageSize(Number((e.target as HTMLInputElement).value));
          }}
        />
        <span style={{ marginLeft: 48, marginRight: 48 }}>
          {page.length > 0 ? pageIndex * pageSize + 1 : 0} -{" "}
          {pageIndex * pageSize + page.length} van de {currentData?.total || 0}
        </span>
        <IconButton
          icon={"first_page"}
          onClick={() => gotoPage(0)}
          disabled={!canPreviousPage}
        />
        <IconButton
          icon={"chevron_left"}
          onClick={() => previousPage()}
          disabled={!canPreviousPage}
        />
        <IconButton
          icon={"chevron_right"}
          onClick={() => nextPage()}
          disabled={!canNextPage}
        />
        <IconButton
          icon={"last_page"}
          onClick={() => gotoPage(pageCount - 1)}
          disabled={!canNextPage}
        />
      </div>
    </div>
  );
};

export default DevicesTable;
