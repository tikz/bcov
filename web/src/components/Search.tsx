import {
  Autocomplete,
  Button,
  Chip,
  CircularProgress,
  Grid,
  TextField,
  Typography,
} from "@mui/material";
import { Box } from "@mui/system";
import React, { ReactElement } from "react";
import { Gene, ISearchResult, Kit, Variant } from "../definitions";
import api from "../services";
import { stringToColor } from "../theme";
import Results from "./Results/Results";

export default () => {
  const [value, setValue] = React.useState<ISearchResult[]>([]);
  const [inputValue, setInputValue] = React.useState<string>("");
  const [searchOptions, setSearchOptions] = React.useState<ISearchResult[]>([]);
  const [inProgress, setInProgress] = React.useState<boolean>(false);
  const [helperText, setHelperText] = React.useState<ReactElement | null>(null);
  const [open, setOpen] = React.useState<boolean>(false);

  const lengthKits = value.filter((v) => v instanceof Kit).length;
  const lengthGenes = value.filter((v) => v instanceof Gene).length;
  const lengthVariants = value.filter((v) => v instanceof Variant).length;

  React.useEffect(() => {
    setSearchOptions([]);
    if (inputValue.length >= 3) {
      setInProgress(true);

      api.search(inputValue).then((r) => {
        if (inputValue === ":kits") {
          setValue(value.filter((v) => !(v instanceof Kit)).concat(r));
          setInputValue("");
        }
        setSearchOptions(r);
        setInProgress(false);
      });
    }
  }, [inputValue, value]);

  React.useEffect(() => {
    if (lengthKits === 0 && lengthGenes + lengthVariants > 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Enter at least one DNA capture kit. To see all available kits, type{" "}
          <b>:kits</b>
        </Typography>
      );
    }

    if (lengthKits > 0 && lengthGenes + lengthVariants === 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Enter a gene or dbSNP id.
        </Typography>
      );
    }

    if (lengthKits > 0 && lengthGenes + lengthVariants > 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Ready! you can also keep adding more DNA capture kits for comparison.
        </Typography>
      );
    }

    if (lengthKits === 0 && lengthGenes + lengthVariants === 0) {
      setHelperText(null);
    }
  }, [value, lengthGenes, lengthKits, lengthVariants]);

  const handleInputChange = (
    event: React.SyntheticEvent,
    newInputValue: string
  ) => setInputValue(newInputValue);

  const handleChange = (
    event: React.SyntheticEvent,
    newValue: (ISearchResult | string)[]
  ) => {
    setSearchOptions([]);
    setInputValue("");

    // Don't allow more than one gene or rs
    const lastAdded = newValue[newValue.length - 1];
    if (lastAdded instanceof Gene || lastAdded instanceof Variant) {
      setValue(
        newValue.filter(
          (v) =>
            !(v instanceof Gene || v instanceof Variant) ||
            ((v instanceof Gene || v instanceof Variant) && v === lastAdded)
        ) as ISearchResult[]
      );
    } else {
      setValue(newValue as ISearchResult[]);
    }
  };

  const handleClose = () => setOpen(false);
  const handleSubmit = () => setOpen(true);

  return (
    <>
      <Grid container alignItems="center" spacing={1}>
        <Grid item sx={{ width: 48 }}>
          {inProgress && <CircularProgress />}
        </Grid>
        <Grid item>
          <Autocomplete
            multiple
            ListboxProps={{
              style: { maxHeight: "15rem" },
            }}
            id="search"
            filterOptions={(options, state) => options}
            options={searchOptions.map((option) => option)}
            isOptionEqualToValue={(option, value) => option.id === value.id}
            filterSelectedOptions
            freeSolo
            autoHighlight={true}
            value={value}
            inputValue={inputValue}
            onInputChange={handleInputChange}
            onChange={handleChange}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip
                  label={option.name}
                  {...ChipStyle(option)}
                  {...getTagProps({ index })}
                />
              ))
            }
            renderInput={(params) => (
              <TextField
                {...params}
                label="Filter by gene name, capture kit name, HGNC, Ensembl or dbSNP accessions"
                variant="outlined"
                sx={{ width: 500 }}
              />
            )}
            getOptionLabel={(option: ISearchResult | string) =>
              typeof option === "string" ? option : option.name
            }
            renderOption={(props, option: ISearchResult, { inputValue }) => {
              return (
                <li {...props} key={option.id + option.name}>
                  <Grid
                    container
                    alignItems="center"
                    justifyContent="space-between"
                    spacing={1}
                  >
                    <Grid item>
                      <Chip label={option.name} {...ChipStyle(option)} />
                    </Grid>
                    <Grid item>
                      <Typography
                        variant="caption"
                        color="textSecondary"
                        sx={{
                          fontSize: 10,
                          color: option instanceof Kit ? "#f05a63" : "inherit",
                        }}
                      >
                        {option.description}
                      </Typography>
                    </Grid>
                  </Grid>
                </li>
              );
            }}
          />
        </Grid>
        <Grid item>
          <Button
            variant="contained"
            color="primary"
            sx={{ height: 56 }}
            disabled={!(lengthGenes + lengthVariants > 0 && lengthKits > 0)}
            onClick={handleSubmit}
          >
            See coverages
          </Button>
        </Grid>
      </Grid>
      <Box sx={{ height: 40 }}>{helperText}</Box>
      {open && (
        <Results
          open={open}
          onClose={handleClose}
          genes={value.filter((v): v is Gene => v instanceof Gene)}
          kits={value.filter((v): v is Kit => v instanceof Kit)}
          variants={value.filter((v): v is Variant => v instanceof Variant)}
        />
      )}
    </>
  );
};

const ChipStyle = (option: ISearchResult) => {
  return {
    variant: "outlined" as "outlined",
    sx: {
      backgroundColor:
        stringToColor(
          option instanceof Variant ? option.gene.name : option.name
        ) + "44",
      cursor: "pointer",
      border: 0,
      boxShadow: option instanceof Kit ? "0 0 2px #f05a63" : "inherit",
    },
  };
};
