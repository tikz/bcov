import AlternateEmailIcon from "@mui/icons-material/AlternateEmail";
import GitHubIcon from "@mui/icons-material/GitHub";
import { Grid } from "@mui/material";
import logo from "../assets/b.svg";

export default () => (
  <div id="footer">
    <Grid container justifyContent="center" spacing={3}>
      <Grid item>
        <a href="https://github.com/tikz/bcov" target="_blank" rel="noreferrer">
          <GitHubIcon id="github" className="item" />
        </a>
      </Grid>
      <Grid item>
        <a href="mailto:bcov@bitgenia.com" target="_blank" rel="noreferrer">
          <AlternateEmailIcon id="mailto" className="item" />
        </a>
      </Grid>
      <Grid item>
        <a href="https://www.bitgenia.com" target="_blank" rel="noreferrer">
          <img src={logo} id="bitgenia" alt="Bitgenia" className="item" />
        </a>
      </Grid>
    </Grid>
  </div>
);
