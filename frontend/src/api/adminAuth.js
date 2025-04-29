import axios from "axios";

/**
  * Logs in the user with given credentials.
  */
export async function adminLogin(credentials, navigate, setError) {
  try {
    const response = await axios.post(
      "http://localhost:8080/api/v1/auth/login",
      credentials,
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    // On successful login, store the token in local storage
    const token = response.data.token;
    localStorage.setItem("token", token);

    // Redirect to admin dashboard
    navigate("/admin/dashboard");
    setError(null);
  } catch (error) {
    setError("نام کاربری یا رمز عبور اشتباه است");
  }
};

/**
 * Checks whether the user is currently logged in.
 */
export function getIsLoggedIn() {
  const token = localStorage.getItem("token");
  return token ? true : false;
}

/**
  * Retrieves the JWT token.
  * Used to authenticate API requests.
  */
export function getToken() {
  return localStorage.getItem("token");
}

/**
  * Checks if the currently logged-in user is an admin.
  * Sends a request to a protected admin endpoint to verify access.
  */
let adminStatus = null;
export async function checkAdminStatus() {
  try {
    const response = await axios.get("http://localhost:8080/api/v1/admin/semesters", {
      headers: {
        Authorization: `Bearer ${getToken()}`
      }
    });
    adminStatus = response.status === 200;
    return adminStatus;
  } catch (error) {
    console.error("Error checking admin:", error.response?.data || error.message);
    adminStatus = false;
    return false;
  }
}

/**
  * Returns the current admin status.
  * Should be used after calling checkAdminStatus().
  */
export function getIsAdmin() {
  return adminStatus === true;
}

/**
  * Logs out the current admin user.
  */
export function adminLogout() {
  localStorage.removeItem("token");
  window.location.href = "/";
}
