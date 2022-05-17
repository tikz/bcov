package db

type Depth uint8
type Coverage float64

type Kit struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Name     string    `json:"name"`
	BAMFiles []BAMFile `json:"-"`
}

type BAMFile struct {
	ID                 uint   `gorm:"primarykey"`
	SHA256Sum          string `gorm:"uniqueIndex"`
	Size               uint64
	Name               string
	KitID              uint
	ExonDepthCoverages []ExonDepthCoverage `gorm:"constraint:OnDelete:CASCADE;"`
}

type Gene struct {
	ID                  uint   `gorm:"primarykey" json:"id"`
	HGNCAccession       string `json:"hgncAccession"`
	GeneAccession       string `gorm:"index"`
	Name                string `gorm:"index"`
	Description         string `json:"description"`
	TranscriptAccession string `json:"transcriptAccession"`
	Exons               []Exon `json:"exons"`
}

type Exon struct {
	ID                 uint `gorm:"primarykey" json:"id"`
	GeneID             uint `json:"-"`
	Strand             int
	Chromosome         string              `gorm:"index" json:"chromosome"`
	Start              uint64              `gorm:"index" json:"start"`
	End                uint64              `json:"end"`
	ExonNumber         uint                `json:"exonNumber"`
	ExonDepthCoverages []ExonDepthCoverage `json:"-"`
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

func automigrate() {
	DB.AutoMigrate(&Kit{}, &BAMFile{}, &Gene{}, &Exon{}, &ExonReadCount{}, &ReadCount{}, &ExonDepthCoverage{}, &DepthCoverage{})
}
