import { CssBaseline } from "@mui/material";
import "./../styles/main.scss";

import { ThemeProvider } from "@mui/material/styles";
import LogRocket from "logrocket";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { theme } from "../theme";
import Index from "./Index";

if (process.env.NODE_ENV !== "development") {
  LogRocket.init("wuajdc/bcov");
}

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <Routes>
          <Route path="/">
            <Route index element={<Index />} />
          </Route>
          <Route path="*" element={<h1>Not found</h1>} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
