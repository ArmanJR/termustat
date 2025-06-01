import styles from "./Users.module.css";
import { editUser, deleteUser } from "../../../api/admin/users";
import { getFaculties } from "../../../api/admin/faculties";
import { getUniversities } from "../../../api/admin/universities";
import { IconButton } from "@mui/material";
import { Error } from '@mui/icons-material';

const universityListCache = { data: null, fetched: false };
const facultyListCache = new Map();

const getCachedUniversities = async () => {
  if (universityListCache.fetched && universityListCache.data) {
    return universityListCache.data;
  }
  try {
    const res = await getUniversities();
    const options = res.data.map((u) => ({
      label: u.name_fa,
      value: u.id,
    }));
    universityListCache.data = options;
    universityListCache.fetched = true;
    return options;
  } catch(err) {
    console.error("Failed to fetch universities:", err);
    return [];
  }
};

const getCachedFaculties = async (universityId) => {
  if (facultyListCache.has(universityId)) {
    return facultyListCache.get(universityId);
  }
  try {
    const res = await getFaculties(universityId);
    const options = res.data.map((f) => ({
      label: f.name_fa,
      value: f.id,
    }));
    facultyListCache.set(universityId, options);
    return options;
  } catch(err) {
    console.error(`Failed to fetch faculties for university ID ${universityId}:`, err);
    return [];
  }
};

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
      key: "gender",
      render: (user) => (user.gender === "female" ? "زن" : "مرد"),
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
    {
      key: "student_id",
      render: (user) => user.student_id,
      hideOnTablet: true,
    },
    {
      key: "email",
      render: (user) => (
        <div  className={styles.email}>
          {user.email}
          {!user.email_verified && (
            <IconButton disabled>
              <Error />
            </IconButton>
          )}
        </div>
      ),
      hideOnTablet: true,
    },
    {
      key: "is_admin",
      render: (user) => (
        <div className={styles.admin}>
          {user.is_admin ? "ادمین" : ""}
        </div>
      ),
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
    {
      name: "password",
      label: "رمز عبور جدید (اختیاری)",
      dir: "ltr",
      defaultValue: "",
      dataType: "string",
      inputType: "password",
      required: false,
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
        universityOptions = await getCachedUniversities();
        if (entity?.university_id) {
          facultyOptions = await getCachedFaculties(entity.university_id);
        }
      }
      const handleUniversityChange = async (universityId) => {
        setDialogFields((prev) =>
          prev.map((field) =>
            field.name === "faculty_id" ? { ...field, options: [] } : field
          )
        );
        if (universityId) {
          const newFacultyOptions = await getCachedFaculties(universityId);
          setDialogFields((prev) =>
            prev.map((field) =>
              field.name === "faculty_id"
                ? {
                    ...field,
                    options: newFacultyOptions,
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
              options: universityOptions,
              onChangeHandler: handleUniversityChange,
            };
          }
          if (field.name === "faculty_id") {
            return {
              ...field,
              options: facultyOptions,
            };
          }
          return field;
        })
      );
    } catch(err) {
      console.error("Failed to open dialog and set fields:", err);
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
