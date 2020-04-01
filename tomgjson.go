// Package tomgjson converts time based data sources to Adobe's mgJSON format for After Effects.
//
// Initially, this supports appropriately formatted CSV files and simple GPS files (see sample_sources).
//
// A live version of this app can be found here: https://goprotelemetryextractor.com/csv-gpx-to-mgjson/.
package tomgjson

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// Like math.Max but with ints
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns both sides of a float number as strings
func sides(n float64) (string, string) {
	sides := strings.Split(strconv.FormatFloat(math.Abs(n), 'f', -1, 64), ".")
	if len(sides) == 1 {
		sides = append(sides, "0")
	}
	if len(sides) != 2 {
		log.Panicf("Badly formatted float: %v %v", n, sides)
	}
	return sides[0], sides[1]
}

// Make sure float values are within mgJSON's valid values
const largestMgjsonNum = 2147483648.0

func validValue(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return math.Max(math.Min(v, largestMgjsonNum), -largestMgjsonNum)
}

// Stream contains a slice of values or strings and their label
// The slices must be of the same length as the timing slice in their parent's FormattedData
// Only one of the slices must be present, not both
type Stream struct {
	Label   string
	Values  []float64
	Strings []string
}

// FormattedData is the struct accepted by ToMgjson.
// It consists of a slice of timestamps and a slice with all the streams of labelled values (floats for now)
type FormattedData struct {
	Timing  []time.Time
	Streams []Stream
}

// mgJSON structure. For now, only the fields we are using are specified
type utcInfo struct {
	PrecisionLength int  `json:"precisionLength"`
	IsGMT           bool `json:"isGMT"`
}

type dynamicDataInfo struct {
	UseTimecodeB bool    `json:"useTimecodeB"`
	UtcInfo      utcInfo `json:"utcInfo"`
}

type pattern struct {
	DigitsInteger int  `json:"digitsInteger"`
	DigitsDecimal int  `json:"digitsDecimal"`
	IsSigned      bool `json:"isSigned"`
}

type minmax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type mRange struct {
	Occuring minmax `json:"occuring"`
	Legal    minmax `json:"legal"`
}

type numberStringProperties struct {
	Pattern pattern `json:"pattern"`
	Range   mRange  `json:"range"`
}

type paddedStringProperties struct {
	MaxLen               int  `json:"maxLen"`
	MaxDigitsInStrLength int  `json:"maxDigitsInStrLength"`
	EventMarkerB         bool `json:"eventMarkerB"`
}

type dataType struct {
	Type                   string                 `json:"type"`
	NumberStringProperties numberStringProperties `json:"numberStringProperties"`
	PaddedStringProperties paddedStringProperties `json:"paddedStringProperties"`
}

type singleDataOutline struct {
	ObjectType            string   `json:"objectType"`
	DisplayName           string   `json:"displayName"`
	SampleSetID           string   `json:"sampleSetID"`
	DataType              dataType `json:"dataType"`
	Interpolation         string   `json:"interpolation"`
	HasExpectedFrequencyB bool     `json:"hasExpectedFrequecyB"`
	SampleCount           int      `json:"sampleCount"`
	MatchName             string   `json:"matchName"`
}

type paddedStringValue struct {
	Length string `json:"length"`
	Str    string `json:"str"`
}

type sample struct {
	Time  string      `json:"time"`
	Value interface{} `json:"value"`
}

type dataDynamicSample struct {
	SampleSetID string   `json:"sampleSetID"`
	Samples     []sample `json:"samples"`
}

type mgjson struct {
	Version                string              `json:"version"`
	Creator                string              `json:"creator"`
	DynamicSamplesPresentB bool                `json:"dynamicSamplesPresentB"`
	DynamicDataInfo        dynamicDataInfo     `json:"dynamicDataInfo"`
	DataOutline            []singleDataOutline `json:"dataOutline"`
	DataDynamicSamples     []dataDynamicSample `json:"dataDynamicSamples"`
}

