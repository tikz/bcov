package cov

import (
	"bcov/db"
)

type Position uint64
type Depth uint64
type Coverage uint8

type Exon struct {
	Chromosome     string
	Start          uint64
	End            uint64
	PositionDepth  map[Position]Depth
	DepthCoverages map[Depth]Coverage

	ID         uint
	GeneID     uint
	ExonNumber uint
}

// NewExon instantiates a new Exon struct
func NewExon(chromosome string, start uint64, end uint64) Exon {
	return Exon{
		Chromosome:    chromosome,
		Start:         start,
		End:           end,
		PositionDepth: make(map[Position]Depth),
	}
}

// NewExonsFromDB returns a slice of all the exons that are already loaded in the database.
func NewExonsFromDB() []Exon {
	dbExons := db.GetExons()
	var exons []Exon
	for _, exon := range dbExons {
		exons = append(exons, Exon{
			ID:            exon.ID,
			GeneID:        exon.GeneID,
			Chromosome:    exon.Chromosome,
			Start:         exon.Start,
			End:           exon.End,
			ExonNumber:    exon.ExonNumber,
			PositionDepth: make(map[Position]Depth),
		})
	}

	return exons
}

// AddDepth increases the depth counter for a given position
// If the position falls outside the exon it isn't counted.
func (r *Exon) AddDepth(pos Position, depth Depth) {
	if uint64(pos) < r.Start || uint64(pos) > r.End {
		return
	}
	if _, ok := r.PositionDepth[pos]; !ok {
		r.PositionDepth[pos] = 0
	}
	r.PositionDepth[pos] += depth
}

// AddDepthFromTo increases the depth counter for a given position range.
// Positions that fall outside the exon aren't counted.
func (r *Exon) AddDepthFromTo(fromPos Position, toPos Position, depth Depth) {
	start, end := max(r.Start, uint64(fromPos)), min(r.End, uint64(toPos))
	for i := start; i <= end; i++ {
		r.AddDepth(Position(i), depth)
	}
}

// ComputeDepthCoverage computes the coverage at the given depth.
func (r *Exon) ComputeDepthCoverage(depth Depth) Coverage {
	count := 0
	for _, posDepth := range r.PositionDepth {
		if Depth(posDepth) >= depth {
			count++
		}
	}

	return Coverage(100 * float64(count) / float64(r.End-r.Start+1))
}

// ComputeDepthCoverageRange computes the coverage at the given depths range.
// Stores the result in Exon.DepthCoverages.
func (r *Exon) ComputeDepthCoverageRange(fromDepth Depth, toDepth Depth) {
	r.DepthCoverages = make(map[Depth]Coverage)
	r.DepthCoverages[1] = r.ComputeDepthCoverage(1)
	for i := fromDepth; i <= toDepth; i += 10 {
		r.DepthCoverages[i] = r.ComputeDepthCoverage(i)
	}
}

// StoreDepthCoverages computes and stores the depth coverages in the DB, under the given BAM file ID.
func (r *Exon) StoreDepthCoverages(bamFileID uint) {
	r.ComputeDepthCoverageRange(10, 100)
	var depthCoverages []db.DepthCoverage

	for depth, coverage := range r.DepthCoverages {
		depthCoverages = append(depthCoverages,
			db.DepthCoverage{Depth: uint8(depth),
				Coverage: uint8(coverage),
			})
	}

	db.DB.Create(&db.ExonDepthCoverage{
		ExonID:         r.ID,
		BAMFileID:      bamFileID,
		DepthCoverages: depthCoverages,
	})
}

// StoreReadCounts stores the read counts in the DB, under the given BAM file ID.
func (r *Exon) StoreReadCounts(bamFileID uint) {

	var readCounts []db.ReadCount

	for position, depth := range r.PositionDepth {
		if position%10 == 0 {
			readCounts = append(readCounts, db.ReadCount{Position: uint64(position), Count: uint64(depth)})
		}
	}

	db.DB.Create(&db.ExonReadCount{
		ExonID:     r.ID,
		BAMFileID:  bamFileID,
		ReadCounts: readCounts,
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
