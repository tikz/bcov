import NavigateBeforeIcon from "@mui/icons-material/NavigateBefore";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import {
  Button,
  Container,
  Dialog,
  Fade,
  Grid,
  Paper,
  Tooltip,
  Typography,
} from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { Box } from "@mui/system";
import React from "react";
import { Gene, IExon, Kit, Variant } from "../../definitions";
import api from "../../services";
import { stringToColor } from "../../theme";
import BAMSources from "./BAMSources";
import DepthCoverages from "./DepthCoverages";
import ExonMap from "./ExonMap";
import KitLegends from "./KitLegends";
import ReadCounts from "./ReadCounts";
import Variants from "./Variants";

type ResultsProps = {
  open: boolean;
  onClose: () => void;
  genes: Gene[];
  kits: Kit[];
  variants: Variant[];
};

export default ({ open, onClose, genes, kits, variants }: ResultsProps) => {
  const [gene, setGene] = React.useState<Gene>(
    variants.length > 0 ? variants[0].gene : genes[0]
  );
  const [exon, setExon] = React.useState<IExon | undefined>(undefined);

  React.useEffect(() => {
    api.getGene(gene.id).then((g: Gene) => {
      setGene(g);
      if (variants.length > 0) {
        setExon(
          g.exons.filter((e) => e.exonNumber === variants[0].exonNumber)[0]
        );
      } else {
        setExon(g.exons[0]);
      }
    });
  }, [genes, variants, gene.id]);

  return (
    <Dialog
      open={open}
      onClose={onClose}
      disableAutoFocus={true}
      maxWidth="lg"
      fullWidth
      scroll="body"
    >
      <Container id="results">
        <Paper elevation={2} sx={{ padding: "1rem" }} id="results-paper">
          <Grid
            container
            alignItems="flex-start"
            justifyContent="space-between"
            spacing={1}
          >
            <Grid item>
              <Typography variant="h2">{gene.name}</Typography>
              <Box
                sx={{
                  backgroundColor: stringToColor(gene.name) + "33",
                  height: 10,
                  width: 50,
                }}
              />
              <Typography
                variant="caption"
                sx={{
                  marginTop: "-15px",
                  paddingLeft: "60px",
                  display: "block",
                  color: "text.secondary",
                }}
              >
                {gene.description}
              </Typography>
            </Grid>
            <Grid item>
              <Grid container direction="column" alignItems="flex-end">
                <Grid item>
                  <Typography
                    variant="caption"
                    sx={{ color: "text.secondary" }}
                  >
                    <a
                      href={`https://www.genenames.org/data/gene-symbol-report/#!/hgnc_id/${gene.hgncAccession}`}
                      target="_blank"
                      rel="noreferrer"
                    >
                      {gene.hgncAccession}
                    </a>
                  </Typography>
                </Grid>
                <Grid item>
                  <Typography
                    variant="caption"
                    sx={{ color: "text.secondary" }}
                  >
                    <a
                      href={`https://www.ensembl.org/Homo_sapiens/Gene/Summary?db=core;g=${gene.geneAccession}`}
                      target="_blank"
                      rel="noreferrer"
                    >
                      {gene.geneAccession}
                    </a>
                  </Typography>
                </Grid>

                <Grid item>
                  <Typography
                    variant="caption"
                    sx={{ color: "text.secondary" }}
                  >
                    GRCh38.p3 (hg38)
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
          <br />
          <Fade in={true} timeout={500}>
            <Grid
              container
              justifyContent="space-between"
              alignItems="flex-start"
            >
              <Grid item>
                <Typography variant="h4" gutterBottom>
                  Transcript
                </Typography>
              </Grid>
              <Grid item>
                <Grid container direction="column" alignItems="flex-end">
                  <Grid item>
                    <Typography
                      variant="caption"
                      sx={{ color: "text.secondary" }}
                    >
                      <a
                        href={`https://www.ensembl.org/Homo_sapiens/Transcript/Summary?db=core;g=${gene.geneAccession};t=${gene.transcriptAccession}`}
                        target="_blank"
                        rel="noreferrer"
                      >
                        {gene.transcriptAccession}
                      </a>{" "}
                      MANE Select
                    </Typography>
                  </Grid>

                  {exon && (
                    <Grid item>
                      <Typography
                        variant="caption"
                        sx={{ color: "text.secondary" }}
                      >
                        {gene.exons.length > 0 && gene.exons[0].strand === 1
                          ? "forward"
                          : "reverse"}{" "}
                        strand
                      </Typography>
                    </Grid>
                  )}
                </Grid>
              </Grid>
            </Grid>
          </Fade>
          {exon && <ExonMap gene={gene} exon={exon} setExon={setExon} />}
          <Paper
            elevation={4}
            sx={{
              padding: "1rem",
              marginTop: 2,
              backgroundColor: stringToColor(gene.name) + 11,
            }}
          >
            {exon && (
              <Grid
                container
                justifyContent="space-between"
                alignItems="center"
              >
                <Grid item>
                  <Grid container alignItems="center" spacing={1}>
                    <Grid item>
                      <IconButton
                        aria-label="delete"
                        disabled={
                          (exon.strand === 1 && exon.exonNumber === 1) ||
                          (exon.strand === -1 &&
                            exon.exonNumber === gene.exons.length)
                        }
                        onClick={() => {
                          setExon(
                            exon.strand === 1
                              ? gene.exons[exon.exonNumber - 2]
                              : gene.exons[
                                  gene.exons.length - exon.exonNumber - 1
                                ]
                          );
                        }}
                      >
                        <NavigateBeforeIcon />
                      </IconButton>
                    </Grid>
                    <Grid item>
                      <Box sx={{ width: 150 }}>
                        <Typography
                          variant="h5"
                          align="center"
                          className="no-select"
                        >
                          Exon #{exon.exonNumber}
                        </Typography>
                      </Box>
                    </Grid>
                    <Grid item>
                      <IconButton
                        aria-label="delete"
                        disabled={
                          (exon.strand === 1 &&
                            exon.exonNumber === gene.exons.length) ||
                          (exon.strand === -1 && exon.exonNumber === 1)
                        }
                        onClick={() => {
                          setExon(
                            exon.strand === 1
                              ? gene.exons[exon.exonNumber]
                              : gene.exons[
                                  gene.exons.length - exon.exonNumber + 1
                                ]
                          );
                        }}
                      >
                        <NavigateNextIcon />
                      </IconButton>
                    </Grid>
                  </Grid>
                </Grid>
                <Grid item>
                  <Grid container direction="column" alignItems="flex-end">
                    <Grid item>
                      <Typography
                        variant="caption"
                        sx={{ color: "text.secondary" }}
                      >
                        chr{exon.chromosome}:{exon.start.toLocaleString()}-
                        {exon.end.toLocaleString()}
                      </Typography>
                    </Grid>
                    <Grid item>
                      <Typography
                        variant="caption"
                        sx={{ color: "text.secondary" }}
                      >
                        {exon.end - exon.start + 1} bp long
                      </Typography>
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
            )}
          </Paper>

          <KitLegends kits={kits} />

          <Grid container spacing={1}>
            <Grid item xs={6}>
              <Paper elevation={4} sx={{ padding: "1rem", marginTop: 2 }}>
                {exon && <ReadCounts kits={kits} exonId={exon.id} />}
              </Paper>
            </Grid>
            <Grid item xs={6}>
              <Paper elevation={4} sx={{ padding: "1rem", marginTop: 2 }}>
                {exon && <DepthCoverages kits={kits} exonId={exon.id} />}
              </Paper>
            </Grid>
          </Grid>

          <Grid container spacing={1}>
            <Grid item xs={12}>
              <Paper elevation={4} sx={{ padding: "1rem", marginTop: 2 }}>
                {exon && (
                  <Variants
                    kits={kits}
                    exonId={exon.id}
                    variantFilter={
                      variants.length > 0 ? "rs" + variants[0].id : ""
                    }
                  />
                )}
              </Paper>
            </Grid>
          </Grid>
          <Grid container justifyContent="space-between" sx={{ marginTop: 1 }}>
            <Grid item>
              <BAMSources kits={kits} />
            </Grid>
            <Grid item>
              <Tooltip
                title={`Export all kit depths for variants of all ${gene.name} exons`}
                arrow
              >
                <Button
                  variant="outlined"
                  href={`/api/variants-csv/${gene.name}`}
                >
                  Export CSV
                </Button>
              </Tooltip>
            </Grid>
          </Grid>
        </Paper>
      </Container>
    </Dialog>
  );
};
