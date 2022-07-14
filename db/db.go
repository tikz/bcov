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

// ConnectDB handles the connection and migrations depending on the configured database engine.
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

// TODO: DELETE THIS feature<<<<<, change to MySQL for compatibility with already written manual queries
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

// ConnectSQLite manages the creation and connection to a local SQLite file.
// Includes extra PRAGMA statements and settings for improved write performance.
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

// StoreFile stores in the database an associated record to an original BAM file source
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

// StoreGene stores in the database a record for a gene
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

// StoreSynonyms stores in the database a record for a gene synonym mapping
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

// StoreExons stores in the database a record for a gene exon
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

// GetExons returns a slice of all exons in the database
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

// StoreExons stores in the database a single variant (rs) from ClinVar
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

// StoreKit stores in the database a commercial exon capture kit name and associated details
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
