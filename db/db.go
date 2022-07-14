package db

import (
	"bcov/ensembl"
	"bcov/utils"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/glebarez/sqlite"
	"github.com/tikz/bio/clinvar"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	if os.Getenv("BCOV_DB_ENGINE") == "postgres" {
		ConnectPostgres()
		fmt.Println("Database engine:\t", "postgres")
	} else {
		ConnectSQLite()
		fmt.Println("Database engine:\t", "sqlite")
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
		log.Fatalf("failed to connect database: %s", err)
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

func StoreGene(hgncAccession string, geneAccession string, name string, description string, transcriptAccession string) (geneID uint, created bool) {
	var geneDB Gene
	result := DB.Where("gene_accession = ?", geneAccession).First(&geneDB)
	if result.RowsAffected == 0 {
		created = true
		geneDB = Gene{
			HGNCAccession:       hgncAccession,
			GeneAccession:       geneAccession,
			Name:                name,
			Description:         description,
			TranscriptAccession: transcriptAccession,
		}
		DB.Create(&geneDB)
	}

	return geneDB.ID, created
}

func StoreSynonyms(geneID uint, synonyms []ensembl.Synonym) (created bool) {
	var geneSynonymDB []GeneSynonym

	for _, synonym := range synonyms {
		geneSynonymDB = append(geneSynonymDB, GeneSynonym{
			GeneID:  geneID,
			Synonym: synonym.Synonym})
	}

	result := DB.Create(&geneSynonymDB)
	if result.RowsAffected == 0 {
		created = true
	}

	return created
}

func StoreExons(geneID uint, exons []ensembl.Exon) {
	var exonsDB []Exon

	for _, exon := range exons {
		exonsDB = append(exonsDB, Exon{
			GeneID:     geneID,
			Strand:     exon.Strand,
			Chromosome: exon.Chromosome,
			Start:      exon.Start,
			End:        exon.End,
			ExonNumber: exon.ExonNumber})
	}
	DB.Create(&exonsDB)
}

func GetExons() []Exon {
	var exons []Exon
	DB.Find(&exons)

	// karyotypic order
	// sort by chromosome then start position
	sort.SliceStable(exons, func(i, j int) bool {
		ci := utils.ChromosomeIndex(exons[i].Chromosome)
		cj := utils.ChromosomeIndex(exons[j].Chromosome)

		if ci < cj {
			return true
		}
		if ci > cj {
			return false
		}
		return exons[i].Start < exons[j].Start
	})

	return exons
}

func StoreVariant(allele clinvar.Allele, exonId uint) (variantID uint) {
	variant := Variant{
		VariantID:     allele.VariantID,
		Name:          allele.Name,
		ClinSig:       allele.ClinSig,
		ClinSigSimple: allele.ClinSigSimple,
		ProteinChange: allele.ProteinChange,
		ReviewStatus:  allele.ReviewStatus,
		Phenotypes:    allele.Phenotypes,
		Chromosome:    allele.Chromosome,
		Start:         allele.Start,
		End:           allele.End,
		ExonID:        exonId,
	}
	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&variant)

	return variant.ID
}

func StoreKit(name string) (kitID uint, created bool) {
	var kitDB Kit
	result := DB.Where("name = ?", name).First(&kitDB)
	if result.RowsAffected == 0 {
		created = true
		kitDB = Kit{Name: name}
		DB.Create(&kitDB)
	}

	return kitDB.ID, created
}
