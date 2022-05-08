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
const stc = require("string-to-color");

export default function Search() {
  const [value, setValue] = React.useState(null);
  const [inputValue, setInputValue] = React.useState("");
  const [searchOptions, setSearchOptions] = React.useState([]);
  const [inProgress, setInProgress] = React.useState(false);

  React.useEffect(() => {
    if (inputValue.length < 3) {
      setSearchOptions([]);
    } else {
      setInProgress(true);
      fetch("http://localhost:8080/api/search/genes/" + inputValue)
        .then((response) => response.json())
        .then((data) => {
          setSearchOptions(!("error" in data) ? data : []);
          setInProgress(false);
        });
    }
  }, [inputValue]);

  const filterOptions = createFilterOptions({
    matchFrom: "any",
    stringify: (option) => option.name + option.description,
  });

  return (
    <Grid container alignItems="center" spacing={1}>
      <Grid item sx={{ width: 48 }}>
        {inProgress && <CircularProgress />}
      </Grid>
      <Grid item>
        <Autocomplete
          multiple
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
          onChange={(event, newValue) => {
            setSearchOptions([]);
            setInputValue("");
            setValue(newValue);
            console.log("changed", newValue);
          }}
          renderTags={(value, getTagProps) =>
            value.map((option, index) => (
              <Chip
                variant="outlined"
                label={option.name}
                sx={{ backgroundColor: stc(option.name) + 33 }}
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
          getOptionLabel={(option) => option.name}
          renderOption={(props, option, { inputValue }) => {
            return (
              <li {...props}>
                <Grid
                  container
                  alignItems="center"
                  justifyContent="space-between"
                  spacing={1}
                >
                  <Grid item>
                    <Chip
                      variant="outlined"
                      label={option.name}
                      sx={{
                        backgroundColor: stc(option.name) + 33,
                        cursor: "pointer",
                      }}
                    />
                  </Grid>
                  <Grid item>
                    <Typography
                      variant="caption"
                      color="textSecondary"
                      sx={{ fontSize: 10 }}
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
        >
          See coverages
        </Button>
      </Grid>
    </Grid>
  );
}
