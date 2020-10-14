import { CellProps, IdType } from "react-table";
import React, { ChangeEvent, useEffect, useState } from "react";
import { Device } from "./DevicesTable";

export interface EditableCellProps<T extends object> extends CellProps<T> {
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

export default EditableCell;
