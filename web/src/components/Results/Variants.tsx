import React from "react";
import { Kit, Variant, Variants } from "../../definitions";
import api from "../../services";

import CircleIcon from "@mui/icons-material/Circle";
import {
  Box,
  Checkbox,
  Divider,
  Fade,
  FormControlLabel,
  FormGroup,
  Grid,
  LinearProgress,
  Pagination,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { stringToColor } from "../../theme";

type VariantsProps = {
  kits: Kit[];
  exonId: number;
  variant: string;
};

export default ({ kits, exonId, variant }: VariantsProps) => {
  const [kitVariants, setKitVariants] = React.useState<Variant[][]>([]);
  const [variants, setVariants] = React.useState<Variant[]>([]);
  const [variantFilter, setVariantFilter] = React.useState<string>(variant);
  const [pathogenicFilter, setPathogenicFilter] =
    React.useState<Boolean>(false);

  // const [depthFilter, setDepthFilter] = React.useState<number[]>([0, 200]);

  const [loaded, setLoaded] = React.useState<Boolean>(false);

  const [totalCount, setTotalCount] = React.useState<number>(0);
  const [totalPages, setTotalPages] = React.useState<number>(1);
  const [page, setPage] = React.useState<number>(1);

  React.useEffect(() => {
    setLoaded(false);
    (async () => {
      const kvs: Variants[] = await Promise.all(
        kits.map((kit) =>
          api.getVariants(
            kit.id,
            exonId,
            page,
            variantFilter.replace("rs", ""),
            pathogenicFilter
          )
        )
      );
      setKitVariants(kvs.map((kv) => kv.variants));
      setVariants(kvs[0].variants);
      setTotalPages(kvs[0].pages);
      setTotalCount(kvs[0].totalCount);
      setLoaded(true);
    })();
  }, [kits, exonId, page, variantFilter, pathogenicFilter]);

  React.useEffect(() => {
    setPage(1);
  }, [exonId, pathogenicFilter]);

  const filterID = (event: React.ChangeEvent<HTMLInputElement>) => {
    setVariantFilter(event.target.value);
    setPage(1);
  };

  const handlePathogenicFilterChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setPathogenicFilter(event.target.checked);
  };

  // const handleDepthFilterChange = (
  //   event: Event,
  //   newValue: number | number[]
  // ) => {
  //   setDepthFilter(newValue as number[]);
  // };

  return (
    <div id="variants">
      <Grid container justifyContent="space-between" alignItems="center">
        <Grid item md={2}>
          <Typography variant="h6" gutterBottom>
            Variants {`(${totalCount})`}
          </Typography>
        </Grid>
        <Grid item md={4}>
          <Grid container justifyContent="flex-end" spacing={2}>
            <Grid item>
              <Pagination
                count={totalPages}
                page={page}
                onChange={(
                  event: React.ChangeEvent<unknown>,
                  value: number
                ) => {
                  setPage(value);
                }}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      <Grid container alignItems="center">
        {/* <Grid item xs={6}>
          <Typography gutterBottom>Depth range</Typography>
          <Slider
            getAriaLabel={() => "Depth range"}
            value={depthFilter}
            onChange={handleDepthFilterChange}
            min={0}
            max={200}
            valueLabelDisplay="auto"
            marks={[
              { value: 0, label: "0" },
              { value: 200, label: "> 200" },
            ]}
          />
        </Grid> */}
        <Grid item>
          <FormGroup>
            <FormControlLabel
              control={
                <Checkbox
                  value={pathogenicFilter}
                  onChange={handlePathogenicFilterChange}
                />
              }
              label="Pathogenic"
            />
          </FormGroup>
        </Grid>
        <Grid item>
          <TextField
            label="dbSNP ID"
            variant="outlined"
            onChange={filterID}
            value={variantFilter}
            size="small"
          />
        </Grid>
      </Grid>

      <Grid item>
        <Box sx={{ height: "5px" }}>{!loaded && <LinearProgress />}</Box>
        <Fade in={true}>
          <TableContainer component={Paper}>
            <Table size="small" aria-label="variants">
              <TableHead>
                <TableRow>
                  <TableCell>dbSNP ID</TableCell>
                  <TableCell>Coordinates</TableCell>
                  <TableCell>Protein change</TableCell>
                  <TableCell>Clinical significance</TableCell>
                  <TableCell align="right">Mean read depth</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {variants.map((v, iv) => (
                  <TableRow
                    key={v.id}
                    sx={{
                      "&:last-child td, &:last-child th": { border: 0 },
                    }}
                  >
                    <TableCell component="th" scope="row">
                      <Stack
                        direction="row"
                        divider={<Divider orientation="vertical" flexItem />}
                        spacing={1}
                      >
                        <a
                          href={`https://www.ncbi.nlm.nih.gov/snp/rs${v.id}`}
                          target="_blank"
                          rel="noreferrer"
                          key={v.id}
                        >
                          rs{v.id}
                        </a>
                      </Stack>
                    </TableCell>
                    <TableCell className="variant-coords">
                      {v.start === v.end
                        ? `chr${v.chromosome}:${v.start.toLocaleString()}`
                        : `chr${
                            v.chromosome
                          }:${v.start.toLocaleString()}-${v.end.toLocaleString()} (${
                            v.end - v.start + 1
                          } bp)`}
                    </TableCell>
                    <TableCell>{v.proteinChange}</TableCell>
                    <TableCell className={variantColorClass(v.clinSig)}>
                      {v.clinSig}
                    </TableCell>
                    <TableCell align="right">
                      <Stack
                        direction="row"
                        divider={<Divider orientation="vertical" flexItem />}
                        spacing={1}
                        justifyContent="flex-end"
                      >
                        {kits.map((k, ik) => {
                          const depth = kitVariants[ik][iv].depth;
                          return (
                            <Tooltip title={k.name} arrow>
                              <div
                                className={`variant-depth ${
                                  depth < 30 ? "variant-depth-muted" : ""
                                }`}
                              >
                                <CircleIcon
                                  sx={{
                                    color: stringToColor(k.name) + 88,
                                  }}
                                />
                                {depth}
                              </div>
                            </Tooltip>
                          );
                        })}
                      </Stack>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Fade>

        {loaded && variants.length === 0 && (
          <div id="variants">
            <Typography variant="caption" align="center">
              No variants found for this exon.
            </Typography>
          </div>
        )}
      </Grid>
    </div>
  );
};

const variantColorClass = (clinSig: string): string => {
  if (clinSig.includes("athogenic")) {
    return "variant-pathogenic";
  }

  if (clinSig.includes("enign")) {
    return "variant-benign";
  }

  return "variant-uncertain";
};
