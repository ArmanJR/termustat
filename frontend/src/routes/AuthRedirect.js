import { Navigate } from "react-router-dom";
import { getIsLoggedIn } from "../api/adminAuth";

/**
 * A route wrapper for public pages.
 * Redirects authenticated users away from some public pages if needed (e.g., login).
 */
const AuthRedirect = ({ children, redirectTo }) => {
  // If user is logged in, redirect to the specified path
  // Otherwise, render the public content
  return getIsLoggedIn() ? <Navigate to={redirectTo} /> : children;
};

export default AuthRedirect;
