import { Gene, Kit, KitReadCounts } from "./definitions";

const search = async (pattern: string) => {
  const [kits, genes] = await Promise.all([
    fetch(`/api/search/kits/${pattern}`).then((r) => r.json()),
    fetch(`/api/search/genes/${pattern}`).then((r_1) => r_1.json()),
  ]);
  return [
    ...kits.map((k: any) => new Kit(k)),
    ...genes.map((g: any) => new Gene(g)),
  ];
};

const getKit = async (id: number) => {
  const r = await fetch(`/api/kit/${id}`);
  const k = await r.json();
  return new Kit(k);
};

const getGene = async (id: number) => {
  const r = await fetch(`/api/gene/${id}`);
  const g = await r.json();
  return new Gene(g);
};

const getReadCounts = async (kitId: number, regionId: number) => {
  const r = await fetch(`/api/reads/${kitId}/${regionId}`);
  const r_1 = await r.json();
  return new KitReadCounts(r_1);
};

const api = { search, getKit, getGene, getReadCounts };
export default api;
