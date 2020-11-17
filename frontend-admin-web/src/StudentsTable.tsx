import React, { useMemo } from "react";
import { Icon } from "@rmwc/icon";
import { useMutation } from "react-query";
import { fetchAuthnd } from "./DevicesPage";
import "@rmwc/linear-progress/styles";
import { Column, IdType } from "react-table";
import { queryCache } from "./App";
import GeneralTable from "./GeneralTable";
import EditableCell from "./EditableCell";
import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import "@rmwc/dialog/styles";
import "@rmwc/textfield/styles";
import { dialogQueue } from "./dialogQueue";

export interface Student {
  id: string;
  zermelo: StudentZermeloInfo;
  hasCardId: boolean;
}

export interface StudentZermeloInfo {
  user: string;
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

interface StudentsTableProps {
  setSelectedItems: (items: Student[]) => void;
}

interface StudentPatch {
  id: string;
  zermelo?: StudentZermeloInfo;
  cardId?: string;
}

const updateStudent = (patch: StudentPatch) =>
  fetchAuthnd(`student/${patch.id}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });

const boolToYesNoStringDutch = (b: boolean) => (b ? "Ja" : "Nee");

const StudentsTable: React.FC<StudentsTableProps> = ({ setSelectedItems }) => {
  const [updateStudentMut] = useMutation(updateStudent, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("students");
    },
  });

  const columns = useMemo<Array<Column<Student>>>(
    () => [
      {
        id: "zermeloUser",
        Header: "Zermelo-gebruiker",
        accessor: (student) => student.zermelo.user,
        Cell: EditableCell,
      },
      {
        id: "hasCardCode",
        Header: "Heeft pascode",
        accessor: (student) => (
          <div style={{ display: "flex", alignItems: "center" }}>
            <UserHasCardCodeIcon hasCardCode={student.hasCardId} />
            &nbsp;
            {boolToYesNoStringDutch(student.hasCardId)}
            &nbsp;&nbsp;&nbsp;
            <Theme use={"onSurface"} wrap>
              <Button
                onClick={() => {
                  dialogQueue
                    .prompt({
                      title: "Pascode toewijzen",
                      body: (
                        <>
                          Deze zal toegewezen worden aan <br />
                          de leerling met als Zermelo-gebruiker
                          <br />
                          <code style={{ fontWeight: "bold" }}>
                            {student.zermelo.user}
                          </code>
                        </>
                      ),
                      acceptLabel: "Toewijzen",
                      cancelLabel: "Annuleren",
                      inputProps: {
                        outlined: true,
                      },
                    })
                    .then((res) => {
                      return (
                        res &&
                        updateStudentMut({
                          id: student.id,
                          cardId: res,
                        })
                      );
                    });
                }}
              >
                Toewijzen
              </Button>
            </Theme>
          </div>
        ),
      },
    ],
    [updateStudentMut]
  );

  const updateData = async (
    columnId: IdType<Student>,
    item: Student,
    value: string,
    setSkipPageReset: (skipPageReset: boolean) => void
  ) => {
    if (columnId === "zermeloUser") {
      setSkipPageReset(true);

      updateStudentMut({
        id: item.id,
        zermelo: {
          user: value,
        },
      })
        .then()
        .catch();
    }
  };

  const fetchData = (page: number, pageSize: number) => {
    return fetchAuthnd(
      `student?offset=${page * pageSize}&maxAmount=${pageSize}`
    ).then((res) => res.json());
  };

  return (
    <>
      <GeneralTable
        setSelectedItems={setSelectedItems}
        columns={columns}
        fetchData={fetchData}
        queryKey={"students"}
        updateData={updateData}
      />
    </>
  );
};

export default StudentsTable;
