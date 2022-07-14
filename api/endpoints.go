package api

import (
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
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

	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})

	return r
}
