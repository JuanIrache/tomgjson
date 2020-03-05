package tomgjson

import (
	"encoding/csv"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
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

func floatToTime(f float64) time.Time {
	seconds := f / 1000
	fullSeconds := math.Floor(seconds)
	nanoseconds := (seconds - fullSeconds) * 1e+9
	return time.Unix(int64(fullSeconds), int64(nanoseconds))
}

func floatsToTimes(xf []float64) []time.Time {
	xt := []time.Time{}
	for _, f := range xf {
		mTime := floatToTime(f)
		xt = append(xt, mTime)
	}
	return xt
}

// FromCSV formats a compatible CSV as a struct ready for mgJSON. The optional frame rate (fr) is only used if timing data is not present
func FromCSV(src []byte, fr float64) FormattedData {
	var data FormattedData

	r := csv.NewReader(strings.NewReader(string(src)))
	lines, err := r.ReadAll()
	check(err)

	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0][0], 64); err != nil {
		headers := lines[0]
		lines = lines[1:]
		floatsTable := stringsTableToFloats(lines)
		if headers[0] == "milliseconds" && len(headers[1]) > 1 {
			data.timing = floatsToTimes(floatsTable[0])
			floatsTable = floatsTable[1:]
			headers = headers[1:]
		}
		for i, vv := range floatsTable {
			data.streams = append(data.streams, Stream{
				label:  headers[i],
				values: vv,
			})
		}
	} else {
		data.streams = []Stream{{
			label:  "Data",
			values: stringsTableToFloats(lines)[0],
		}}
	}

	if len(data.streams[0].values) < 1 {
		log.Panic("No valid data found")
	}

	if len(data.timing) < 1 {
		for i := 0; i < len(data.streams[0].values); i++ {
			data.timing = append(data.timing, floatToTime(fr*float64(i)))
		}
	}

	return data
}