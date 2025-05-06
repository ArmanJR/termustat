import { createTheme } from "@mui/material/styles";

const AdminLayoutTheme = createTheme({
  typography: {
    fontFamily: "Vazirmatn, sans-serif",
    h1: {
      fontSize: "2rem",
    },
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
        color: "inherit",
        edge: "start",
      },
      styleOverrides: {
        root: {
          marginRight: "16px",
        },
      },
    },
  },
});

export default AdminLayoutTheme;
