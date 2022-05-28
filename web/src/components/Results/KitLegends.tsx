import { Circle } from "@mui/icons-material";
import { Grid } from "@mui/material";
import { Kit } from "../../definitions";
import { stringToColor } from "../../theme";

type KitLegendsProps = {
  kits: Kit[];
};

export default ({ kits }: KitLegendsProps) => {
  return (
    <>
      <Grid
        container
        alignItems="center"
        justifyContent="center"
        id="kit-legends"
        spacing={2}
      >
        {kits.map((k) => (
          <Grid item>
            <Circle sx={{ color: stringToColor(k.name) + "88" }} />
            {k.name}
          </Grid>
        ))}
      </Grid>
    </>
  );
};
