import React, { useCallback, useEffect, useState } from "react";
import { usePaginatedQuery } from "react-query";
import { queryCache } from "./App";
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
import { LinearProgress } from "@rmwc/linear-progress";
import { Select } from "@rmwc/select";
import { IconButton } from "@rmwc/icon-button";
import { Checkbox } from "@rmwc/checkbox";

interface GeneralTableProps<T extends object> {
  setSelectedItems: (items: T[]) => void;
  columns: Column<T>[];
  updateData?: (
    columnId: IdType<T>,
    item: T,
    value: string,
    setSkipPageReset: (skipPageReset: boolean) => void
  ) => void;
  fetchData: (page: number, pageSize: number) => Promise<Paginated<T>>;
  queryKey: string;
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
        width: "1%",
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

const GeneralTable = <T extends object>({
  setSelectedItems,
  columns,
  updateData,
  fetchData,
  queryKey,
}: GeneralTableProps<T>) => {
  const [currentData, setCurrentData] = useState({
    total: 0,
    offset: 0,
    data: [],
    maxAmount: 0,
  } as Paginated<T>);
  const [currentPageCount, setCurrentPageCount] = useState(0);

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
  } = useTable<T>(
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
      updateData: (index, id, value) =>
        updateData &&
        updateData(id, currentData.data[index], value, setSkipPageReset),
    },
    usePagination,
    useRowSelect,
    selectionHook
  );

  const fetchDevices = useCallback(
    async (key, page = 0, pageSize = 50) => {
      return fetchData(page, pageSize);
    },
    [fetchData]
  );

  const { latestData, resolvedData, error, isFetching } = usePaginatedQuery<
    Paginated<T>
  >([queryKey, pageIndex, pageSize], fetchDevices);

  React.useEffect(() => {
    if (
      latestData &&
      latestData.data &&
      latestData.offset + latestData.data.length < latestData.total
    ) {
      (async () =>
        queryCache.prefetchQuery(
          [queryKey, pageIndex + 1, pageSize],
          fetchDevices
        ))();
    }
  }, [latestData, fetchDevices, pageIndex, pageSize, queryKey]);

  useEffect(() => {
    if (resolvedData && resolvedData.data) {
      setCurrentData(resolvedData);
      setCurrentPageCount(Math.ceil(resolvedData.total / pageSize));
    }
  }, [pageIndex, pageSize, resolvedData]);

  useEffect(() => {
    setSkipPageReset(false);
  }, [resolvedData]);

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
                  <DataTableHeadCell
                    {...column.getHeaderProps()}
                    style={column.style}
                  >
                    {column.render("Header")}
                  </DataTableHeadCell>
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

export default GeneralTable;
