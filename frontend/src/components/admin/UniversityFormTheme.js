import { createTheme } from "@mui/material/styles";

const UniversityFormTheme = createTheme({
  components: {
    MuiDialog: {
      styleOverrides: {
        paper: {
          width: "100%",
          maxWidth: "460px",
        },
      },
    },
    MuiSnackbar: {
      defaultProps: {
        anchorOrigin: { vertical: "bottom", horizontal: "center" },
        autoHideDuration: 3000,
      },
    },
    MuiAlert: {
      defaultProps: {
        variant: "filled",
      },
      styleOverrides: {
        filled: {
          width: "100%",
        },
      },
    },
  },
});

export default UniversityFormTheme;
