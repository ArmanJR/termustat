import { useState, useEffect } from "react";
import Input from "../form/Input";
import Button from "../form/Button";
import styles from "./UniversityForm.module.css";
import UniversityFormTheme from "./UniversityFormTheme";
import { Dialog, DialogContent, ThemeProvider, Snackbar, Alert } from "@mui/material";
import { addUniversity, deleteUniversity, editUniversity } from "../../api/admin/universities";

const UniversityForm = ({ open, handleClose, university, mode, refetchUniversities }) => {
  const [id, setId] = useState(null);

  const [formData, setFormData] = useState({
    name_fa: "",
    name_en: "",
    is_active: true,
  });

  const [error, setError] = useState({
    name_fa: "",
    name_en: "",
  });

  const [snackbar, setSnackbar] = useState({
    open: false,
    message: "",
    severity: "success",
  });

  const showSnackbar = (message, severity = "success") => {
    setSnackbar({
      open: true,
      message,
      severity,
    });
  };

  useEffect(() => {
    setError({
      name_fa: "",
      name_en: "",
    });
    if ((mode == "edit" || mode == "delete") && university) {
      setId(university.id);
      setFormData({
        name_fa: university.name_fa,
        name_en: university.name_en,
        is_active: university.is_active,
      });
    } else if (mode === "add") {
      setId(null);
      setFormData({
        name_fa: "",
        name_en: "",
        is_active: true,
      });
    }
  }, [university, mode, open]);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]:
        e.target.name == "is_active"
          ? e.target.value == "true"
            ? true
            : false
          : e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
  
    const containsPersian = (text) => /[\u0600-\u06FF]/.test(text);
    const containsEnglish = (text) => /[A-Za-z]/.test(text);

    const newErrors = {
      name_fa: "",
      name_en: "",
    };
  
    if (containsEnglish(formData.name_fa)) {
      newErrors.name_fa = "عنوان فارسی نباید شامل حروف انگلیسی باشد.";
    }
  
    if (containsPersian(formData.name_en)) {
      newErrors.name_en = "عنوان انگلیسی نباید شامل حروف فارسی باشد.";
    }
  
    setError(newErrors);
  
    if (newErrors.name_fa || newErrors.name_en) {
      return;
    }
  
    try {
      if (mode == "add") {
        await addUniversity(formData);
        showSnackbar("دانشگاه با موفقیت افزوده شد!", "success");
      } else if (mode == "delete") {
        await deleteUniversity(id);
        showSnackbar("دانشگاه با موفقیت حذف شد!", "success");
      } else if (mode == "edit") {
        await editUniversity(id, formData);
        showSnackbar("دانشگاه با موفقیت ویرایش شد!", "success");
      }
      refetchUniversities();
      handleClose();
    } catch (error) {
      if (error.response?.status === 409)
        showSnackbar("این دانشگاه قبلا در سیستم ثبت شده است.", "warning");
      else
        showSnackbar("خطا در افزودن دانشگاه. لطفا دوباره تلاش کنید.", "error");
    }
  };
  
  return (
    <ThemeProvider theme={UniversityFormTheme}>
      <Dialog open={open} onClose={handleClose}>
        <h1>
          <span style={{ color: "#309a9a" }}>&#9699; &nbsp;</span>
          {mode == "edit"
            ? "ویرایش دانشگاه"
            : mode == "delete"
            ? "حذف دانشگاه"
            : "دانشگاه جدید"}
        </h1>
        <DialogContent>
          <form onSubmit={handleSubmit} className={styles.form}>
            {mode == "edit" || mode == "add" ? (
              <>
                <Input
                  type="text"
                  name="name_fa"
                  label="عنوان (فارسی)"
                  value={formData.name_fa}
                  onChange={handleChange}
                  required
                />
                <div className={styles.errorMessage}>{error.name_fa}</div>

                <Input
                  type="text"
                  name="name_en"
                  label="عنوان (انگلیسی)"
                  value={formData.name_en}
                  dir="ltr"
                  onChange={handleChange}
                  required
                />
                <div className={styles.errorMessage}>{error.name_en}</div>

                <div>
                  <label>
                    <input
                      type="radio"
                      name="is_active"
                      value="true"
                      checked={formData.is_active + "" === "true"}
                      onChange={handleChange}
                      required
                    />
                    فعال
                  </label>
                  <label>
                    <input
                      type="radio"
                      name="is_active"
                      value="false"
                      checked={formData.is_active + "" === "false"}
                      onChange={handleChange}
                    />
                    غیرفعال
                  </label>
                </div>
              </>
            ) : (
              <p>
                آیا از حذف دانشگاه
                <span style={{ fontWeight: "bold" }}>
                  &nbsp;"{formData.name_fa}"&nbsp;
                </span>
                مطمئن هستید؟
              </p>
            )}
            <Button type="submit" value="تأیید" />
            <Button
              onClick={handleClose}
              value="انصراف"
              className={styles.whiteButton}
            />
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
    </ThemeProvider>
  );
};

export default UniversityForm;
