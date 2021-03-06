import {
  BAMFile,
  Gene,
  IBAMFile,
  IGene,
  IKit,
  IVariant,
  Kit,
  KitDepthCoverages,
  KitReadCounts,
  Variant,
  Variants,
} from "./definitions";

const search = async (pattern: string) => {
  const [kits, genes, variants] = await Promise.all([
    fetch(`/api/search/kits/${pattern}`).then((r) => r.json()),
    fetch(`/api/search/genes/${pattern}`).then((r) => r.json()),
    fetch(`/api/search/variant/${pattern}`).then((r) => r.json()),
  ]);
  return [
    ...kits.map((k: IKit) => new Kit(k)),
    ...genes.map((g: IGene) => new Gene(g)),
    ...variants.map((v: IVariant) => new Variant(v)),
  ];
};

const searchGenes = async (pattern: string) => {
  const r = await fetch(`/api/search/genes/${pattern}`);
  return (await r.json()).map((g: IGene) => new Gene(g));
};

const searchKits = async (pattern: string) => {
  const r = await fetch(`/api/search/kits/${pattern}`);
  return (await r.json()).map((k: IKit) => new Kit(k));
};

const getKit = async (id: number) => {
  const r = await fetch(`/api/kit/${id}`);
  return new Kit(await r.json());
};

const getGene = async (id: number) => {
  const r = await fetch(`/api/gene/${id}`);
  return new Gene(await r.json());
};

const getReadCounts = async (kitId: number, exonId: number) => {
  const r = await fetch(`/api/reads/${kitId}/${exonId}`);
  return new KitReadCounts(await r.json());
};

const getDepthCoverages = async (kitId: number, exonId: number) => {
  const r = await fetch(`/api/depth-coverages/${kitId}/${exonId}`);
  return new KitDepthCoverages(await r.json());
};

const getVariants = async (
  kitId: number,
  exonId: number,
  page: number,
  variantFilter: string,
  pathogenicFilter: Boolean
) => {
  const r = await fetch(
    `/api/variants/${kitId}/${exonId}?page=${page}&filter_id=${variantFilter}&pathogenic=${pathogenicFilter}`
  );
  return new Variants(await r.json());
};

const getBAMs = async (kitId: number) => {
  const r = await fetch(`/api/bams/${kitId}`);
  return (await r.json()).map((b: IBAMFile) => new BAMFile(b));
};

const api = {
  search,
  searchGenes,
  searchKits,
  getKit,
  getGene,
  getReadCounts,
  getDepthCoverages,
  getVariants,
  getBAMs,
};
export default api;
