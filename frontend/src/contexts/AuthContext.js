import { createContext, useContext, useState } from "react";
import axios from "axios";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [accessToken, setAccessToken] = useState(null);
  const [isLoggedIn, setIsLoggedIn] = useState(null);
  const [isLoggingOut, setIsLoggingOut] = useState(null);
  const [isAdmin, setIsAdmin] = useState(null);

  const login = async (credentials, setError, navigate) => {
    try {
      const response = await axios.post(
        "http://localhost:8080/api/v1/auth/login",
        credentials,
        {
          headers: { "Content-Type": "application/json" },
          withCredentials: true,
        }
      );
      const token = response.data.access_token;
      setAccessToken(token);
      setIsLoggedIn(true);
      const payload = JSON.parse(atob(token.split(".")[1]));
      setIsAdmin(payload.scp[0] == "admin-dashboard");
      setError(null);
      navigate("/admin/dashboard");
    } catch (error) {
      if (error.response?.status === 401)
        setError("نام کاربری یا رمز عبور اشتباه است");
      if (error.response?.status === 403)
        setError("ایمیل شما هنوز تأیید نشده است. لطفاً برای ادامه، ایمیل خود را تأیید کنید.");
      setIsLoggedIn(false);
    }
  };

  const logout = async () => {
    try {
      await axios.post(
        "http://localhost:8080/api/v1/auth/logout",
        {},
        {
          withCredentials: true,
        }
      );
      setIsLoggingOut(true);
      setIsLoggedIn(null);
      setAccessToken(null);
      setIsAdmin(null);
      window.location.href = "/";
    } catch (error) {
      console.log("Logout failed:", error);
    }
  };

  const tryRefreshToken = async () => {
    try {
      const response = await axios.post(
        "http://localhost:8080/api/v1/auth/refresh",
        {},
        { withCredentials: true }
      );
      const token = response.data.access_token;
      setAccessToken(token);
      setIsLoggedIn(true);
      const payload = JSON.parse(atob(token.split(".")[1]));
      setIsAdmin(payload.scp[0] == "admin-dashboard");
      return token;
    } catch (error) {
      setAccessToken(null);
      setIsLoggedIn(false);
      return null;
    }
  };

  return (
    <AuthContext.Provider
      value={{
        accessToken,
        isLoggedIn,
        isLoggingOut,
        isAdmin,
        login,
        logout,
        tryRefreshToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
