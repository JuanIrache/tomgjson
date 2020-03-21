package tomgjson

import (
	"encoding/csv"
	"fmt"
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

func millisecondsToTime(f float64) time.Time {
	seconds := f / 1000
	fullSeconds := math.Floor(seconds)
	nanoseconds := (seconds - fullSeconds) * 1e+9
	return time.Unix(int64(fullSeconds), int64(nanoseconds))
}

func floatsToTimes(xf []float64) []time.Time {
	xt := []time.Time{}
	for _, f := range xf {
		mTime := millisecondsToTime(f)
		xt = append(xt, mTime)
	}
	return xt
}

// FromCSV formats a compatible CSV as a FormattedData struct ready for mgJSON and returns it. Or returns an error
// The optional frame rate (fr) is used if timing data is not present
func FromCSV(src []byte, fr float64) (FormattedData, error) {
	var data FormattedData

	r := csv.NewReader(strings.NewReader(string(src)))
	lines, err := r.ReadAll()
	if err != nil {
		return data, fmt.Errorf("Error reading CSV: %v", err.Error())
	}

	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0][0], 64); err != nil {
		headers := lines[0]
		lines = lines[1:]
		floatsTable := stringsTableToFloats(lines)
		if headers[0] == "milliseconds" && len(headers[1]) > 1 {
			data.Timing = floatsToTimes(floatsTable[0])
			floatsTable = floatsTable[1:]
			headers = headers[1:]
		}
		for i, vv := range floatsTable {
			data.Streams = append(data.Streams, Stream{
				Label:  headers[i],
				Values: vv,
			})
		}
	} else {
		data.Streams = []Stream{{
			Label:  "Data",
			Values: stringsTableToFloats(lines)[0],
		}}
	}

	if len(data.Streams[0].Values) < 1 {
		return data, fmt.Errorf("No valid data found")
	}

	if len(data.Timing) < 1 {
		for i := 0; i < len(data.Streams[0].Values); i++ {
			data.Timing = append(data.Timing, millisecondsToTime(float64(i)*1000.0/fr))
		}
	}

	return data, nil
}
