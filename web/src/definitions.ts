export interface ISearchResult {
  id: number;
  name: string;
  description: string;
}

interface IKit extends ISearchResult {
  id: number;
  name: string;
  description: string;
}

interface IGene extends ISearchResult {
  id: number;
  name: string;
  description: string;
  accession: string;
  ensemblId: string;
  regions: IRegion[];
}

interface IRegion {
  id: number;
  chromosome: string;
  start: number;
  end: number;
  exonNumber: number;
}

export class Kit implements ISearchResult {
  public id: number;
  public name: string;
  public description: string;
  constructor(obj: IKit) {
    this.id = obj.id ?? 0;
    this.name = obj.name ?? "";
    this.description = "DNA capture kit";
  }
}

export class Gene implements IGene {
  public id: number;
  public name: string;
  public description: string;
  public accession: string;
  public ensemblId: string;
  public regions: IRegion[];

  constructor(obj: IGene) {
    this.id = obj?.id ?? 0;
    this.name = obj?.name ?? "";
    this.description = obj?.description ?? "";
    this.accession = obj?.accession ?? "";
    this.ensemblId = obj?.ensemblId ?? "";
    this.regions = obj?.regions ?? [];
  }
}
