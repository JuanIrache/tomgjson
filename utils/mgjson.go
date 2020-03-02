package utils

type stream struct {
	label  string
	values []float64
}

// SourceData structures data relevant for mgJSON files
type SourceData struct {
	timing  []float64
	streams []stream
}

type mgjson struct {
}

func FormatMgjson(sd SourceData) mgjson {
	var data mgjson

	return data
}
