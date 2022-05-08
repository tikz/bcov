package db

type Depth uint8
type Coverage float64

type Kit struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Name     string    `json:"name"`
	BAMFiles []BAMFile `json:"-"`
}

type BAMFile struct {
	ID                   uint   `gorm:"primarykey"`
	SHA256Sum            string `gorm:"uniqueIndex"`
	Size                 uint64
	Name                 string
	KitID                uint
	RegionDepthCoverages []RegionDepthCoverage `gorm:"constraint:OnDelete:CASCADE;"`
}

type Gene struct {
	ID          uint `gorm:"primarykey" json:"id"`
	Accession   string `json:"accession"`
	Name        string `gorm:"index" json:"name"`
	Description string `json:"description"`
	EnsemblID   string `gorm:"index" json:"ensemblId"`
	Regions     []Region `json:"regions"`
}

type Region struct {
	ID                   uint                  `gorm:"primarykey" json:"id"`
	GeneID               uint                  `json:"-"`
	Chromosome           string                `gorm:"index" json:"chromosome"`
	Start                uint64                `gorm:"index" json:"start"`
	End                  uint64                `json:"end"`
	ExonNumber           uint                  `json:"exonNumber"`
	RegionDepthCoverages []RegionDepthCoverage `json:"-"`
}

type RegionReadCount struct {
	ID         uint        `gorm:"primarykey" json:"-"`
	RegionID   uint        `gorm:"index" json:"-"`
	BAMFileID  uint        `gorm:"index" json:"-"`
	BAMFile    BAMFile     `json:"bamFile"`
	ReadCounts []ReadCount `gorm:"constraint:OnDelete:CASCADE;" json:"readCounts"`
}

type ReadCount struct {
	ID                uint   `gorm:"primarykey" json:"-"`
	RegionReadCountID uint   `gorm:"index" json:"-"`
	Position          uint64 `json:"depth"`
	Count             uint64 `json:"coverage"`
}

type RegionDepthCoverage struct {
	ID             uint            `gorm:"primarykey" json:"-"`
	RegionID       uint            `gorm:"index" json:"-"`
	BAMFileID      uint            `gorm:"index" json:"-"`
	BAMFile        BAMFile         `json:"bamFile"`
	DepthCoverages []DepthCoverage `gorm:"constraint:OnDelete:CASCADE;" json:"depthCoverages"`
}

type DepthCoverage struct {
	ID                    uint  `gorm:"primarykey" json:"-"`
	RegionDepthCoverageID uint  `gorm:"index" json:"-"`
	Depth                 uint8 `json:"depth"`
	Coverage              uint8 `json:"coverage"`
}

func automigrate() {
	DB.AutoMigrate(&Kit{}, &BAMFile{}, &Gene{}, &Region{}, &RegionReadCount{}, &ReadCount{}, &RegionDepthCoverage{}, &DepthCoverage{})
}
