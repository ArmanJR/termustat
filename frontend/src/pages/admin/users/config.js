import { editUser, deleteUser } from "../../../api/admin/users";
import { getFaculties } from "../../../api/admin/faculties";
import { getUniversities } from "../../../api/admin/universities";

const config = {
  title: "کاربران",
  entityName: "کاربر",

  tableColumns: [
    {
      key: "first_name",
      render: (user) => user.first_name,
    },
    {
      key: "last_name",
      render: (user) => user.last_name,
    },
    {
      key: "email",
      render: (user) => user.email,
      hideOnTablet: true,
    },
    {
      key: "email_verified",
      render: (user) => user.email_verified + "",
      hideOnTablet: true,
    },
    {
      key: "gender",
      render: (user) => user.gender,
      hideOnTablet: true,
    },
    {
      key: "is_admin",
      render: (user) => user.is_admin + "",
      hideOnTablet: true,
    },
    {
      key: "student_id",
      render: (user) => user.student_id,
      hideOnTablet: true,
    },
    {
      key: "university_name",
      render: (user) => user.university_name,
      hideOnTablet: true,
    },
    {
      key: "faculty_name",
      render: (user) => user.faculty_name,
      hideOnTablet: true,
    },
  ],

  dialogFields: [
    {
      name: "first_name",
      label: "نام",
      defaultValue: "",
      dataType: "string",
      inputType: "text",
    },
    {
      name: "last_name",
      label: "نام خانوادگی",
      defaultValue: "",
      dataType: "string",
      inputType: "text",
    },
    {
      name: "gender",
      defaultValue: "female",
      dataType: "string",
      inputType: "radio",
      options: [
        { label: "زن", value: "female" },
        { label: "مرد", value: "male" },
      ],
    },
    {
      name: "university_id",
      label: "دانشگاه",
      defaultValue: "",
      dataType: "string",
      inputType: "select",
      options: [],
    },
    {
      name: "faculty_id",
      label: "دانشکده",
      defaultValue: "",
      dataType: "string",
      inputType: "select",
      options: [],
    },
  ],

  onSubmit: async ({ id, data, mode }) => {
    if (mode === "edit") await editUser(id, data);
    if (mode === "delete") await deleteUser(id);
  },

  onOpenDialog: async ({ mode, entity, setDialogFields }) => {
    try {
      let universityOptions = [];
      let facultyOptions = [];
      if (mode === "edit") {
        const universityRes = await getUniversities();
        universityOptions = universityRes?.data || [];
        if (entity?.university_id) {
          const facultyRes = await getFaculties(entity.university_id);
          facultyOptions = facultyRes?.data || [];
        }
      }
      const handleUniversityChange = async (universityId) => {
        setDialogFields((prev) =>
          prev.map((field) =>
            field.name === "faculty_id" ? { ...field, options: [] } : field
          )
        );
        if (universityId) {
          const facultyRes = await getFaculties(universityId);
          const newFacultyOptions = facultyRes?.data || [];
          setDialogFields((prev) =>
            prev.map((field) =>
              field.name === "faculty_id"
                ? {
                    ...field,
                    options: newFacultyOptions.map((f) => ({
                      label: f.name_fa,
                      value: f.id,
                    })),
                  }
                : field
            )
          );
        }
      };
      setDialogFields((prevFields) =>
        prevFields.map((field) => {
          if (field.name === "university_id") {
            return {
              ...field,
              options: universityOptions.map((u) => ({
                label: u.name_fa,
                value: u.id,
              })),
              onChangeHandler: handleUniversityChange,
            };
          }
          if (field.name === "faculty_id") {
            return {
              ...field,
              options: facultyOptions.map((f) => ({
                label: f.name_fa,
                value: f.id,
              })),
            };
          }
          return field;
        })
      );
    } catch {
      setDialogFields((prevFields) =>
        prevFields.map((field) =>
          field.name === "university_id" || field.name === "faculty_id"
            ? { ...field, options: [] }
            : field
        )
      );
    }
  },
};

export default config;
