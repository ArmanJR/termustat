import {
  addSemester,
  editSemester,
  deleteSemester,
} from "../../../api/admin/semesters";

const config = {
  title: "نیمسال تحصیلی",
  entityName: "نیمسال تحصیلی",

  tableColumns: [
    {
      key: "term",
      render: (sem) => sem.term,
    },
    {
      key: "year",
      render: (sem) => sem.year,
    },
  ],

  dialogFields: [
    {
      name: "year",
      label: "سال تحصیلی",
      dir: "ltr",
      defaultValue: 0,
      dataType: "number",
      inputType: "number",
    },
    {
      name: "term",
      dir: "ltr",
      defaultValue: "spring",
      dataType: "string",
      inputType: "radio",
      options: [
        { label: "spring", value: "spring" },
        { label: "fall", value: "fall" },
      ],
    },
  ],

  onSubmit: async ({ id, data, mode }) => {
    if (mode === "add") await addSemester(data);
    if (mode === "edit") await editSemester(id, data);
    if (mode === "delete") await deleteSemester(id);
  },
};

export default config;
