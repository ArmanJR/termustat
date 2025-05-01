import { useEffect, useState, useRef } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

/**
 * A route wrapper for public pages.
 * Redirects authenticated users away from some public pages if needed (e.g., login).
 */
const AuthRedirect = ({ children, redirectTo }) => {
  const isFirstRender = useRef(true);
  const { tryRefreshToken, isLoggedIn } = useAuth();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return; // Skip the effect on first render (loading state)
    }
    const refreshToken = async () => {
      await tryRefreshToken();
      setLoading(false);
    };
    refreshToken();
  }, []);

  if (loading) {
    return <div></div>;
  }
  return isLoggedIn ? <Navigate to={redirectTo} /> : children;
};

export default AuthRedirect;
