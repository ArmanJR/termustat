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
  },
});

export default UniversityFormTheme;
