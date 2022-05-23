package db

type Depth uint8
type Coverage float64

type Kit struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Name     string    `json:"name"`
	BAMFiles []BAMFile `json:"-"`
}

type BAMFile struct {
	ID                 uint                `gorm:"primarykey"`
	SHA256Sum          string              `gorm:"uniqueIndex"`
	Size               uint64              `json:"size"`
	Name               string              `json:"name"`
	KitID              uint                `json:"-"`
	ExonDepthCoverages []ExonDepthCoverage `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

type Gene struct {
	ID                  uint   `gorm:"primarykey" json:"id"`
	HGNCAccession       string `gorm:"unique" json:"hgncAccession"`
	GeneAccession       string `gorm:"uniqueIndex" json:"geneAccession"`
	Name                string `gorm:"uniqueIndex" json:"name"`
	Description         string `json:"description"`
	TranscriptAccession string `json:"transcriptAccession"`
	Exons               []Exon `json:"exons"`
}

type Exon struct {
	ID                 uint                `gorm:"primarykey" json:"id"`
	GeneID             uint                `json:"-"`
	Strand             int                 `json:"strand"`
	Chromosome         string              `gorm:"index" json:"chromosome"`
	Start              uint64              `gorm:"index" json:"start"`
	End                uint64              `json:"end"`
	ExonNumber         uint                `json:"exonNumber"`
	ExonDepthCoverages []ExonDepthCoverage `json:"-"`
	Variants           []Variant           `json:"-"`
}

type ExonReadCount struct {
	ID         uint        `gorm:"primarykey" json:"-"`
	ExonID     uint        `gorm:"index" json:"-"`
	BAMFileID  uint        `gorm:"index" json:"-"`
	BAMFile    BAMFile     `json:"bamFile"`
	ReadCounts []ReadCount `gorm:"constraint:OnDelete:CASCADE;" json:"readCounts"`
}

type ReadCount struct {
	ID              uint   `gorm:"primarykey" json:"-"`
	ExonReadCountID uint   `gorm:"index" json:"-"`
	Position        uint64 `json:"depth"`
	Count           uint64 `json:"coverage"`
}

type ExonDepthCoverage struct {
	ID             uint            `gorm:"primarykey" json:"-"`
	ExonID         uint            `gorm:"index" json:"-"`
	BAMFileID      uint            `gorm:"index" json:"-"`
	BAMFile        BAMFile         `json:"bamFile"`
	DepthCoverages []DepthCoverage `gorm:"constraint:OnDelete:CASCADE;" json:"depthCoverages"`
}

type DepthCoverage struct {
	ID                  uint  `gorm:"primarykey" json:"-"`
	ExonDepthCoverageID uint  `gorm:"index" json:"-"`
	Depth               uint8 `json:"depth"`
	Coverage            uint8 `json:"coverage"`
}

type Variant struct {
	ID            uint   `gorm:"primarykey" json:"-"`
	VariantID     string `gorm:"uniqueIndex" json:"variantId"`
	Name          string `json:"-"`
	ClinSig       string `json:"clinSig"`
	ClinSigSimple int    `json:"-"`
	ProteinChange string `json:"-"`
	ReviewStatus  string `json:"-"`
	Phenotypes    string `json:"-"`
	Chromosome    string `gorm:"index" json:"chromosome"`
	Start         uint64 `gorm:"index" json:"start"`
	End           uint64 `gorm:"index" json:"end"`
	ExonID        uint   `gorm:"index" json:"-"`
}

func automigrate() {
	DB.AutoMigrate(&Kit{}, &BAMFile{}, &Gene{}, &Exon{}, &ExonReadCount{}, &ReadCount{}, &ExonDepthCoverage{}, &DepthCoverage{}, &Variant{})
}
