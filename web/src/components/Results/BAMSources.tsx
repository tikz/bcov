import React from "react";
import { BAMFile, Kit } from "../../definitions";

import { Button, Dialog, TableBody } from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import api from "../../services";

import CircleIcon from "@mui/icons-material/Circle";
import { stringToColor } from "../../theme";

type BAMSourcesProps = {
  kits: Kit[];
};

export default ({ kits }: BAMSourcesProps) => {
  const [open, setOpen] = React.useState<boolean>(false);
  const [sources, setSources] = React.useState<BAMFile[][]>([]);

  React.useEffect(() => {
    (async () => {
      const kitBAMs: BAMFile[][] = await Promise.all(
        kits.map((kit) => api.getBAMs(kit.id))
      );
      setSources(kitBAMs.map((kb) => kb));
    })();
  }, [kits]);

  console.log(sources);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <Button color="secondary" onClick={handleOpen}>
        Sources
      </Button>
      <Dialog
        open={open}
        onClose={handleClose}
        maxWidth="md"
        fullWidth
        id="bam-sources"
      >
        <TableContainer component={Paper}>
          <Table size="small" aria-label="bams">
            <TableHead>
              <TableRow>
                <TableCell>File</TableCell>
                <TableCell>Size</TableCell>
                <TableCell>SHA256 sum</TableCell>
                <TableCell align="right">DNA capture kit</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sources.map((sources, i) =>
                sources.map((bam) => (
                  <TableRow>
                    <TableCell>{bam.name}</TableCell>
                    <TableCell className="text-muted">
                      {(bam.size / 1024 / 1024 / 1024).toFixed(3)} GB
                    </TableCell>
                    <TableCell className="text-muted text-monospace">
                      {bam.SHA256Sum}
                    </TableCell>
                    <TableCell align="right">
                      <CircleIcon
                        sx={{
                          color: stringToColor(kits[i].name) + 88,
                        }}
                      />
                      {kits[i].name}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Dialog>
    </>
  );
};
