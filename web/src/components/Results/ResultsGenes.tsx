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
import BAMSources from "./BAMSources";
import DepthCoverages from "./DepthCoverages";
import ExonMap from "./ExonMap";
import KitLegends from "./KitLegends";
import ReadCounts from "./ReadCounts";
import Variants from "./Variants";

import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { stringToColor } from "../../theme";
import CircleIcon from "@mui/icons-material/Circle";

type ResultsProps = {
  open: boolean;
  onClose: () => void;
  genes: Gene[];
  kits: Kit[];
};

export default ({ open, onClose, genes, kits }: ResultsProps) => {
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
              <Typography variant="h2">Coverage by gene</Typography>
            </Grid>
            <Grid item>
              <Grid container direction="column" alignItems="flex-end">
                <Grid item>.</Grid>
              </Grid>
            </Grid>
          </Grid>
          <br />
          <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
              <TableHead>
                <TableRow>
                  <TableCell>Gene</TableCell>
                  <TableCell>Coordinates</TableCell>
                  <TableCell align="right">Mean coverage at 20X</TableCell>
                  <TableCell align="right">Covered variants at 20X</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                <TableRow
                  key={1}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    <a href="#">BRCA1</a>
                  </TableCell>
                  <TableCell>chr17:43,044,295-43,125,364</TableCell>
                  <TableCell align="right">
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("asd") + 88,
                        }}
                      />
                      55.65%
                    </span>
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("sss") + 88,
                        }}
                      />
                      15.65%
                    </span>
                  </TableCell>
                  <TableCell align="right">
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("asd") + 88,
                        }}
                      />
                      123 (34%)
                    </span>
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("sss") + 88,
                        }}
                      />
                      70 (14%)
                    </span>
                    <span className="total-variants">
                      <Typography variant="caption">of 323 variants</Typography>
                    </span>
                  </TableCell>
                </TableRow>
                <TableRow
                  key={1}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    <a href="#">DNMT1</a>
                  </TableCell>
                  <TableCell>chr3:44,295-125,364</TableCell>
                  <TableCell align="right">
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("asd") + 88,
                        }}
                      />
                      22.15%
                    </span>
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("sss") + 88,
                        }}
                      />
                      35.65%
                    </span>
                  </TableCell>
                  <TableCell align="right">
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("asd") + 88,
                        }}
                      />
                      27 (44%)
                    </span>
                    <span className="depth">
                      <CircleIcon
                        sx={{
                          color: stringToColor("sss") + 88,
                        }}
                      />
                      42 (74%)
                    </span>

                    <span className="total-variants">
                      <Typography variant="caption">of 55 variants</Typography>
                    </span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      </Container>
    </Dialog>
  );
};
