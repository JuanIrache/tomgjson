package utils

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"
)

func stringsTableToFloats(xxs [][]string) [][]float64 {
	xxf := [][]float64{}
	for _, xs := range xxs {
		for i, s := range xs {
			val, err := strconv.ParseFloat(s, 64)
			check(err)
			if len(xxf) < i+1 {
				xxf = append(xxf, []float64{})
			}
			xxf[i] = append(xxf[i], val)
		}
	}
	return xxf
}

// ReadCSV formats a compatible CSV as a struct ready for mgJSON. The optional frame rate (fr) is only used if timing data is not present
func ReadCSV(src []byte, fr float64) SourceData {
	var data SourceData

	r := csv.NewReader(strings.NewReader(string(src)))
	lines, err := r.ReadAll()
	check(err)

	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0][0], 64); err != nil {
		headers := lines[0]
		lines = lines[1:]
		floatsTable := stringsTableToFloats(lines)
		if headers[0] == "milliseconds" && len(headers[1]) > 1 {
			data.timing = floatsTable[0]
			floatsTable = floatsTable[1:]
			headers = headers[1:]
		}
		for i, vv := range floatsTable {
			data.streams = append(data.streams, stream{
				label:  headers[i],
				values: vv,
			})
		}
	} else {
		data.streams = []stream{{
			label:  "Data",
			values: stringsTableToFloats(lines)[0],
		}}
	}

	if len(data.streams[0].values) < 1 {
		log.Panic("No valid data found")
	}

	if len(data.timing) < 1 {
		for i := 0; i < len(data.streams[0].values); i++ {
			data.timing = append(data.timing, float64(i)*fr)
		}
	}

	return data
}
