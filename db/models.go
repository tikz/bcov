package db

type Depth uint8
type Coverage float64

type Kit struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	BAMFiles []BAMFile
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
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"index"`
	EnsemblID string `gorm:"index"`
	Regions   []Region
}

type Region struct {
	ID                   uint `gorm:"primarykey"`
	GeneID               uint
	Chromosome           string `gorm:"index"`
	Start                uint64 `gorm:"index"`
	End                  uint64
	ExonNumber           uint
	RegionDepthCoverages []RegionDepthCoverage
}

type RegionDepthCoverage struct {
	ID             uint            `gorm:"primarykey"`
	RegionID       uint            `gorm:"index"`
	BAMFileID      uint            `gorm:"index"`
	DepthCoverages []DepthCoverage `gorm:"constraint:OnDelete:CASCADE;"`
}

type DepthCoverage struct {
	ID                    uint `gorm:"primarykey"`
	RegionDepthCoverageID uint `gorm:"index"`
	Depth                 uint8
	Coverage              uint8
}

func automigrate() {
	DB.AutoMigrate(&Kit{}, &BAMFile{}, &Gene{}, &Region{}, &DepthCoverage{})
}
