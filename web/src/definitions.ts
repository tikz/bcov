export interface ISearchResult {
  id: number;
  name: string;
  description: string;
}

export interface IKit extends ISearchResult {
  id: number;
  name: string;
  description: string;
}

export interface IGene extends ISearchResult {
  id: number;
  name: string;
  description: string;
  hgncAccession: string;
  geneAccession: string;
  transcriptAccession: string;
  exons: IExon[];
}

export interface IVariant extends ISearchResult {
  variantId: string;
  name: string;
  exonNumber: number;
  gene: IGene;
  geneName: string;
  geneId: number;
  clinSig: string;
  proteinChange: string;
  chromosome: string;
  start: number;
  end: number;
  depth: number;
}

export interface IExon {
  id: number;
  chromosome: string;
  start: number;
  end: number;
  exonNumber: number;
  strand: number;
}

export class Kit implements ISearchResult {
  public id: number;
  public name: string;
  public description: string;
  constructor(obj: IKit) {
    this.id = obj.id ?? 0;
    this.name = obj.name ?? "";
    this.description = "DNA exome capture probe kit";
  }
}

export class Gene implements IGene {
  public id: number;
  public name: string;
  public description: string;
  public hgncAccession: string;
  public geneAccession: string;
  public transcriptAccession: string;
  public exons: IExon[];

  constructor(obj: IGene) {
    this.id = obj?.id ?? 0;
    this.name = obj?.name ?? "";
    this.description = obj?.description ?? "";
    this.hgncAccession = obj?.hgncAccession ?? "";
    this.geneAccession = obj?.geneAccession ?? "";
    this.transcriptAccession = obj?.transcriptAccession ?? "";
    this.exons = obj?.exons ?? [];
  }
}

export interface IReadCount {
  position: number;
  avgCount: number;
}
export class ReadCount {
  public position: number;
  public avgCount: number;

  constructor(obj: IReadCount) {
    this.position = obj?.position ?? 0;
    this.avgCount = obj?.avgCount ?? 0;
  }
}

export interface IKitReadCounts {
  kitName: string;
  readCounts: ReadCount[];
}

export class KitReadCounts {
  public kitName: string;
  public readCounts: ReadCount[];

  constructor(obj: IKitReadCounts) {
    this.kitName = obj?.kitName ?? "";
    this.readCounts = obj?.readCounts ?? [];
  }
}

export interface IDepthCoverage {
  depth: number;
  coverage: number;
}
export class DepthCoverage {
  public depth: number;
  public coverage: number;

  constructor(obj: IDepthCoverage) {
    this.depth = obj?.depth ?? 0;
    this.coverage = obj?.coverage ?? 0;
  }
}

export interface IKitDepthCoverages {
  kitName: string;
  depthCoverages: DepthCoverage[];
}

export class KitDepthCoverages {
  public kitName: string;
  public depthCoverages: DepthCoverage[];

  constructor(obj: IKitDepthCoverages) {
    this.kitName = obj?.kitName ?? "";
    this.depthCoverages = obj?.depthCoverages ?? [];
  }
}

export interface IVariants {
  totalCount: number;
  pages: number;
  currentPage: number;
  variants: Variant[];
}

export class Variants {
  public totalCount: number;
  public pages: number;
  public currentPage: number;
  public variants: Variant[];

  constructor(obj: IVariants) {
    this.totalCount = obj?.totalCount ?? 0;
    this.pages = obj?.pages ?? 0;
    this.currentPage = obj?.currentPage ?? 0;
    this.variants = obj?.variants ?? [];
  }
}

export class Variant {
  public id: number;
  public name: string;
  public exonNumber: number;
  public gene: Gene;
  public geneName: string;
  public geneId: number;
  public description: string;
  public clinSig: string;
  public proteinChange: string;
  public chromosome: string;
  public start: number;
  public end: number;
  public depth: number;

  constructor(obj: IVariant) {
    this.id = obj?.id ?? 0;
    this.name = "rs" + this.id;
    this.exonNumber = obj?.exonNumber ?? 0;
    this.gene = obj?.gene ?? ({} as IGene);
    this.geneName = obj?.geneName ?? "";
    this.geneId = obj?.geneId ?? 0;
    this.description = `${this.geneName} Exon ${this.exonNumber} variant`;
    this.clinSig = obj?.clinSig ?? "";
    this.proteinChange = obj?.proteinChange ?? "";
    this.chromosome = obj?.chromosome ?? "";
    this.start = obj?.start ?? 0;
    this.end = obj?.end ?? 0;
    this.depth = obj?.depth ?? 0;
  }
}

export interface IBAMFile {
  name: string;
  size: number;
  SHA256Sum: string;
}
export class BAMFile {
  public name: string;
  public size: number;
  public SHA256Sum: string;

  constructor(obj: IBAMFile) {
    this.name = obj?.name ?? "";
    this.size = obj?.size ?? 0;
    this.SHA256Sum = obj?.SHA256Sum ?? "";
  }
}
