import { useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";

const AuthListener = () => {
  const { setIsLoggedIn, setIsLoggingOut, setAccessToken, setIsAdmin } = useAuth();

  useEffect(() => {
    const channel = new BroadcastChannel('auth-channel');
    channel.onmessage = (event) => {
      if (event.data === 'logout') {
        setIsLoggingOut(true);
        setIsLoggedIn(null);
        setAccessToken(null);
        setIsAdmin(null);
        window.location.href = "/";
      }
    };
    return () => {
      channel.close();
    };
  }, []);

  return null;
};

export default AuthListener;
