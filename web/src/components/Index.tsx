import { Container, Grid, Typography } from "@mui/material";
import Footer from "./Footer";
import Logo from "./Logo";
import Search from "./Search";
import SplashBackground from "./SplashBackground";
import Version from "./Version";

export default () => (
  <>
    <SplashBackground />
    <div className="App">
      <header className="App-header">
        <Container>
          <Grid
            container
            alignItems="center"
            direction="column"
            spacing={0}
            id="search"
          >
            <Grid item>
              <Logo />
            </Grid>

            <Grid item>
              <Typography variant="body1" id="desc">
                A tool for the indexing of gene depth coverages across human
                exome sequencing data.
              </Typography>
            </Grid>
            <Grid item>
              <Version />
            </Grid>
            <Grid item>
              <Search />
            </Grid>
          </Grid>
        </Container>
        <Footer />
      </header>
    </div>
  </>
);
