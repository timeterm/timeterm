import React, { useMemo } from "react";
import { Icon } from "@rmwc/icon";
import { useMutation } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import "@rmwc/linear-progress/styles";
import { Column, IdType } from "react-table";
import { queryCache } from "./App";
import GeneralTable from "./GeneralTable";
import EditableCell from "./EditableCell";

export interface User {
  id: string;
  zermeloUser: string;
  hasCardCode: boolean;
}

interface PrimaryDeviceStatusIconProps {
  hasCardCode: boolean;
}

const UserHasCardCodeIcon: React.FC<PrimaryDeviceStatusIconProps> = ({
  hasCardCode,
}) => {
  return (
    <Icon
      icon={hasCardCode ? "check_circle" : "warning"}
      style={{
        color: hasCardCode ? "#4ecd6a" : "#ffab00",
      }}
    />
  );
};

interface UsersTableProps {
  setSelectedItems: (items: User[]) => void;
}

interface UserPatch {
  id: string;
  zermeloUser: string;
}

const updateUser = (patch: UserPatch) =>
  fetchAuthnd(`/api/user/${patch.id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });

const boolToYesNoStringDutch = (b: boolean) => (b ? "Ja" : "Nee");

const UsersTable: React.FC<UsersTableProps> = ({ setSelectedItems }) => {
  const [updateUserMut] = useMutation(updateUser, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("users");
    },
  });

  const columns = useMemo<Array<Column<User>>>(
    () => [
      {
        id: "zermeloUser",
        Header: "Zermelo-gebruiker",
        accessor: (user) => user.zermeloUser,
        Cell: EditableCell,
      },
      {
        id: "hasCardCode",
        Header: "Heeft pascode",
        accessor: (user) => (
          <div style={{ display: "flex", alignItems: "center" }}>
            <UserHasCardCodeIcon hasCardCode={user.hasCardCode} />
            &nbsp;
            {boolToYesNoStringDutch(user.hasCardCode)}
          </div>
        ),
      },
    ],
    []
  );

  const updateData = async (
    columnId: IdType<User>,
    item: User,
    value: string,
    setSkipPageReset: (skipPageReset: boolean) => void
  ) => {
    if (columnId === "zermeloUser") {
      setSkipPageReset(true);

      updateUserMut({
        id: item.id,
        zermeloUser: value,
      })
        .then()
        .catch();
    }
  };

  const fetchData = (page: number, pageSize: number) => {
    return fetchAuthnd(
      `/api/user?offset=${page * pageSize}&maxAmount=${pageSize}`
    ).then((res) => res.json());
  };

  return (
    <GeneralTable
      setSelectedItems={setSelectedItems}
      columns={columns}
      fetchData={fetchData}
      queryKey={"users"}
      updateData={updateData}
    />
  );
};

export default UsersTable;
