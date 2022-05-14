import * as React from "react";
import type { NavigateOptions } from "react-router-dom";
import { useSearchParams } from "react-router-dom";
const JSURL = require("jsurl");

export interface Query {
  geneIds: number[];
  kitIds: number[];
}

export function generateQueryURL(query: Query) {
  let searchParams = new URLSearchParams();
  searchParams.set("q", JSURL.stringify(query));
  return searchParams.toString();
}

export function useQueryParam<T>(
  key: string
): [T | undefined, (newQuery: T, options?: NavigateOptions) => void] {
  let [searchParams, setSearchParams] = useSearchParams();
  let paramValue = searchParams.get(key);

  let value = React.useMemo(() => JSURL.parse(paramValue), [paramValue]);

  let setValue = React.useCallback(
    (newValue: T, options?: NavigateOptions) => {
      let newSearchParams = new URLSearchParams(searchParams);
      newSearchParams.set(key, JSURL.stringify(newValue));
      setSearchParams(newSearchParams, options);
    },
    [key, searchParams, setSearchParams]
  );

  return [value, setValue];
}
