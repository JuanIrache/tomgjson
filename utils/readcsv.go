package readcsv

import (
	"log"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// SourceData structures data relevant for mgJSON files
type SourceData struct {
	label  string
	timing []float64
	values []float64
}

func stringsToFloats(xs []string) []float64 {
	xf := []float64{}
	for _, s := range xs {
		if len(s) > 0 {
			val, err := strconv.ParseFloat(s, 64)
			check(err)
			xf = append(xf, val)
		}
	}
	return xf
}

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

// ReadCSV formats a compatible CSV as a struct
func ReadCSV(src []byte) SourceData {
	var data SourceData
	//To-Do check if can split bytes before converting to string
	lines := strings.Split(string(src), "\r\n")
	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0], 64); err != nil {
		header := lines[0]
		lines = lines[1:]
		if headers := strings.Split(header, ","); headers[0] == "milliseconds" && len(headers[1]) > 1 {
			//Have timing data
			data.label = headers[1]
			splitStrings := splitStringsToFloats(lines)
			data.timing = splitStrings[0]
			data.values = splitStrings[1]
		} else {
			data.label = header
			data.values = stringsToFloats(lines)
		}
	} else {
		//To-Do use file name?
		data.label = "Data"
		data.values = stringsToFloats(lines)
	}

	//Fill timing if missing
	return data
}