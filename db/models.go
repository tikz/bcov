package db

type Depth uint8
type Coverage float64

type BAMFile struct {
	ID             uint   `gorm:"primarykey"`
	SHA256Sum      string `gorm:"uniqueIndex"`
	Name           string
	DepthCoverages []DepthCoverage
}

type Gene struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"index"`
	EnsemblID string `gorm:"index"`
	Regions   []Region
}

type Region struct {
	ID             uint `gorm:"primarykey"`
	GeneID         uint
	Chromosome     string `gorm:"index"`
	Start          uint64 `gorm:"index"`
	End            uint64
	ExonNumber     uint
	DepthCoverages []DepthCoverage
}

type DepthCoverage struct {
	ID             uint `gorm:"primarykey"`
	RegionID       uint `gorm:"index"`
	BAMFileID      uint `gorm:"index"`
	DepthCoverages string
	Depth          uint8
	Coverage       float64
}

func automigrate() {
	db.AutoMigrate(&BAMFile{}, &Gene{}, &Region{}, &DepthCoverage{})
}
