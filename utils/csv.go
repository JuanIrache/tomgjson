package utils

import (
	"log"
	"strconv"
	"strings"
)

func splitStringsToFloats(xs []string) [][]float64 {
	xxf := [][]float64{}
	for _, s0 := range xs {
		if len(s0) > 0 {
			ss := strings.Split(s0, ",")
			for i, s := range ss {
				val, err := strconv.ParseFloat(s, 64)
				check(err)
				if len(xxf) < i+1 {
					xxf = append(xxf, []float64{})
				}
				xxf[i] = append(xxf[i], val)
			}
		}
	}
	return xxf
}

// ReadCSV formats a compatible CSV as a struct ready for mgJSON
func ReadCSV(src []byte, fr float64) SourceData {
	var data SourceData
	//To-Do check if can split bytes before converting to string
	lines := strings.Split(string(src), "\r\n")
	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0], 64); err != nil {
		header := lines[0]
		lines = lines[1:]
		splitStrings := splitStringsToFloats(lines)
		headers := strings.Split(header, ",")
		if headers[0] == "milliseconds" && len(headers[1]) > 1 {
			data.timing = splitStrings[0]
			splitStrings = splitStrings[1:]
			headers = headers[1:]
		}
		for i, vv := range splitStrings {
			data.streams = append(data.streams, stream{
				label:  headers[i],
				values: vv,
			})
		}
	} else {
		data.streams = []stream{{
			label:  "Data",
			values: splitStringsToFloats(lines)[0],
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
