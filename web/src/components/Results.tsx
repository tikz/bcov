import { Query, useQueryParam } from "../query";

// class GeneResults extends Gene {
//   public regions: string[];
//   constructor(
//     id: number
//   ) {
//     super(id, name, description, access, ensemblId);
//     this.regions = regions;
//   }

//   fetch() {
//     (async () => {
//       let gene = await fetch("/api/gene/" + this.id).then((response) =>
//         response.json()
//       );

//       console.log(gene)
//     })();
//   }
// }
// class Region {
//   constructor(
//     public id: number,
//     public chromosome: string,
//     public start: number,
//     public end: number,
//     public exonNumber: number
//   ) {
//     this.id = id;
//     this.chromosome = chromosome;
//     this.start = start;
//     this.end = end;
//     this.exonNumber = exonNumber;
//   }
// }
export default () => {
  let [query, setQuery] = useQueryParam<Query>("q");

  return (
    <>
      <h1>Results</h1>

      <p>The current form values are:</p>
      <pre>{JSON.stringify(query || {}, null, 2)}</pre>
    </>
  );
};
