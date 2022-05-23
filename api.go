package main

import (
	"bcov/db"
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Variant struct {
	VariantIDs    string `json:"variantIds"`
	ClinSig       string `json:"clinSig"`
	ProteinChange string `json:"proteinChange"`
	Chromosome    string `json:"chromosome"`
	Start         uint64 `json:"start"`
	End           uint64 `json:"end"`
	Depth         uint64 `json:"depth"`
}

type VariantsResult struct {
	TotalCount  int       `json:"totalCount"`
	Pages       int       `json:"pages"`
	CurrentPage int       `json:"currentPage"`
	Variants    []Variant `json:"variants"`
}

func KitsEndpoint(c *gin.Context) {
	var kits []db.Kit
	result := db.DB.Find(&kits)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, kits)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func GeneEndpoint(c *gin.Context) {
	id := c.Param("id")

	var gene db.Gene
	result := db.DB.Where("id = ?", id).Preload("Exons").First(&gene)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gene)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func ReadsEndpoint(c *gin.Context) {
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))

	var exon db.Exon
	db.DB.Where("id = ?", exonId).First(&exon)

	type ReadCount struct {
		Position uint64  `json:"position"`
		AvgCount float64 `json:"avgCount"`
	}

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

	var kit db.Kit
	db.DB.Where("id = ?", kitId).First(&kit)

	type Response struct {
		KitName    string      `json:"kitName"`
		ReadCounts []ReadCount `json:"readCounts"`
	}

	c.JSON(http.StatusOK, Response{KitName: kit.Name, ReadCounts: filledReadCounts})
}

func KitEndpoint(c *gin.Context) {
	id := c.Param("id")

	var kit db.Kit
	result := db.DB.Where("id = ?", id).First(&kit)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, kit)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func SearchGenesEndpoint(c *gin.Context) {
	name := c.Param("name")
	if len(name) < 3 {
		c.JSON(http.StatusNotFound, gin.H{"error": "too few characters to search"})
		return
	}

	var genes []db.Gene

	// Exact matchs
	db.DB.Where("gene_accession = ? OR hgnc_accession = ? OR transcript_accession = ?", name, name, name).Find(&genes)
	if len(genes) > 0 {
		c.JSON(http.StatusOK, genes)
		return
	}

	search := fmt.Sprintf("%%%s%%", name)
	db.DB.Where("name LIKE ? OR description LIKE ?", search, search).Find(&genes)

	// Prioritize sorting first by gene name, then by description
	nameUpper := strings.ToUpper(name)
	sort.SliceStable(genes, func(i, j int) bool {
		return strings.Index(genes[i].Name, nameUpper) > strings.Index(genes[j].Name, nameUpper)
	})
	if len(genes) > 20 {
		genes = genes[:20]
	}

	c.JSON(http.StatusOK, genes)
}

func SearchKitsEndpoint(c *gin.Context) {
	name := c.Param("name")
	if len(name) < 3 {
		c.JSON(http.StatusNotFound, gin.H{"error": "too few characters to search"})
		return
	}

	var kits []db.Kit

	if name == ":kits" {
		db.DB.Find(&kits)
		c.JSON(http.StatusOK, kits)
		return
	}

	search := fmt.Sprintf("%%%s%%", name)
	db.DB.Where("name LIKE ? ", search).Find(&kits)
	c.JSON(http.StatusOK, kits)
}

func DepthCoveragesEndpoint(c *gin.Context) {
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))

	var exon db.Exon
	db.DB.Where("id = ?", exonId).First(&exon)

	type DepthCoverage struct {
		Depth    uint64  `json:"depth"`
		Coverage float64 `json:"coverage"`
	}

	var depthCoverages []DepthCoverage
	db.DB.Raw(`
			SELECT depth, avg(coverage) coverage from depth_coverages dc
			INNER JOIN exon_depth_coverages edc on edc.id = dc.exon_depth_coverage_id
			INNER JOIN bam_files bf on bf.id = edc.bam_file_id
			INNER JOIN kits k on k.id = bf.kit_id
			
			WHERE exon_id = ? AND kit_id = ?
			GROUP BY depth
	`, exonId, kitId).Scan(&depthCoverages)

	var kit db.Kit
	db.DB.Where("id = ?", kitId).First(&kit)

	type Response struct {
		KitName        string          `json:"kitName"`
		DepthCoverages []DepthCoverage `json:"depthCoverages"`
	}

	c.JSON(http.StatusOK, Response{KitName: kit.Name, DepthCoverages: depthCoverages})
}

