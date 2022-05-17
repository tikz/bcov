import CircleIcon from "@mui/icons-material/Circle";
import { Container, Grid, Modal, Paper, Typography } from "@mui/material";
import React from "react";
import { Gene, IRegion, Kit } from "../../definitions";
import api from "../../services";
import ExonMap from "./ExonMap";
import ReadCounts from "./ReadCounts";

const stc = require("string-to-color");
const style = {
  position: "absolute" as "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  border: 0,
  zIndex: 10,
};

type ResultsProps = {
  open: boolean;
  onClose: () => void;
  genes: Gene[];
  kits: Kit[];
};

export default ({ open, onClose, genes, kits }: ResultsProps) => {
  const [currentGeneIndex, setCurrentGeneIndex] = React.useState<number>(0);
  const [gene, setGene] = React.useState<Gene>(genes[currentGeneIndex]);
  const [region, setRegion] = React.useState<IRegion>({} as IRegion);

  React.useEffect(() => {
    api.getGene(genes[currentGeneIndex].id).then((g: Gene) => {
      setGene(g);
    });
  }, [currentGeneIndex, genes]);
  console.log(region);
  return (
    <Modal open={open} onClose={onClose} disableAutoFocus={true}>
      <Container sx={style}>
        <Paper elevation={2} sx={{ padding: "1rem" }}>
          <Grid container alignItems="center" spacing={1}>
            <Grid item>
              <CircleIcon sx={{ color: stc(gene.name) }} />
            </Grid>
            <Grid item>
              <Typography variant="h2">{gene.name}</Typography>
              <Typography
                variant="caption"
                sx={{
                  paddingLeft: 4,
                  marginTop: -0.5,
                  display: "block",
                  color: "text.secondary",
                }}
              >
                {gene.description}
              </Typography>
            </Grid>
          </Grid>
          <br />
          <Typography variant="h4" gutterBottom>
            Exons
          </Typography>
          {gene.regions.length > 0 && (
            <ExonMap gene={gene} setRegion={setRegion} />
          )}
          <Typography variant="h4" gutterBottom>
            Read counts
          </Typography>
          {region.id && <ReadCounts kits={kits} regionId={region.id} />}

          {/* <>
            <h1>Results</h1>

            <p>The current form values are:</p>
            <pre>{JSON.stringify(query || {}, null, 2)}</pre>
          </> */}
        </Paper>
      </Container>
    </Modal>
  );
};
