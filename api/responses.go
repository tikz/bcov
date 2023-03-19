package api

import "github.com/tikz/bcov/db"

type ReadCount struct {
	Position uint64  `json:"position"`
	AvgCount float64 `json:"avgCount"`
}

type ReadCountsResponse struct {
	KitName    string      `json:"kitName"`
	ReadCounts []ReadCount `json:"readCounts"`
}

type DepthCoverage struct {
	Depth    uint64  `json:"depth"`
	Coverage float64 `json:"coverage"`
}

type DepthCoveragesResponse struct {
	KitName        string          `json:"kitName"`
	DepthCoverages []DepthCoverage `json:"depthCoverages"`
}

type KitGeneDepthCoverage struct {
	ID       uint64        `json:"id"`
	Name     string        `json:"name"`
	Depth    uint64        `json:"depth"`
	Coverage float64       `json:"coverage"`
	Variants VariantsDepth `json:"variants" gorm:"-"`
}

type VariantSearch struct {
	ID         uint64  `json:"id"`
	Gene       db.Gene `json:"gene"`
	GeneName   string  `json:"geneName"`
	GeneID     uint64  `json:"geneId"`
	ExonNumber int     `json:"exonNumber"`
}
type Variant struct {
	ID            string `json:"id"`
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

type VariantDepth struct {
	ID    int     `json:"id"`
	Depth float64 `json:"depth"`
}

type VariantsDepth struct {
	Total             int            `json:"total"`
	Covered           int            `json:"covered"`
	UncoveredVariants []VariantDepth `json:"uncoveredVariants"`
}
