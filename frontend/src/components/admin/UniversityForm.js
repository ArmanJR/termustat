import { useState, useEffect } from "react";
import Input from "../form/Input";
import Button from "../form/Button";
import styles from "./UniversityForm.module.css";

import { Dialog, DialogContent } from "@mui/material";

const UniversityForm = ({ open, handleClose, university, mode }) => {
  const [formData, setFormData] = useState({
    name_fa: "",
    name_en: "",
    is_active: true,
  });
  const [error, setError] = useState({
    name_fa: "",
    name_en: "",
  });

  const containsPersian = (text) => /[\u0600-\u06FF]/.test(text);
  const containsEnglish = (text) => /[A-Za-z]/.test(text);

  useEffect(() => {
    setError({
      name_fa: "",
      name_en: "",
    });
    if ((mode == "edit" || mode == "delete") && university) {
      setFormData({
        name_fa: university.name_fa,
        name_en: university.name_en,
        is_active: university.is_active,
      });
    } else if (mode === "add") {
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

  const handleSubmit = (e) => {
    e.preventDefault();
  
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
  
    alert(JSON.stringify(formData));
    handleClose();
  };
  

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      sx={{
        "& .MuiDialog-paper": {
          width: "100%",
          maxWidth: "460px",
        },
      }}
    >
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
  );
};

export default UniversityForm;
