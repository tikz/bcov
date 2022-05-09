import { CssBaseline } from "@mui/material";
import "./styles/main.scss";

import { ThemeProvider } from "@mui/material/styles";
import Index from "./components/Index";
import theme from "./theme";

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Index />
    </ThemeProvider>
  );
}

export default App;
