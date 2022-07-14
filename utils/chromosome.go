package utils

// Hardcoded chromosome lengths for GRCh38 - hg38
// https://www.ncbi.nlm.nih.gov/grc/human/data
var (
	CHROMOSOME_LENGTHS = [25]uint64{248956422, 242193529, 198295559, 190214555, 181538259,
		170805979, 159345973, 145138636, 138394717, 133797422,
		135086622, 133275309, 114364328, 107043718, 101991189,
		90338345, 83257441, 80373285, 58617616, 64444167,
		46709983, 50818468, 156040895, 57227415, 16569}
)

// ChromosomeIndex returns an index used for sorting the chromosomes based on karyotypic order
// https://stackoverflow.com/questions/46789259/map-vs-switch-performance-in-go
func ChromosomeIndex(chromosome string) int {
	switch chromosome {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "11":
		return 11
	case "12":
		return 12
	case "13":
		return 13
	case "14":
		return 14
	case "15":
		return 15
	case "16":
		return 16
	case "17":
		return 17
	case "18":
		return 18
	case "19":
		return 19
	case "20":
		return 20
	case "21":
		return 21
	case "22":
		return 22
	case "X":
		return 23
	case "Y":
		return 24
	case "MT":
		return 25
	default:
		return 99
	}
}
