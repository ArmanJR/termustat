import {
  addUniversity,
  editUniversity,
  deleteUniversity,
} from "../../../api/admin/universities";

import { IconButton } from "@mui/material";
import { Block } from "@mui/icons-material";

const config = {
  title: "دانشگاه‌ها",
  entityName: "دانشگاه",

  tableColumns: [
    {
      key: "name_fa",
      render: (uni) => uni.name_fa,
      customRender: (uni) =>
        !uni.is_active ? (
          <IconButton disabled>
            <Block />
          </IconButton>
        ) : null,
    },
    {
      key: "name_en",
      render: (uni) => uni.name_en,
      hideOnTablet: true,
    },
    {
      key: "is_active",
      render: (uni) => (uni.is_active ? "فعال" : "غیرفعال"),
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
    if (mode === "add") await addUniversity(data);
    if (mode === "edit") await editUniversity(id, data);
    if (mode === "delete") await deleteUniversity(id);
  },
};

export default config;
