import styles from "./Input.module.css";
import { useState } from "react";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import IconButton from "@mui/material/IconButton";

const Input = ({ type, name, label, value, onChange, required }) => {
  const [showPassword, setShowPassword] = useState(false);
  const isPassword = type === "password";

  // Toggle between showing and hiding the password
  const handleToggle = () => {
    setShowPassword((prev) => !prev);
  };

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
        />
      </div>
    </div>
  );
};

export default Input;
