import { Gene, Kit } from "./definitions";

const search = (pattern: string) => {
  return Promise.all([
    fetch(`/api/search/kits/${pattern}`).then((r) => r.json()),
    fetch(`/api/search/genes/${pattern}`).then((r) => r.json()),
  ]).then(([kits, genes]) => [
    ...kits.map((k: any) => new Kit(k)),
    ...genes.map((g: any) => new Gene(g)),
  ]);
};

const getKit = (id: number) => {
  return fetch(`/api/kit/${id}`)
    .then((r) => r.json())
    .then((k: any) => new Kit(k));
};

const getGene = (id: number) => {
  return fetch(`/api/gene/${id}`)
    .then((r) => r.json())
    .then((g: any) => new Gene(g));
};

const api = { search, getKit, getGene };
export default api;
