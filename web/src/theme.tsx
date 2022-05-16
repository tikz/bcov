import { createTheme } from "@mui/material/styles";
import "./styles/main.scss";

const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#f05a63",
    },
    secondary: {
      main: "#8c92a4",
    },
    background: {
      default: "#252525",
      paper: "#2b2b2d",
    },
  },
});
export default theme;
