import styles from "./Input.module.css";
import { useState, useEffect } from "react";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import IconButton from "@mui/material/IconButton";

const Input = ({ type, name, label, value, dir = "rtl", onChange, required }) => {
  const [direction, setDirection] = useState(dir);
  const [showPassword, setShowPassword] = useState(false);
  const isPassword = type === "password";

  // Toggle between showing and hiding the password
  const handleToggle = () => {
    setShowPassword((prev) => !prev);
  };

  // Detect input direction based on the first typed character
  useEffect(() => {
    if (type === "text" || type === "password") {
      const firstChar = value.trim().charAt(0);
      const rtlChars = /[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC]/;
      if (firstChar) {
        setDirection(rtlChars.test(firstChar) ? "rtl" : "ltr");
      }
    }
  }, [value]);

  return (
    <div className={styles.container}>
      <label className={styles.label}>{label}</label>
      <div className={styles.inputWrapper}>
        {isPassword && (
          <IconButton
            onClick={handleToggle}
            edge="end"
            size="small"
          >
            {showPassword ? (
              <VisibilityOff className={styles.icon} />
            ) : (
              <Visibility className={styles.icon} />
            )}
          </IconButton>
        )}
        <input
          className={styles.input}
          type={isPassword && showPassword ? "text" : type}
          name={name}
          value={value}
          onChange={onChange}
          required={required}
          dir={direction}
        />
      </div>
    </div>
  );
};

export default Input;
