package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/tikz/bcov/db"

	"github.com/gin-gonic/gin"
)

func Endpoints() *gin.Engine {
	cacheDuration := 24 * 7 * time.Hour
	memoryStore := persist.NewMemoryStore(cacheDuration)

	r := gin.Default()
	r.Use(cors.Default())

	r.Use(static.Serve("/", static.LocalFile("web/build", true)))

	r.GET("/api/kits", KitsEndpoint)
	r.GET("/api/kit/:id", KitEndpoint)
	r.GET("/api/gene/:id", GeneEndpoint)

	r.GET("/api/search/genes/:name", cache.CacheByRequestURI(memoryStore, cacheDuration), SearchGenesEndpoint)
	r.GET("/api/search/kits/:name", cache.CacheByRequestURI(memoryStore, cacheDuration), SearchKitsEndpoint)
	r.GET("/api/search/variant/:id", cache.CacheByRequestURI(memoryStore, cacheDuration), SearchVariantEndpoint)

	r.GET("/api/reads/:kit_id/:exon_id", cache.CacheByRequestURI(memoryStore, cacheDuration), ReadsEndpoint)
	r.GET("/api/depth-coverages/:kit_id/:exon_id", cache.CacheByRequestURI(memoryStore, cacheDuration), DepthCoveragesEndpoint)
	r.GET("/api/variants/:kit_id/:exon_id", cache.CacheByRequestURI(memoryStore, cacheDuration), VariantsEndpoint)
	r.GET("/api/bams/:kit_id", BAMsEndpoint)
	r.GET("/api/variants-csv/:gene_name", cache.CacheByRequestURI(memoryStore, cacheDuration), VariantsCSVEndpoint)
	r.GET("/api/gene-coverage/:gene_id", GeneCoverageEndpoint)

	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})

	return r
}

// Endpoints logic

// KitsEndpoint handles the response for all the kits in the database.
func KitsEndpoint(c *gin.Context) {
	var kits []db.Kit
	result := db.DB.Find(&kits)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, kits)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

// GeneEndpoint handles the response for details about a given gene passed by database ID.
func GeneEndpoint(c *gin.Context) {
	// API inputs
	id := c.Param("id")
	////

	var gene db.Gene
	result := db.DB.Where("id = ?", id).Preload("Exons").First(&gene)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gene)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

// ReadsEndpoint returns all reads for a given kit ID and exon ID.
func ReadsEndpoint(c *gin.Context) {
	// API inputs
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))
	////

	readCounts := queryKitReads(exonId, kitId)

	var kit db.Kit
	db.DB.Where("id = ?", kitId).First(&kit)

	c.JSON(http.StatusOK, ReadCountsResponse{KitName: kit.Name, ReadCounts: readCounts})
}

