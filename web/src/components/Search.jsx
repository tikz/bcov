import { Autocomplete, Button, Chip, Grid, TextField } from "@mui/material";
import React from "react";

export default function Search() {
  const [inputValue, setInputValue] = React.useState("");
  const [searchOptions, setSearchOptions] = React.useState([]);

  React.useEffect(() => {
    if (inputValue.length > 2) {
      fetch("http://localhost:8080/api/search/genes/" + inputValue)
        .then((response) => response.json())
        .then((data) => {
          setSearchOptions(!("error" in data) ? data : []);
        });
    }
  }, [inputValue]);

  return (
    <Grid container alignItems="center" spacing={1}>
      <Grid item>
        <Autocomplete
          multiple
          id="search"
          options={searchOptions.map((option) => option.name)}
          freeSolo
          inputValue={inputValue}
          onInputChange={(event, newInputValue) => {
            setInputValue(newInputValue);
          }}
          onChange={(event, newValue) => {
            setSearchOptions([]);
            setInputValue("");
          }}
          renderTags={(value, getTagProps) =>
            value.map((option, index) => (
              <Chip
                variant="outlined"
                label={option}
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
        />
      </Grid>
      <Grid item>
        <Button variant="contained" color="primary" sx={{ height: 56 }}>
          See coverages
        </Button>
      </Grid>
    </Grid>
  );
}
