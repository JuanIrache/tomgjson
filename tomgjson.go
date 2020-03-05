/*
# tomgjson

Experimenting with converting other time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson)

The first iteration will be very limited in terms of customization and supported formats, but the provided parser funcs can potentially be adapted to any source that associates time with values.

For instructions on how to use mgJSON files in After Effects, see [Work with Data-driven animation](https://helpx.adobe.com/after-effects/using/data-driven-animations.html)

## Sample project templates

You can find sample After Effects projects that use mgJSON files on the [GoPro Telemetry Extractor page](https://goprotelemetryextractor.com). Look for the **Lite templates**.

## Videos made with mgJSON

- [mgJSON playlist on YouTube](https://www.youtube.com/playlist?list=PLgoeWSWqXedI7FbZccAEudt2_t8qPX0Px)

If you create something with mgJSON, let me know and I'll add it to the list.

## Software that supports mgJSON

These apps can output mgJSON files:

- [GoPro Telemetry Extractor](https://goprotelemetryextractor.com)
- [DJI Telemetry Overlay](https://djitelemetryoverlay.com)
*/

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

// Stream contains a slice of values and their label
// The slices of floats must be of the same length as the timing slice in their parent's FormattedData
type Stream struct {
	Label  string
	Values []float64
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

type dataType struct {
	Type                   string                 `json:"type"`
	NumberStringProperties numberStringProperties `json:"numberStringProperties"`
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

type sample struct {
	Time  string `json:"time"`
	Value string `json:"value"`
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

// ToMgjson receives a formatted source data (FormattedData) and a creator or author name
// and returns formatted mgjson ready to write to a file
func ToMgjson(sd FormattedData, creator string) []byte {

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

	largestMgjsonNum := 2147483648.0

	for i, stream := range sd.Streams {
		sName := fmt.Sprintf("Stream%d", i)
		min := math.Inf(1)
		max := math.Inf(-1)
		digitsInteger := 0
		digitsDecimal := 0
		for _, v := range stream.Values {
			min = math.Min(min, v)
			max = math.Max(min, v)
			integer, decimal := sides(v)
			digitsInteger = mMax(digitsInteger, len(integer))
			digitsDecimal = mMax(digitsDecimal, len(decimal))
		}
		data.DataOutline = append(data.DataOutline, singleDataOutline{
			ObjectType:  "dataDynamic",
			DisplayName: stream.Label,
			SampleSetID: sName,
			DataType: dataType{
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
			},
			Interpolation:         "linear",
			HasExpectedFrequencyB: false,
			SampleCount:           len(stream.Values),
			MatchName:             sName,
		})

		streamSamples := []sample{}
		for i, v := range stream.Values {
			paddedValue := fmt.Sprintf("%+0*.*f", digitsInteger+digitsDecimal+2, digitsDecimal, v)
			timeStr := sd.Timing[i].Format("2006-01-02T15:04:05.000Z")
			streamSamples = append(streamSamples, sample{
				Time:  timeStr,
				Value: paddedValue,
			})
		}

		data.DataDynamicSamples = append(data.DataDynamicSamples, dataDynamicSample{
			SampleSetID: sName,
			Samples:     streamSamples,
		})
	}

	doc, err := json.Marshal(data)
	check(err)
	return doc
}
