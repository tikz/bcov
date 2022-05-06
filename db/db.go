package db

import (
	"bcov/bed"
	"bcov/ensembl"
	"bcov/utils"
	"fmt"
	"log"
	"os"
	"sort"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	fmt.Println("Database engine:\t", os.Getenv("BCOV_DB_ENGINE"))
	if os.Getenv("BCOV_DB_ENGINE") == "postgres" {
		ConnectPostgres()
	} else {
		ConnectSQLite()
	}

	automigrate()
}

func ConnectPostgres() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("BCOV_DB_HOST"),
		os.Getenv("BCOV_DB_USER"),
		os.Getenv("BCOV_DB_PASSWORD"),
		os.Getenv("BCOV_DB_NAME"),
		os.Getenv("BCOV_DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to connect database")
	}
}

func ConnectSQLite() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to connect database")
	}

	if res := DB.Exec("PRAGMA synchronous = OFF; PRAGMA foreign_keys = ON;", nil); res.Error != nil {
		panic(res.Error)
	}
}

func StoreFile(file string, hash string, kitID uint, size int64) (bamFileID uint, created bool) {
	var bamFileDB BAMFile
	result := DB.First(&bamFileDB, "sha256_sum = ?", hash)
	if result.RowsAffected == 0 {
		bamFileDB = BAMFile{Name: file, SHA256Sum: hash, KitID: kitID, Size: uint64(size)}
		DB.Create(&bamFileDB)
		created = true
	}

	return bamFileDB.ID, created
}

func StoreRegion(region *bed.Region) (regionID uint) {
	var regionDB Region
	result := DB.Where("chromosome = ? AND start = ? AND end = ?", region.Chromosome, region.Start, region.End).First(&regionDB)
	if result.RowsAffected == 0 {
		regionDB = Region{Chromosome: region.Chromosome, Start: region.Start, End: region.End}
		DB.Create(&regionDB)
	}

	return regionDB.ID
}

func StoreGene(name string, ensemblID string) (geneID uint, created bool) {
	var geneDB Gene
	result := DB.Where("ensembl_id = ?", ensemblID).First(&geneDB)
	if result.RowsAffected == 0 {
		created = true
		geneDB = Gene{Name: name, EnsemblID: ensemblID}
		DB.Create(&geneDB)
	}

	return geneDB.ID, created
}

func StoreRegions(geneID uint, regions []ensembl.Region) {
	var regionsDB []Region

	for _, region := range regions {
		regionsDB = append(regionsDB, Region{GeneID: geneID, Chromosome: region.Chromosome, Start: region.Start, End: region.End, ExonNumber: region.ExonNumber})
	}
	DB.Create(&regionsDB)
}

func GetRegions() []Region {
	var regions []Region
	DB.Find(&regions)

	// karyotypic order
	// sort by chromosome then start position
	sort.SliceStable(regions, func(i, j int) bool {
		ci := utils.ChromosomeIndex(regions[i].Chromosome)
		cj := utils.ChromosomeIndex(regions[j].Chromosome)

		if ci < cj {
			return true
		}
		if ci > cj {
			return false
		}
		return regions[i].Start < regions[j].Start
	})

	return regions
}

func StoreKit(name string) (kitID uint, created bool) {
	var kitDB Kit
	result := DB.Where("d = ?", kitID).First(&kitDB)
	if result.RowsAffected == 0 {
		created = true
		kitDB = Kit{Name: name}
		DB.Create(&kitDB)
	}

	return kitDB.ID, created
}
