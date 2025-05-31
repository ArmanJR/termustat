import { useState, useEffect } from "react";
import styles from "./EntityDialog.module.css";
import Input from "../form/Input";
import Button from "../form/Button";
import { Dialog, DialogContent, Snackbar, Alert } from "@mui/material";

const EntityDialog = ({
  open,
  mode,
  entityName,
  entity,
  fields,
  validate,
  onSubmit,
  onClose,
}) => {
  const modeLabels = { add: "افزودن", edit: "ویرایش", delete: "حذف" };
  const [formData, setFormData] = useState({});
  const [errors, setErrors] = useState({});

  useEffect(() => {
    setErrors({});
    if (mode === "edit" || mode === "delete") {
      const initialData = {};
      fields.forEach((field) => {
        initialData[field.name] = entity[field.name];
      });
      setFormData(initialData);
    } else if (mode === "add") {
      const emptyData = {};
      fields.forEach((f) => {
        emptyData[f.name] = f.defaultValue;
      });
      setFormData(emptyData);
    }
  }, [mode]);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]:
        fields.find((field) => field.name === e.target.name).dataType === "boolean"
          ? e.target.value === "true"
            ? true
            : false
          : fields.find((field) => field.name === e.target.name).dataType === "number"
          ? parseInt(e.target.value, 10)
          : e.target.value,
    });
  };

  const handleSelectChange = (e, field) => {
    const value = e.target.value;
    setFormData((prev) => {
      const hasFacultyField = fields.some((f) => f.name === "faculty_id");
      return {
        ...prev,
        [field.name]: value,
        ...(field.name === "university_id" && hasFacultyField
          ? { faculty_id: "" }
          : {}),
      };
    });
    if (field.onChangeHandler) {
      field.onChangeHandler(value);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const validationErrors = validate ? validate(formData) : {};
    setErrors(validationErrors);
    if (Object.values(validationErrors).some((msg) => msg)) {
      return;
    }

    try {
      await onSubmit({ id: entity?.id, data: formData, mode });
      showSnackbar(`${modeLabels[mode]} ${entityName} با موفقیت انجام شد`, "success");
      onClose();
    } catch (error) {
      if ((entityName === "دانشگاه" || entityName === "نیمسال تحصیلی") && error.response?.status === 409)
        showSnackbar(`این ${entityName} قبلا در سیستم ثبت شده است.`, "warning");
      else if (entityName === "دانشکده" && error.response?.status === 409)
        showSnackbar("کد دانشکده تکراری است.", "warning");
      else
        showSnackbar(`خطا در ${modeLabels[mode]} ${entityName}. لطفا دوباره تلاش کنید.`, "error");
    }
  };

  const [snackbar, setSnackbar] = useState({
    open: false,
    message: "",
    severity: "success",
  });

  const showSnackbar = (message, severity = "success") => {
    setSnackbar({ open: true, message, severity });
  };

  return (
    <>
      <Dialog open={open} onClose={onClose}>
        <h1>
          <span style={{ color: "#309a9a" }}>&#9699; &nbsp;</span>
          {`${modeLabels[mode]} ${entityName}`}
        </h1>
        <DialogContent>
          <form onSubmit={handleSubmit} className={styles.form}>
            {mode === "edit" || mode === "add" ? (
              fields.map((field) =>
                field.inputType === "radio" ? (
                  <div
                    dir={field.dir || "rtl"}
                    style={(field.dir === "ltr" && { textAlign: "left" }) || {}}
                  >
                    {field.options?.map((option) => (
                      <label>
                        <input
                          type="radio"
                          name={field.name}
                          value={option.value}
                          checked={formData[field.name]?.toString() === option.value.toString()}
                          onChange={handleChange}
                        />
                        {option.label}
                      </label>
                    ))}
                  </div>
                ) : field.inputType === "text" ||
                  field.inputType === "number" ? (
                  <>
                    <Input
                      type={field.inputType}
                      name={field.name}
                      label={field.label}
                      value={formData[field.name] || ""}
                      onChange={handleChange}
                      dir={field.dir || ""}
                      required
                    />
                    {errors[field.name] && (
                      <div className={styles.errorMessage}>
                        {errors[field.name]}
                      </div>
                    )}
                  </>
                ) : field.inputType === "select" ? (
                  <div key={field.name}>
                    <select
                      name={field.name}
                      value={formData[field.name] || ""}
                      onChange={(e) => handleSelectChange(e, field)}
                    >
                      <option key="" value="" disabled>
                        {field.label}
                      </option>
                      {field.options?.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                  </div>
                ) : (
                  <></>
                )
              )
            ) : (
              <p>
                آیا از حذف
                <>
                  &nbsp;{entityName}
                  <span style={{ fontWeight: "bold" }}>
                    &nbsp;"
                    {entityName === "دانشگاه" || entityName === "دانشکده"
                      ? formData.name_fa
                      : entityName === "نیمسال تحصیلی"
                      ? `${formData.term} ${formData.year}`
                      : entityName === "کاربر"
                      ? `${formData.first_name} ${formData.last_name}`
                      : ""}
                    "&nbsp;
                  </span>
                </>
                مطمئن هستید؟
              </p>
            )}
            <Button type="submit" value="تأیید" />
            <Button onClick={onClose} value="انصراف" className={styles.whiteButton} />
          </form>
        </DialogContent>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          severity={snackbar.severity}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </>
  );
};

export default EntityDialog;
