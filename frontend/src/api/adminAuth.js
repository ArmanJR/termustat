import axios from "axios";

/**
  * Logs in the user with given credentials.
  */
export async function adminLogin(credentials, navigate) {
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
  } catch (error) {
    console.error("Login failed:", error.response?.data || error.message);
    alert("Login failed:", error.response?.data || error.message);
  }
};

/**
 * Checks whether the user is currently logged in.
 */
export function getIsLoggedIn() {
  const token = localStorage.getItem("token");
  return token ? true : false;
}
