package ensembl

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// https://www.ensembl.org/info/data/mysql.html
// https://www.ensembl.org/info/docs/api/core/core_schema.html

type Exon struct {
	HGNCAccession       string `db:"hgnc_accession"`
	GeneAccession       string `db:"gene_accession"`
	GeneName            string `db:"gene_name"`
	GeneDescription     string `db:"gene_description"`
	TranscriptAccession string `db:"transcript_accession"`
	RefSeqAccession     string `db:"refseq_accession"`
	Chromosome          string `db:"chromosome"`
	ExonNumber          uint   `db:"exon_number"`
	Strand              int    `db:"strand"`
	Start               uint64 `db:"start"`
	End                 uint64 `db:"end"`
	ExonAccession       string `db:"exon_accession"`
}

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "anonymous@(ensembldb.ensembl.org:3306)/homo_sapiens_core_106_38")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetExons(db *sqlx.DB) ([]Exon, error) {
	exons := []Exon{}
	query := `
		SELECT
		xref_gene.dbprimary_acc AS hgnc_accession,
		gene.stable_id AS gene_accession,
		xref_gene.display_label AS gene_name,
		xref_gene.description AS gene_description,
		transcript.stable_id AS transcript_accession,
		transcript_attrib.value AS refseq_accession,
		seq_region.name AS chromosome,
		exon_transcript.rank AS exon_number,
		exon.seq_region_strand AS strand,
		exon.seq_region_start AS start,
		exon.seq_region_end AS end,
		exon.stable_id AS exon_accession
		
		FROM exon
		
		INNER JOIN exon_transcript ON exon_transcript.exon_id = exon.exon_id
		INNER JOIN transcript ON transcript.transcript_id = exon_transcript.transcript_id
		INNER JOIN gene ON gene.gene_id = transcript.gene_id
		INNER JOIN transcript_attrib ON transcript_attrib.transcript_id = transcript.transcript_id AND attrib_type_id = 535
		INNER JOIN xref xref_gene ON xref_gene.xref_id = gene.display_xref_id AND xref_gene.external_db_id = 1100
		INNER JOIN seq_region ON seq_region.seq_region_id = gene.seq_region_id
		
		WHERE LEFT(gene.stable_id, 4) <> 'LRG_' AND xref_gene.display_label IS NOT NULL
		AND seq_region.name IN ("1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "X", "Y", "MT")
		
		ORDER BY chromosome, start
	`

	err := db.Select(&exons, query)
	if err != nil {
		return nil, err
	}

	return exons, nil
}

// TODO: synonyms
// SELECT xref_gene.display_label AS gene_name, synonym FROM external_synonym

// INNER JOIN gene on gene.display_xref_id = external_synonym.xref_id
// INNER JOIN xref xref_gene ON xref_gene.xref_id = gene.display_xref_id AND xref_gene.external_db_id = 1100

// WHERE xref_gene.display_label = "MAPK1";
