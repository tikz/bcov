// import logo from "./logo.svg";
import AlternateEmailIcon from "@mui/icons-material/AlternateEmail";
import GitHubIcon from "@mui/icons-material/GitHub";
import { Container, CssBaseline, Grid, Typography } from "@mui/material";
import "./App.css";
import logo from "./assets/b.svg";
import Search from "./components/Search";

import { createTheme, ThemeProvider } from "@mui/material/styles";
import Splash from "./components/Splash";

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

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Splash />
      <div className="App">
        <header className="App-header">
          {/* <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edist <code>src/App.tsx</code> and save to reload.
        </p> */}
          <Container>
            <Grid
              container
              alignItems="center"
              direction="column"
              spacing={0}
              id="search"
            >
              <Grid item>
                <div id="logo">
                  <Typography variant="h1">
                    <span style={{ color: "#fff" }}>B</span> cov
                  </Typography>
                  <div id="underscore"></div>
                </div>
              </Grid>

              <Grid item>
                <Typography variant="body1" id="desc">
                  A tool for the indexing of gene depth coverages across exome
                  sequencing data.
                </Typography>
              </Grid>
              <Grid item>
                <Grid container spacing={2} className="version">
                  <Grid item>
                    <Typography variant="overline">v1.2.3</Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="overline">8fb2c5c</Typography>
                  </Grid>
                </Grid>
              </Grid>
              <Grid item>
                <Search />
              </Grid>
            </Grid>
          </Container>
          <div id="footer">
            <Grid container justifyContent="center" spacing={3}>
              <Grid item>
                <a
                  href="https://github.com/tikz/bcov"
                  target="_blank"
                  rel="noreferrer"
                >
                  <GitHubIcon id="github" className="item" />
                </a>
              </Grid>
              <Grid item>
                <a
                  href="mailto:bcov@bitgenia.com"
                  target="_blank"
                  rel="noreferrer"
                >
                  <AlternateEmailIcon id="mailto" className="item" />
                </a>
              </Grid>
              <Grid item>
                <a
                  href="https://www.bitgenia.com"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img
                    src={logo}
                    id="bitgenia"
                    alt="Bitgenia"
                    className="item"
                  />
                </a>
              </Grid>
            </Grid>
          </div>
        </header>
      </div>
    </ThemeProvider>
  );
}

export default App;
