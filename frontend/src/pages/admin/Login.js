import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import styles from "./Login.module.css";
import logo from "../../images/logo.png";
import Input from "../../components/form/Input";
import Button from "../../components/form/Button";

const AdminLoginPage = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [error, setError] = useState(null);

  // Update state whenever user types into the input fields
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name.slice(5)]: e.target.value,
    });
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    login(formData, setError, navigate);
  };

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.loginBox}>
        <Link to="/">
          <img
            className={styles.logo}
            src={logo}
            alt="لوگو - بازگشت به صفحه اصلی"
            title="صفحه اصلی"
          />
        </Link>
        <form className={styles.form} onSubmit={handleSubmit}>
          <Input
            label="ایمیل"
            type="email"
            name="user_email"
            value={formData.email}
            dir="ltr"
            onChange={handleChange}
            required
          />
          <Input
            label="رمز عبور"
            type="password"
            name="user_password"
            value={formData.password}
            dir="ltr"
            onChange={handleChange}
            required
          />
          <div
            className={styles.errorMessage}
            style={error && {visibility: "visible"}}
          >
            {error}
          </div>
          <Button type="submit" value="ورود" />
        </form>
      </div>
    </div>
  );
};

export default AdminLoginPage;
