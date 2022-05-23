import { Grid, Typography } from "@mui/material";
import packageJson from "../../package.json";

export default () => (
  <Grid container spacing={2} className="version">
    <Grid item>
      <Typography variant="overline">Core {packageJson.version}</Typography>
    </Grid>
    <Grid item>
      <Typography variant="overline">{process.env.REACT_APP_COMMIT}</Typography>
    </Grid>
  </Grid>
);
