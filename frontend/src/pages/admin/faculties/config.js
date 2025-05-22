import {
  addFaculty,
  editFaculty,
  deleteFaculty,
  getFaculties,
} from "../../../api/admin/faculties";

import { IconButton } from "@mui/material";
import { Block } from "@mui/icons-material";

const config = {
  title: "",
  entityName: "دانشکده",
  fetchFunction: getFaculties,

  tableColumns: [
    {
      key: "name_fa",
      render: (fac) => fac.name_fa,
      customRender: (fac) =>
        !fac.is_active ? (
          <IconButton disabled>
            <Block />
          </IconButton>
        ) : null,
    },
    {
      key: "name_en",
      render: (fac) => fac.name_en,
      hideOnTablet: true,
    },
    {
      key: "short_code",
      render: (fac) => fac.short_code,
      hideOnTablet: true,
    },
    {
      key: "is_active",
      render: (fac) => (fac.is_active ? "فعال" : "غیرفعال"),
      hideOnTablet: true,
    },
  ],

  dialogFields: [
    {
      name: "name_fa",
      label: "عنوان (فارسی)",
      defaultValue: "",
      dataType: "string",
      inputType: "text",
    },
    {
      name: "name_en",
      label: "عنوان (انگلیسی)",
      dir: "ltr",
      defaultValue: "",
      dataType: "string",
      inputType: "text",
    },
    {
      name: "short_code",
      label: "کد",
      dir: "ltr",
      defaultValue: "",
      dataType: "string",
      inputType: "text",
    },
    {
      name: "is_active",
      defaultValue: true,
      dataType: "boolean",
      inputType: "radio",
      options: [
        { label: "فعال", value: true },
        { label: "غیرفعال", value: false },
      ],
    },
  ],

  validate: (data) => {
    const containsPersian = /[\u0600-\u06FF]/.test(data.name_en);
    const containsEnglish = /[A-Za-z]/.test(data.name_fa);

    return {
      name_fa: containsEnglish
        ? "عنوان فارسی نباید شامل حروف انگلیسی باشد."
        : "",
      name_en: containsPersian
        ? "عنوان انگلیسی نباید شامل حروف فارسی باشد."
        : "",
    };
  },

  onSubmit: async ({ id, data, mode }) => {
    if (mode === "add") await addFaculty(data);
    if (mode === "edit") await editFaculty(id, data);
    if (mode === "delete") await deleteFaculty(id);
  },
};

export default config;