// KitEndpoint returns details about a given kit passed by database ID.
func KitEndpoint(c *gin.Context) {
	// API inputs
	id := c.Param("id")
	////

	var kit db.Kit
	result := db.DB.Where("id = ?", id).First(&kit)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

// SearchGenesEndpoint handles the response for a single gene that matches exactly or partially the HGNC name.
func SearchGenesEndpoint(c *gin.Context) {
	// API inputs
	name := c.Param("name")
	////

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

	// Limit number of results
	if len(genes) > 20 {
		genes = genes[:20]
	}

	c.JSON(http.StatusOK, genes)
}

// SearchKitsEndpoint handles the response for details about a given kit that matches exactly or partially a given name.
func SearchKitsEndpoint(c *gin.Context) {
	// API inputs
	name := c.Param("name")
	////

	if len(name) < 3 {
		c.JSON(http.StatusNotFound, gin.H{"error": "too few characters to search"})
		return
	}

	var kits []db.Kit

	if name == "@kits" {
		db.DB.Find(&kits)
		c.JSON(http.StatusOK, kits)
		return
	}

	search := fmt.Sprintf("%%%s%%", name)
	db.DB.Where("name LIKE ? ", search).Find(&kits)
	c.JSON(http.StatusOK, kits)
}

// SearchVariantEndpoint handles the response for details about a given variant by ID.
func SearchVariantEndpoint(c *gin.Context) {
	// API inputs
	id, _ := strconv.Atoi(strings.ReplaceAll(c.Param("id"), "rs", ""))
	////

	variants := make([]VariantSearch, 0)
	db.DB.Raw(`
				SELECT v.variant_id as id, e.exon_number, g.name as gene_name, g.id as gene_id FROM variants v
				INNER JOIN exons e on e.id = v.exon_id
				INNER JOIN genes g on g.id = e.gene_id

				WHERE v.variant_id = ?
	`, id).Scan(&variants)

	for i := range variants {
		var gene db.Gene
		db.DB.Where("id = ?", variants[i].GeneID).First(&gene)
		variants[i].Gene = gene
	}

	c.JSON(http.StatusOK, variants)
}

// DepthCoveragesEndpoint handles the response for depth coverage for a given exon ID and kit ID.
func DepthCoveragesEndpoint(c *gin.Context) {
	// API inputs
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))
	////

	var exon db.Exon
	db.DB.Where("id = ?", exonId).First(&exon)

	depthCoverages := queryExonDepthCoverages(exonId, kitId)

	var kit db.Kit
	db.DB.Where("id = ?", kitId).First(&kit)

	c.JSON(http.StatusOK, DepthCoveragesResponse{KitName: kit.Name, DepthCoverages: depthCoverages})
}

// BAMsEndpoint handles the response for all loaded BAM files of a given kit ID.
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

// VariantsEndpoint handles the response for a list of variants that fall inside given a kit ID and exon ID.
func VariantsEndpoint(c *gin.Context) {
	// API inputs
	kitId, _ := strconv.Atoi(c.Param("kit_id"))
	exonId, _ := strconv.Atoi(c.Param("exon_id"))
	filterId := c.DefaultQuery("filter_id", "")
	pageParam := c.DefaultQuery("page", "1")
	pathogenicParam, _ := strconv.ParseBool(c.DefaultQuery("pathogenic", "0"))
	page, _ := strconv.Atoi(pageParam)
	perPage := 20
	////

	variants := queryExonVariants(kitId, exonId, filterId, pathogenicParam)

	start := (page - 1) * perPage
	if start > len(variants) {
		start = len(variants) - 1
	}
	end := page * perPage
	if end > len(variants) {
		end = len(variants)
	}

	c.JSON(http.StatusOK, VariantsResult{TotalCount: len(variants), Pages: int(math.Ceil(float64(len(variants)) / float64(perPage))), CurrentPage: page, Variants: variants[start:end]})
}

// VariantsCSVEndpoint handles the response for a given HGNC gene name returning a CSV file.
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
			kitVariants = append(kitVariants, queryExonVariants(int(kit.ID), int(exon.ID), "", false)...)
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
		line := []string{"rs" + variant.ID, variant.ClinSig, variant.ProteinChange, variant.Chromosome, fmt.Sprint(variant.Start), fmt.Sprint(variant.End)}
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

// GeneCoverageEndpoint returns average coverage statistics for a given gene ID.
func GeneCoverageEndpoint(c *gin.Context) {
	// API inputs
	geneId, _ := strconv.Atoi(c.Param("gene_id"))
	kitIdsInput := strings.Split(c.Query("kit_ids"), ",")

	// Parse list of kit ids from query
	var kitIds []int
	for _, kitIdStr := range kitIdsInput {
		kitId, err := strconv.Atoi(kitIdStr)
		if err == nil {
			kitIds = append(kitIds, kitId)
		}
	}

	var kits []KitGeneDepthCoverage

	kitGeneCoverages := queryGeneDepthCoverage(geneId)
	for _, kitGeneCoverage := range kitGeneCoverages {
		for _, kitId := range kitIds {
			if uint64(kitId) == kitGeneCoverage.ID {
				kitGeneCoverage.Variants = queryGeneVariantsDepth(kitId, geneId)
				kits = append(kits, kitGeneCoverage)
			}
		}
	}

	c.JSON(http.StatusOK, kits)
}
