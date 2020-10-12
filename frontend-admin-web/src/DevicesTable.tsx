import React, {
  ChangeEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
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
import { useMutation, usePaginatedQuery } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import { LinearProgress } from "@rmwc/linear-progress";
import "@rmwc/linear-progress/styles";
import {
  CellProps,
  Column,
  HeaderProps,
  Hooks,
  IdType,
  usePagination,
  useRowSelect,
  useTable,
} from "react-table";
import { queryCache } from "./App";
import { Theme } from "@rmwc/theme";

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

interface EditableCellProps<T extends object> extends CellProps<T> {
  updateData: (index: number, id: IdType<T>, data: string) => void;
}

const EditableCell: React.FC<EditableCellProps<Device>> = ({
  value: initialValue,
  row: { index },
  column: { id },
  updateData,
}) => {
  const [value, setValue] = useState(initialValue);

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  const onBlur = () => {
    updateData(index, id, value);
  };

  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  return (
    <input
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      style={{
        padding: 0,
        margin: 0,
        border: 0,
        width: "100%",
        color: "inherit",
        fontSize: "inherit",
        fontWeight: "inherit",
        letterSpacing: "inherit",
        textTransform: "inherit",
        fontFamily: "inherit",
        background: "inherit",
      }}
    />
  );
};

interface DevicePatch {
  id: string;
  name?: string;
}

const updateDevice = (patch: DevicePatch) =>
  fetchAuthnd(`/api/device/${patch.id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });

const DevicesTable: React.FC<DevicesTableProps> = ({ setSelectedItems }) => {
  const [currentData, setCurrentData] = useState({
    total: 0,
    offset: 0,
    data: [],
    maxAmount: 0,
  } as Paginated<Device>);
  const [currentPageCount, setCurrentPageCount] = useState(0);

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
            <DeviceStatusIcon status={dev.status} />
            &nbsp;
            {deviceStatusString(dev.status)}
          </div>
        ),
      },
    ],
    []
  );

  const [skipPageReset, setSkipPageReset] = useState(false);

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
  } = useTable<Device>(
    {
      columns: columns,
      data: currentData.data,
      manualPagination: true,
      autoResetPage: !skipPageReset,
      initialState: {
        pageIndex: 0,
        pageSize: 50,
      },
      pageCount: currentPageCount,
      updateData: async (rowIndex, columnId, value) => {
        if (columnId === "name") {
          setSkipPageReset(true);

          try {
            await updateDeviceMut({
              id: currentData.data[rowIndex].id,
              name: value,
            });
          } catch (e) {}
        }
      },
    },
    usePagination,
    useRowSelect,
    selectionHook
  );

  useEffect(() => {
    setSkipPageReset(false);
  }, [currentData]);

  const fetchDevices = useCallback(async (key, page = 0, pageSize = 50) => {
    return fetchAuthnd(
      `/api/device?offset=${page * pageSize}&maxAmount=${pageSize}`
    ).then((res) => res.json());
  }, []);

  const { latestData, resolvedData, error, isFetching } = usePaginatedQuery<
    Paginated<Device>
  >(["devices", pageIndex, pageSize], fetchDevices);

  React.useEffect(() => {
    if (
      latestData &&
      latestData.offset + latestData.data.length < latestData.total
    ) {
      (async () => {
        await queryCache.prefetchQuery(
          ["devices", pageIndex + 1],
          fetchDevices
        );
      })();
    }
  }, [latestData, fetchDevices, pageIndex]);

  useEffect(() => {
    if (resolvedData && resolvedData.data) {
      setCurrentData(resolvedData);
      setCurrentPageCount(Math.ceil(resolvedData.total / pageSize));
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
        flexWrap: "nowrap",
        justifyContent: "space-between",
        height: "100%",
        overflow: "hidden",
      }}
    >
      <DataTable
        {...getTableProps()}
        style={{
          width: "100%",
          borderRadius: "4px 4px 0 0",
          borderTop: 0,
          borderLeft: 0,
          borderRight: 0,
          overflowY: "scroll",
          flex: "1 1 auto",
        }}
      >
        <DataTableContent>
          <DataTableHead>
            {headerGroups.map((headerGroup) => (
              <DataTableRow {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map((column) => (
                  <Theme use={"surface"} wrap>
                    <DataTableHeadCell
                      {...column.getHeaderProps()}
                      style={column.style}
                    >
                      {column.render("Header")}
                    </DataTableHeadCell>
                  </Theme>
                ))}
              </DataTableRow>
            ))}
            <DataTableRow style={{ padding: 0, margin: 0 }}>
              <DataTableHeadCell
                colSpan={10000}
                style={{ padding: 0, margin: 0, height: 0 }}
              >
                {error ? (
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "center",
                      padding: 4,
                    }}
                  >
                    Er is een fout opgetreden bij het ophalen van de data
                  </div>
                ) : (
                  <LinearProgress
                    style={{
                      maxHeight: isFetching ? 4 : 0,
                      transition: "0.5s",
                    }}
                  />
                )}
              </DataTableHeadCell>
            </DataTableRow>
          </DataTableHead>
          <DataTableBody {...getTableBodyProps()} style={{ overflow: "auto" }}>
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
          flex: "0 0 auto",
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
