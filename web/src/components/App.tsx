import { CssBaseline } from "@mui/material";
import "./../styles/main.scss";

import { ThemeProvider } from "@mui/material/styles";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { theme } from "../theme";
import Index from "./Index";

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
