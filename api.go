package main

import (
	"bcov/db"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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
	name := c.Param("id")

	var gene db.Gene
	result := db.DB.Where("id = ?", name).Preload("Regions").First(&gene)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gene)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
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
	search := fmt.Sprintf("%%%s%%", name)
	db.DB.Where("name LIKE ? ", search).Find(&kits)
	c.JSON(http.StatusOK, kits)
}

// SELECT k.id, k.name, dc.depth, avg(dc.coverage) as coverage FROM depth_coverages dc
// INNER JOIN region_depth_coverages rdc on dc.region_depth_coverage_id = rdc.id
// INNER JOIN regions r on rdc.region_id = r.id
// INNER JOIN genes g on r.gene_id = g.id
// INNER JOIN bam_files bf on rdc.bam_file_id = bf.id
// INNER JOIN kits k on bf.kit_id = k.id
// WHERE g.id = 32849
// GROUP BY dc.depth, k.id

// SELECT k.id, k.name, dc.depth, avg(dc.coverage) as coverage FROM depth_coverages dc
// INNER JOIN region_depth_coverages rdc on dc.region_depth_coverage_id = rdc.id
// INNER JOIN bam_files bf on rdc.bam_file_id = bf.id
// INNER JOIN kits k on bf.kit_id = k.id
// WHERE rdc.region_id = 209642
// GROUP BY dc.depth, k.id
func DepthCoveragesEndpoint(c *gin.Context) {
	kitID := c.Param("kitID")
	regionID := c.Param("regionID")
	fmt.Println(kitID)

	var regionDepthCoverages []db.ExonDepthCoverage
	result := db.DB.Where("region_id = ?", regionID).Preload("DepthCoverages").Preload("BAMFile").Preload("BAMFile.Kit").Find(&regionDepthCoverages)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, regionDepthCoverages)
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
	r.GET("/api/depth-coverages/kit/:kitID/region/:regionID", DepthCoveragesEndpoint)

	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})
	r.Run()
}
