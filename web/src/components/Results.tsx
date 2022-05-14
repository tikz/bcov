import * as React from "react";
import type { NavigateOptions } from "react-router-dom";
import { useSearchParams } from "react-router-dom";

const JSURL = require("jsurl");

interface Query {
  genes: string[];
  kits: string[];
}

function useQueryParam<T>(
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

export default () => {
  let [query, setQuery] = useQueryParam<Query>("query");

  React.useEffect(() => {
    let query: Query = {
      genes: ["abc", "test"],
      kits: ["def", "ghi"],
    };
    console.log(query);

    setQuery(query, {
      replace: true,
    });
  }, [setQuery]);

  console.log(query);

  return (
    <>
      <h1>Results</h1>

      <p>The current form values are:</p>

      <pre>{JSON.stringify(query || {}, null, 2)}</pre>
    </>
  );
};
