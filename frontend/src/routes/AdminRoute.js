import { Navigate, Outlet } from "react-router-dom";
import { getIsLoggedIn, getIsAdmin, checkAdminStatus, adminLogout } from "../api/adminAuth";
import { useEffect, useState } from "react";

/**
 * A protected route component for admin-only pages.
 * Checks if the user is logged in and has admin privileges.
 */
const AdminRoute = () => {
  const [adminChecked, setAdminChecked] = useState(false); // Tracks if admin status has been verified

  useEffect(() => {
    const check = async () => {
      await checkAdminStatus(); // Verify if the logged-in user is an admin
      setAdminChecked(true);    // Mark the check as completed
    };

    // Only check admin status if user is logged in
    if (getIsLoggedIn()) {
      check();
    } else {
      setAdminChecked(true);
    }
  }, []);

  // Redirect to /admin/login if the user is not logged in
  if (!getIsLoggedIn()) {
    return <Navigate to="/admin/login" />;
  }

  // Prevent rendering until admin status check is complete
  if (!adminChecked) {
    return <div></div>;
  }

  // Show permission error if user is not an admin
  if (!getIsAdmin()) {
    return (
      <div>
        You do not have permission to access this page. <br/>
        Your session may have expired.
        <button onClick={() => adminLogout()}>Log out</button>
      </div>
    );
  }

  // Render child routes if all checks pass
  return <Outlet />;
};

export default AdminRoute;
