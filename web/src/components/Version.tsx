import { Grid, Typography } from "@mui/material";

export default () => (
  <Grid container spacing={2} className="version">
    <Grid item>
      <Typography variant="overline">v1.2.3</Typography>
    </Grid>
    <Grid item>
      <Typography variant="overline">8fb2c5c</Typography>
    </Grid>
  </Grid>
);
