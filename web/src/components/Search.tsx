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
import { useNavigate } from "react-router-dom";
import { Gene, ISearchResult, Kit } from "../definitions";
import { Query, useQueryParam } from "../query";
import api from "../services";
import { stringToColor } from "../theme";
import Results from "./Results/Results";

export default () => {
  const navigate = useNavigate();
  const [value, setValue] = React.useState<ISearchResult[]>([]);
  const [inputValue, setInputValue] = React.useState<string>("");
  const [searchOptions, setSearchOptions] = React.useState<ISearchResult[]>([]);
  const [inProgress, setInProgress] = React.useState<boolean>(false);
  const [helperText, setHelperText] = React.useState<ReactElement | null>(null);
  const [open, setOpen] = React.useState<boolean>(false);
  const handleClose = () => setOpen(false);

  let [query, setQuery] = useQueryParam<Query>("q");

  const handleInputChange = (
    event: React.SyntheticEvent,
    newInputValue: string
  ) => {
    setInputValue(newInputValue);
  };

  const handleChange = (
    event: React.SyntheticEvent,
    newValue: (ISearchResult | string)[]
  ) => {
    setSearchOptions([]);
    setInputValue("");

    // Don't allow more than one gene
    const lastAdded = newValue[newValue.length - 1];
    if (lastAdded instanceof Gene) {
      setValue(
        newValue.filter(
          (v) => !(v instanceof Gene) || (v instanceof Gene && v === lastAdded)
        ) as ISearchResult[]
      );
    } else {
      setValue(newValue as ISearchResult[]);
    }
  };

  React.useEffect(() => {
    const query: Query = {
      geneIds: value
        .filter((v): v is Gene => v instanceof Gene)
        .map((v) => v.id),
      kitIds: value.filter((v): v is Kit => v instanceof Kit).map((v) => v.id),
    };

    if (query.geneIds.length + query.kitIds.length > 0) {
      setQuery(query);
    }
  }, [setQuery, value]);

  React.useEffect(() => {
    if (inputValue.length < 3) {
      setSearchOptions([]);
    } else {
      setInProgress(true);

      api.search(inputValue).then((r) => {
        setSearchOptions(r);
        setInProgress(false);
      });
    }
  }, [inputValue]);

  React.useEffect(() => {
    const lengthKits = value.filter((v) => v instanceof Kit).length;
    const lengthGenes = value.filter((v) => v instanceof Gene).length;

    if (lengthKits === 0 && lengthGenes > 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Enter at least one DNA capture kit. To see all available kits, type{" "}
          <b>:kits</b>
        </Typography>
      );
    }

    if (lengthKits > 0 && lengthGenes === 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Enter a gene or dbSNP id.
        </Typography>
      );
    }

    if (lengthKits > 0 && lengthGenes > 0) {
      setHelperText(
        <Typography variant="caption" color="secondary">
          Ready! you can also keep adding more DNA capture kits for comparison.
        </Typography>
      );
    }

    if (lengthKits === 0 && lengthGenes === 0) {
      setHelperText(null);
    }
  }, [value]);

  const handleSubmit = () => {
    setOpen(true);
  };

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
            disabled={
              !(
                value.filter((v) => v instanceof Gene).length > 0 &&
                value.filter((v) => v instanceof Kit).length > 0
              )
            }
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
        option instanceof Kit
          ? stringToColor(option.name) + "44"
          : stringToColor(option.name) + "33",
      cursor: "pointer",
      border: 0,
      boxShadow: option instanceof Kit ? "0 0 2px #f05a63" : "inherit",
    },
  };
};
