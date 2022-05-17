import {
  Autocomplete,
  Button,
  Chip,
  CircularProgress,
  Grid,
  TextField,
  Typography,
} from "@mui/material";
import { createFilterOptions } from "@mui/material/Autocomplete";
import React from "react";
import { useNavigate } from "react-router-dom";
import { Gene, ISearchResult, Kit } from "../definitions";
import { Query, useQueryParam } from "../query";
import api from "../services";
import Results from "./Results/Results";
const stc = require("string-to-color");

export default () => {
  const navigate = useNavigate();
  const [value, setValue] = React.useState<(ISearchResult | string)[]>([]);
  const [inputValue, setInputValue] = React.useState<string>("");
  const [searchOptions, setSearchOptions] = React.useState<ISearchResult[]>([]);
  const [inProgress, setInProgress] = React.useState<boolean>(false);
  const [open, setOpen] = React.useState<boolean>(false);
  const handleClose = () => setOpen(false);

  let [query, setQuery] = useQueryParam<Query>("q");

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
      });

      setInProgress(false);
    }
  }, [inputValue]);

  const filterOptions = createFilterOptions({
    matchFrom: "any",
    stringify: (option: ISearchResult) => option.name + option.description,
  });

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
            filterOptions={filterOptions}
            options={searchOptions.map((option) => option)}
            isOptionEqualToValue={(option, value) => option.id === value.id}
            freeSolo
            filterSelectedOptions
            inputValue={inputValue}
            onInputChange={(event, newInputValue) => {
              setInputValue(newInputValue);
            }}
            onChange={(event, newValue: (string | ISearchResult)[]) => {
              setSearchOptions([]);
              setInputValue("");
              setValue(newValue);
            }}
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
                label="Filter by gene name, Ensembl accession, or capture kit name"
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
            disabled={value.length === 0}
            onClick={handleSubmit}
          >
            See coverages
          </Button>
        </Grid>
      </Grid>
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
      backgroundColor: stc(option.name) + 33,
      cursor: "pointer",
      border: 0,
      boxShadow: option instanceof Kit ? "0 0 10px #f05a63" : "inherit",
    },
  };
};
