import { Button } from "@rmwc/button";
import { Theme } from "@rmwc/theme";
import { Elevation } from "@rmwc/elevation";
import React, { useState } from "react";
import { useMutation } from "react-query";
import { queryCache } from "./App";
import { fetchAuthnd } from "./DevicesPage";
import StudentsTable, { Student } from "./StudentsTable";

const removeStudent = (students: Student[]) =>
  fetchAuthnd(`/api/student`, {
    method: "DELETE",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      studentIds: students.map((student) => student.id),
    }),
  });

const DevicesPage: React.FC = () => {
  const [selectedItems, setSelectedItems] = useState([] as Student[]);

  const [deleteStudents] = useMutation(removeStudent, {
    onSuccess: async () => {
      await queryCache.invalidateQueries("students");
    },
  });

  const onDeleteStudents = async () => {
    try {
      await deleteStudents(selectedItems);
    } catch (error) {}
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
        height: "100%",
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
        <h1 style={{ marginTop: 0 }}>Leerlingen</h1>
        <div>
          <Button
            icon={"delete"}
            danger
            raised
            disabled={selectedItems.length === 0}
            onClick={() => onDeleteStudents()}
          >
            Verwijderen
          </Button>
        </div>
      </div>

      <Theme use={"background"} wrap>
        <Elevation
          z={16}
          style={{
            flexGrow: 1,
            margin: 16,
            borderRadius: 4,
            height: "100%",
            overflow: "hidden",
          }}
        >
          <StudentsTable setSelectedItems={setSelectedItems} />
        </Elevation>
      </Theme>
    </div>
  );
};

export default DevicesPage;
