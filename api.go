package main

import (
	"bcov/db"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GeneEndpoint(c *gin.Context) {
	name := c.Param("name")

	var gene db.Gene
	result := db.DB.Where("name = ?", name).Preload("Regions").First(&gene)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gene)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
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

	r.GET("/api/gene/:name", GeneEndpoint)
	r.GET("/api/depth-coverages/kit/:kitID/region/:regionID", DepthCoveragesEndpoint)

	// r.NoRoute(func(c *gin.Context) {
	// 	if !strings.HasPrefix(c.Request.RequestURI, "/auth") && !strings.HasPrefix(c.Request.RequestURI, "/admin") && !strings.HasSuffix(c.Request.RequestURI, ".map") {
	// 		filePath := "web/dist" + c.Request.RequestURI
	// 		if _, err := os.Stat(filePath); err == nil {
	// 			c.File(filePath)
	// 		} else {
	// 			c.File("web/dist/index.html")
	// 		}
	// 	}
	// })
	r.Run()
}
