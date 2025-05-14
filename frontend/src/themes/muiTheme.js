import { createTheme } from "@mui/material/styles";

const muiTheme = createTheme({
  typography: {
    fontFamily: "Vazirmatn, sans-serif",
  },
  components: {
    MuiDrawer: {
      defaultProps: {
        anchor: "right",
      },
      styleOverrides: {
        paper: {
          position: "relative",
          width: 240,
          border: "none",
          boxSizing: "border-box",
        },
      },
    },
    MuiListItemText: {
      styleOverrides: {
        root: {
          textAlign: "right",
        },
      },
    },
    MuiIconButton: {
      defaultProps: {
        edge: "start",
      },
      variants: [
        {
          props: { variant: "addButton" },
          style: {
            backgroundColor: "#309a9a",
            color: "#ffffff",
            "&:hover": {
              backgroundColor: "#42baba",
            },
            "& svg": {
              fontSize: "1.2em",
            },
          },
        },
      ],
    },
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

export default muiTheme;
