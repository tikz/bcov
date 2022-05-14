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

interface SearchResult {
  id: number;
  name: string;
  description: string;
}

class Kit implements SearchResult {
  public description: string;
  constructor(public id: number, public name: string) {
    this.id = id;
    this.name = name;
    this.description = "DNA capture kit";
  }
}
class Gene implements SearchResult {
  constructor(
    public id: number,
    public name: string,
    public description: string,
    public access: string,
    public ensemblId: string
  ) {
    this.id = id;
    this.name = name;
    this.description = description;
    this.access = access;
    this.ensemblId = ensemblId;
  }
}

// function OptionChip({ option: SearchResult }) {
//   console.log(option);
//   return (
//     <Chip
//       key={option.id}
//       variant="outlined"
//       label={option.name}
//       sx={{
//         backgroundColor: stc(option.name) + 33,
//         cursor: "pointer",
//         borderColor: option instanceof Kit ? "#f05a63" : "transparent",
//       }}
//     />
//   );
// }

// type ChipProps = {
//   option: SearchResult;
//   getTagPropsFunc?: AutocompleteRenderGetTagProps;
//   index?: number;
// };

// const OptionChip = ({ option, getTagPropsFunc, index }: ChipProps) => {
//   const additional = {...getTagPropsFunc({ index })}
//   return <Chip
//   variant="outlined"
//   label={option.name}
//   sx={{
//     backgroundColor: stc(option.name) + 33,
//     cursor: "pointer",
//     borderColor: option instanceof Kit ? "#f05a63" : "transparent",
//   }}
//   {additional}
// />
// }

export interface ItemProps {
  variant: "outlined";
}

const ChipStyle = (option: SearchResult) => {
  return {
    variant: "outlined",
    sx: {
      backgroundColor: stc(option.name) + 33,
      cursor: "pointer",
      borderColor: option instanceof Kit ? "#f05a63" : "transparent",
    },
  } as ItemProps;
};

export default () => {
  const [value, setValue] = React.useState<(SearchResult | string)[]>([]);
  const [inputValue, setInputValue] = React.useState<string>("");
  const [searchOptions, setSearchOptions] = React.useState<SearchResult[]>([]);
  const [inProgress, setInProgress] = React.useState<Boolean>(false);

  React.useEffect(() => {
    if (inputValue.length < 3) {
      setSearchOptions([]);
    } else {
      setInProgress(true);

      (async () => {
        let [kits, genes] = await Promise.all([
          fetch("/api/search/kits/" + inputValue).then((response) =>
            response.json()
          ),
          fetch("/api/search/genes/" + inputValue).then((response) =>
            response.json()
          ),
        ]);
        setSearchOptions([
          ...kits.map((k: any) => new Kit(k.id, k.name)),
          ...genes.map(
            (k: any) =>
              new Gene(k.id, k.name, k.description, k.access, k.ensemblId)
          ),
        ]);
        setInProgress(false);
      })();
    }
  }, [inputValue]);

  const filterOptions = createFilterOptions({
    matchFrom: "any",
    stringify: (option: SearchResult) => option.name + option.description,
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
          onChange={(event, newValue: (string | SearchResult)[]) => {
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
          getOptionLabel={(option: SearchResult | string) =>
            typeof option === "string" ? option : option.name
          }
          renderOption={(props, option: SearchResult, { inputValue }) => {
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
        >
          See coverages
        </Button>
      </Grid>
    </Grid>
  );
};
