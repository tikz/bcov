package api

import (
	"fmt"

	"github.com/tikz/bcov/db"
)

// queryExonVariants returns a list of variants that fall inside given a kit and exon database ID, and applying custom filters to the output
func queryExonVariants(kitId int, exonId int, filterSnp string, filterPathogenic bool) []Variant {
	var filterSnpQuery string
	if filterSnp == "" {
		filterSnpQuery = ""
	} else {
		filterSnpQuery = fmt.Sprintf("AND v.variant_id = %s", filterSnp)
	}

	var filterPathogenicQuery string
	if !filterPathogenic {
		filterPathogenicQuery = ""
	} else {
		filterPathogenicQuery = fmt.Sprintf(`AND v.clin_sig LIKE "%%%s%%"`, "pathogenic")
	}

	variants := make([]Variant, 0)
	db.DB.Raw(fmt.Sprintf(`
			SELECT variant_id as id, v.clin_sig, v.protein_change, v.chromosome, v.start, v.end, (SELECT * FROM (SELECT round(avg(rc.count))
			FROM read_counts rc
			INNER JOIN exon_read_counts erc on rc.exon_read_count_id = erc.id AND erc.exon_id = ?
			INNER JOIN bam_files bf on erc.bam_file_id = bf.id AND bf.kit_id = ?
			WHERE position >= v.start
			GROUP BY position
			ORDER BY position
			LIMIT 10) UNION ALL SELECT * FROM (SELECT round(avg(count))
			FROM read_counts rc
			INNER JOIN exon_read_counts erc on rc.exon_read_count_id = erc.id AND erc.exon_id = ?
			INNER JOIN bam_files bf on erc.bam_file_id = bf.id AND bf.kit_id = ?
			WHERE position < v.start
			GROUP BY position
			ORDER BY position DESC
			LIMIT 10)
			LIMIT 1) as depth FROM variants v
			
			WHERE v.exon_id = ? %s %s
			
			ORDER BY v.start
	`, filterSnpQuery, filterPathogenicQuery), exonId, kitId, exonId, kitId, exonId, filterSnp).Scan(&variants)

	return variants
}

// queryKitReads returns average exon reads for a given kit. Missing positions from the database are filled with zeroes.
func queryKitReads(exonId int, kitId int) []ReadCount {
	var exon db.Exon
	db.DB.Where("id = ?", exonId).First(&exon)

	var readCounts []ReadCount
	db.DB.Raw(`
			SELECT position, avg(count) avg_count FROM read_counts rc
			INNER JOIN exon_read_counts edc on rc.exon_read_count_id = edc.id
			INNER JOIN bam_files bf on edc.bam_file_id = bf.id
			WHERE edc.exon_id = ? AND bf.kit_id = ?
			GROUP BY rc.position
			ORDER BY position
	`, exonId, kitId).Scan(&readCounts)

	var start uint64
	if len(readCounts) == 0 {
		start = exon.Start
	} else {
		start = readCounts[0].Position
	}

	readCountsM := make(map[uint64]float64)
	for _, readCount := range readCounts {
		readCountsM[readCount.Position] = readCount.AvgCount
	}

	var filledReadCounts []ReadCount
	for i := uint64(0); i <= exon.End-start; i++ {
		if i%10 == 0 {
			if readCount, ok := readCountsM[start+i]; ok {
				filledReadCounts = append(filledReadCounts, ReadCount{Position: start + i, AvgCount: readCount})
			} else {
				filledReadCounts = append(filledReadCounts, ReadCount{Position: start + i, AvgCount: 0})
			}
		}
	}

	return filledReadCounts
}
