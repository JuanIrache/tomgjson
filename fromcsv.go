package tomgjson

import (
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"bytes"
)

// Returns valid streams with values or strings
func structureData(headers []string, table [][]string) ([]Stream, error) {

	streams := []Stream{}

	for _, xs := range table {
		for i, s := range xs {
			if len(streams) < i+1 {
				streams = append(streams, Stream{
					Label: headers[i],
				})
			}
			val, err := strconv.ParseFloat(s, 64)
			if err == nil {
				streams[i].Values = append(streams[i].Values, val)
			} else {
				if len(streams[i].Values) > 0 {
					return streams, fmt.Errorf("Seems like strings were found in values column")
				}
				streams[i].Strings = append(streams[i].Strings, s)
			}
		}
	}

	return streams, nil
}

var utc *time.Location

func millisecondsToTime(f float64) time.Time {
	seconds := f / 1000
	fullSeconds := math.Floor(seconds)
	nanoseconds := (seconds - fullSeconds) * 1e+9
	t := time.Unix(int64(fullSeconds), int64(nanoseconds))
	return t.In(utc)
}

func floatsToTimes(xf []float64) []time.Time {
	xt := []time.Time{}
	for _, f := range xf {
		mTime := millisecondsToTime(f)
		xt = append(xt, mTime)
	}
	return xt
}

func normalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

// FromCSV formats a compatible CSV as a FormattedData struct ready for mgJSON and returns it. Or returns an error
// The optional frame rate (fr) is used if timing data is not present
func FromCSV(src []byte, fr float64) (FormattedData, error) {
	var data FormattedData

	r := csv.NewReader(strings.NewReader(string(normalizeNewlines(src))))
	lines, err := r.ReadAll()
	if err != nil {
		return data, err
	}

	utc, err = time.LoadLocation("UTC")
	if err != nil {
		return data, err
	}

	//check if first line is headers
	if _, err := strconv.ParseFloat(lines[0][0], 64); err != nil {
		headers := lines[0]
		lines = lines[1:]
		streams, err := structureData(headers, lines)
		if err != nil {
			return data, err
		}
		if headers[0] == "milliseconds" && len(headers[1]) > 1 {
			data.Timing = floatsToTimes(streams[0].Values)
			streams = streams[1:]
			headers = headers[1:]
		}
		data.Streams = streams
	} else {
		streams, err := structureData([]string{"Data"}, lines)
		if err != nil {
			return data, err
		}
		data.Streams = streams
	}

	if len(data.Streams) < 1 {
		return data, fmt.Errorf("No valid data found")
	}

	if len(data.Timing) < 1 {
		for i := 0; i < len(data.Streams[0].Values); i++ {
			data.Timing = append(data.Timing, millisecondsToTime(float64(i)*1000.0/fr))
		}
	}

	return data, nil
}
