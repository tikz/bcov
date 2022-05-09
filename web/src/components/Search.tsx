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

interface Gene {
  id: number;
  name: string;
  access: string;
  description: string;
  ensemblId: string;
}

export default () => {
  const [value, setValue] = React.useState<(Gene | string)[]>([]);
  const [inputValue, setInputValue] = React.useState<string>("");
  const [searchOptions, setSearchOptions] = React.useState<Gene[]>([]);
  const [inProgress, setInProgress] = React.useState<Boolean>(false);

  React.useEffect(() => {
    if (inputValue.length < 3) {
      setSearchOptions([]);
    } else {
      setInProgress(true);
      fetch("/api/search/genes/" + inputValue)
        .then((response) => response.json())
        .then((data) => {
          setSearchOptions(!("error" in data) ? data : []);
          setInProgress(false);
        });
    }
  }, [inputValue]);

  const filterOptions = createFilterOptions({
    matchFrom: "any",
    stringify: (option: Gene) => option.name + option.description,
  });

  return (
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
          onChange={(event, newValue: (string | Gene)[]) => {
            setSearchOptions([]);
            setInputValue("");
            setValue(newValue);
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
          getOptionLabel={(option: Gene | string) =>
            typeof option === "string" ? option : option.name
          }
          renderOption={(props, option: Gene, { inputValue }) => {
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
};
