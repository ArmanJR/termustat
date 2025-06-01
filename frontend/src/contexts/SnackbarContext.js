import { createContext, useContext, useState } from "react";
import { Snackbar, Alert } from "@mui/material";

const SnackbarContext = createContext();

export const useSnackbar = () => useContext(SnackbarContext);

export const SnackbarProvider = ({ children }) => {
  const [message, setMessage] = useState(null);
  const [severity, setSeverity] = useState("success");

  const showSnackbar = (msg, sev = "success") => {
    setMessage(msg);
    setSeverity(sev);
  };

  const hideSnackbar = () => {
    setMessage(null);
  };

  return (
    <SnackbarContext.Provider value={{ showSnackbar, hideSnackbar }}>
      {children}
      {message && (
        <Snackbar
          open={Boolean(message)}
          onClose={hideSnackbar}
        >
          <Alert onClose={hideSnackbar} severity={severity}>
            {message}
          </Alert>
        </Snackbar>
      )}
    </SnackbarContext.Provider>
  );
};
