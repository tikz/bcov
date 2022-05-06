package ensembl

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Region struct {
	Chromosome      string `db:"chrom"`
	Start           uint64 `db:"start"`
	End             uint64 `db:"end"`
	GeneHGNC        string `db:"gene_hgnc"`
	GeneName        string `db:"gene_name"`
	GeneDescription string `db:"gene_description"`
	ExonNumber      uint   `db:"exon_number"`
	StableID        string `db:"stable_id"`
}

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "anonymous@(ensembldb.ensembl.org:3306)/homo_sapiens_core_98_38")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetExons(db *sqlx.DB) ([]Region, error) {
	regions := []Region{}
	q := `
		SELECT DISTINCT
		seq_region.name AS chrom,
		exon.seq_region_start AS start,
		exon.seq_region_end AS end,
		gene_names.display_label AS gene_hgnc,
		gene_names.display_label AS gene_name,
		gene_names.display_label AS gene_description,
		exon_transcript.rank AS exon_number,
		gene.stable_id as stable_id
		FROM gene

		JOIN transcript ON transcript.transcript_id = gene.canonical_transcript_id
		JOIN exon_transcript ON exon_transcript.transcript_id = transcript.transcript_id
		JOIN exon ON exon_transcript.exon_id = exon.exon_id
		JOIN seq_region ON seq_region.seq_region_id = gene.seq_region_id

		LEFT JOIN (
			SELECT xref_id, display_label
			FROM xref
			WHERE xref.external_db_id = 1100
		) AS gene_names ON gene_names.xref_id = gene.display_xref_id

		WHERE LEFT(gene.stable_id, 4) <> 'LRG_' AND gene_names.display_label IS NOT NULL
		AND seq_region.name IN ("1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "X", "Y", "MT")

		ORDER BY CAST(seq_region.name as unsigned), exon.seq_region_start;
	`

	err := db.Select(&regions, q)
	if err != nil {
		return nil, err
	}

	// for _, region := range regions {
	// 	fmt.Printf("%#v\n", region)
	// }
	// fmt.Println(len(regions))
	// fmt.Println(err)

	return regions, nil
}