func queryExonVariants(kitId int, exonId int) []Variant {
	variants := make([]Variant, 0)
	db.DB.Raw(`
			SELECT GROUP_CONCAT(v.variant_id) variant_ids, v.clin_sig, v.protein_change, v.chromosome, v.start, v.end, (SELECT * FROM (SELECT round(avg(rc.count))
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
			
			WHERE v.exon_id = ?
			
			GROUP BY start, end
			ORDER BY v.start
	`, exonId, kitId, exonId, kitId, exonId).Scan(&variants)

	return variants
}

func VariantsEndpoint(c *gin.Context) {
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))
	pageParam := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageParam)
	perPage := 20

	variants := queryExonVariants(kitId, exonId)

	start := (page - 1) * perPage
	if start > len(variants) {
		start = len(variants) - 1
	}
	end := page * perPage
	if end > len(variants) {
		end = len(variants)
	}

	c.JSON(http.StatusOK, VariantsResult{TotalCount: len(variants), Pages: len(variants)/perPage + 1, CurrentPage: page, Variants: variants[start:end]})
}

func VariantsCSVEndpoint(c *gin.Context) {
	name := c.Param("gene_name")

	var kits []db.Kit
	db.DB.Find(&kits)

	var gene db.Gene
	db.DB.Where("name = ?", name).Preload("Exons").Find(&gene)

	variants := make([][]Variant, 0)

	for _, kit := range kits {
		kitVariants := make([]Variant, 0)
		for _, exon := range gene.Exons {
			kitVariants = append(kitVariants, queryExonVariants(int(kit.ID), int(exon.ID))...)
		}

		variants = append(variants, kitVariants)
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"dbSNP ID", "Clinical significance", "Protein change", "Chromosome", "Start", "End"}
	for _, kit := range kits {
		header = append(header, kit.Name+" depth")
	}
	writer.Write(header)

	for i, variant := range variants[0] {
		line := []string{"rs" + variant.VariantIDs, variant.ClinSig, variant.ProteinChange, variant.Chromosome, fmt.Sprint(variant.Start), fmt.Sprint(variant.End)}
		for j := range kits {
			line = append(line, fmt.Sprint(variants[j][i].Depth))
		}
		writer.Write(line)
	}

	writer.Flush()

	filename := fmt.Sprintf("%s_VariantsDepth.csv", name)
	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.String(http.StatusOK, buf.String())
}

func BAMsEndpoint(c *gin.Context) {
	id := c.Param("kit_id")

	var bams []db.BAMFile
	result := db.DB.Where("kit_id = ?", id).Find(&bams)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, bams)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func runWebServer() {
	db.ConnectDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.Use(static.Serve("/", static.LocalFile("web/build", true)))

	r.GET("/api/kits", KitsEndpoint)
	r.GET("/api/kit/:id", GeneEndpoint)
	r.GET("/api/gene/:id", GeneEndpoint)

	r.GET("/api/search/genes/:name", SearchGenesEndpoint)
	r.GET("/api/search/kits/:name", SearchKitsEndpoint)

	r.GET("/api/reads/:kit_id/:exon_id", ReadsEndpoint)
	r.GET("/api/depth-coverages/:kit_id/:exon_id", DepthCoveragesEndpoint)
	r.GET("/api/variants/:kit_id/:exon_id", VariantsEndpoint)
	r.GET("/api/bams/:kit_id", BAMsEndpoint)
	r.GET("/api/variants-csv/:gene_name", VariantsCSVEndpoint)

	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})
	r.Run()
}
