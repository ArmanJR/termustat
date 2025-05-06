import { createTheme } from "@mui/material/styles";

const UniversitiesTheme = createTheme({
  components: {
    MuiIconButton: {
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
  },
});

export default UniversitiesTheme;
