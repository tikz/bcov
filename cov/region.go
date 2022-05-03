package cov

import (
	"bcov/db"
)

type Position uint64
type Depth uint64
type Coverage uint8

type Region struct {
	Chromosome     string
	Start          uint64
	End            uint64
	PositionDepth  map[Position]Depth
	DepthCoverages map[Depth]Coverage

	ID         uint
	GeneID     uint
	ExonNumber uint
}

func NewRegion(chromosome string, start uint64, end uint64) Region {
	return Region{
		Chromosome:    chromosome,
		Start:         start,
		End:           end,
		PositionDepth: make(map[Position]Depth),
	}
}

func NewRegionsFromDB() []Region {
	dbRegions := db.GetRegions()
	var regions []Region
	for _, region := range dbRegions {
		regions = append(regions, Region{
			ID:            region.ID,
			GeneID:        region.GeneID,
			Chromosome:    region.Chromosome,
			Start:         region.Start,
			End:           region.End,
			ExonNumber:    region.ExonNumber,
			PositionDepth: make(map[Position]Depth),
		})
	}

	return regions
}

// AddDepth increases the depth counter for a given position
// If the position falls outside the region it isn't counted.
func (r *Region) AddDepth(pos Position, depth Depth) {
	if uint64(pos) < r.Start || uint64(pos) > r.End {
		return
	}
	if _, ok := r.PositionDepth[pos]; !ok {
		r.PositionDepth[pos] = 0
	}
	r.PositionDepth[pos] += depth
}

// AddDepthFromTo increases the depth counter for a given position range.
// Positions that fall outside the region aren't counted.
func (r *Region) AddDepthFromTo(fromPos Position, toPos Position, depth Depth) {
	start, end := max(r.Start, uint64(fromPos)), min(r.End, uint64(toPos))
	for i := start; i <= end; i++ {
		r.AddDepth(Position(i), depth)
	}
}

// ComputeDepthCoverage computes the coverage at the given depth.
func (r *Region) ComputeDepthCoverage(depth Depth) Coverage {
	count := 0
	for _, posDepth := range r.PositionDepth {
		if Depth(posDepth) >= depth {
			count++
		}
	}

	return Coverage(100 * float64(count) / float64(r.End-r.Start+1))
}

// ComputeDepthCoverageRange computes the coverage at the given depths range.
// Stores the result in Region.DepthCoverages.
func (r *Region) ComputeDepthCoverageRange(fromDepth Depth, toDepth Depth) {
	r.DepthCoverages = make(map[Depth]Coverage)
	r.DepthCoverages[1] = r.ComputeDepthCoverage(1)
	for i := fromDepth; i <= toDepth; i += 10 {
		r.DepthCoverages[i] = r.ComputeDepthCoverage(i)
	}
}

// StoreDepthCoverages computes and stores the depth coverages in the DB, under the given BAM file ID.
func (r *Region) StoreDepthCoverages(bamFileID uint) {
	r.ComputeDepthCoverageRange(10, 100)
	var depthCoverages []db.DepthCoverage

	for depth, coverage := range r.DepthCoverages {
		depthCoverages = append(depthCoverages,
			db.DepthCoverage{Depth: uint8(depth),
				Coverage: uint8(coverage),
			})
	}

	db.DB.Create(&db.RegionDepthCoverage{
		RegionID:       r.ID,
		BAMFileID:      bamFileID,
		DepthCoverages: depthCoverages,
	})
}

func max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func min(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