// ToMgjson receives a formatted source data (FormattedData) and a creator or author name
// and returns formatted mgjson ready to write to a file
// compatible with Adobe After Effects data-driven animations (or an error)
func ToMgjson(sd FormattedData, creator string) ([]byte, error) {

	if len(sd.Streams) < 1 {
		return nil, fmt.Errorf("No streams found")
	}

	if len(sd.Timing) < 1 {
		return nil, fmt.Errorf("No timing data")
	}

	//Hardcode non configurable values (for now)
	data := mgjson{
		Version:                "MGJSON2.0.0",
		Creator:                creator,
		DynamicSamplesPresentB: true,
		DynamicDataInfo: dynamicDataInfo{
			UseTimecodeB: false,
			UtcInfo: utcInfo{
				PrecisionLength: 3,
				IsGMT:           true,
			},
		},
		DataOutline:        []singleDataOutline{},
		DataDynamicSamples: []dataDynamicSample{},
	}

	for i, stream := range sd.Streams {
		sName := fmt.Sprintf("Stream%d", i)
		min := largestMgjsonNum
		max := -largestMgjsonNum
		digitsInteger := 0
		digitsDecimal := 0
		maxLen := 0
		maxDigitsInStrLength := 0

		for _, v := range stream.Values {
			v = validValue(v)
			min = math.Min(min, v)
			max = math.Max(max, v)
			integer, decimal := sides(v)
			digitsInteger = maxInt(digitsInteger, len(integer))
			digitsDecimal = maxInt(digitsDecimal, len(decimal))
		}

		for _, v := range stream.Strings {
			maxLen = maxInt(maxLen, len(v))
			maxDigitsInStrLength = len(strconv.Itoa(maxLen))
		}

		var thisDataType dataType
		var thisInterpolation string
		var thisSampleCount int

		if len(stream.Values) > 0 {

			thisDataType = dataType{
				Type: "numberString",
				NumberStringProperties: numberStringProperties{
					Pattern: pattern{
						DigitsInteger: digitsInteger,
						DigitsDecimal: digitsDecimal,
						IsSigned:      true,
					},
					Range: mRange{
						Occuring: minmax{min, max},
						Legal:    minmax{-largestMgjsonNum, largestMgjsonNum},
					},
				},
			}
			thisInterpolation = "linear"
			thisSampleCount = len(stream.Values)

		} else if len(stream.Strings) > 0 {

			thisDataType = dataType{
				Type: "paddedString",
				PaddedStringProperties: paddedStringProperties{
					MaxLen:               maxLen,
					MaxDigitsInStrLength: maxDigitsInStrLength,
					EventMarkerB:         false,
				},
			}
			thisInterpolation = "hold"
			thisSampleCount = len(stream.Strings)

		}

		if len(sd.Timing) != thisSampleCount {
			return nil, fmt.Errorf("Timing data does not match slice length")
		}

		data.DataOutline = append(data.DataOutline, singleDataOutline{
			ObjectType:            "dataDynamic",
			DisplayName:           stream.Label,
			SampleSetID:           sName,
			DataType:              thisDataType,
			Interpolation:         thisInterpolation,
			HasExpectedFrequencyB: false,
			SampleCount:           thisSampleCount,
			MatchName:             sName,
		})

		streamSamples := []sample{}

		for i, v := range stream.Values {
			v = validValue(v)
			paddedValue := fmt.Sprintf("%+0*.*f", digitsInteger+digitsDecimal+2, digitsDecimal, v)
			timeStr := sd.Timing[i].Format("2006-01-02T15:04:05.000Z")
			streamSamples = append(streamSamples, sample{
				Time:  timeStr,
				Value: paddedValue,
			})
		}

		for i, v := range stream.Strings {
			stringValue := paddedStringValue{
				Length: fmt.Sprintf("%0*d", maxDigitsInStrLength, len(v)),
				Str:    fmt.Sprintf("%-*v", maxLen, v),
			}
			timeStr := sd.Timing[i].Format("2006-01-02T15:04:05.000Z")
			streamSamples = append(streamSamples, sample{
				Time:  timeStr,
				Value: stringValue,
			})
		}

		data.DataDynamicSamples = append(data.DataDynamicSamples, dataDynamicSample{
			SampleSetID: sName,
			Samples:     streamSamples,
		})
	}

	doc, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
