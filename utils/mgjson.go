package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type stream struct {
	label  string
	values []float64
}

// SourceData structures data relevant for mgJSON files
type SourceData struct {
	timing  []float64
	streams []stream
}

// Destructured mgJSON fields. For now, only the fields we are using are specified
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
	sides := strings.Split(strconv.FormatFloat(n, 'f', -1, 64), ".")
	if len(sides) == 1 {
		sides = append(sides, "0")
	}
	if len(sides) != 2 {
		log.Panicf("Badly formatted float: %v %v", n, sides)
	}
	return sides[0], sides[1]
}

func FormatMgjson(sd SourceData, creator string) mgjson {

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

	for i, stream := range sd.streams {
		sName := fmt.Sprintf("%d_%v", i, stream.label)
		min := math.Inf(1)
		max := math.Inf(-1)
		digitsInteger := 0
		digitsDecimal := 0
		for _, v := range stream.values {
			min = math.Min(min, v)
			max = math.Max(min, v)
			integer, decimal := sides(v)
			digitsInteger = mMax(digitsInteger, len(integer))
			digitsDecimal = mMax(digitsDecimal, len(decimal))
		}
		data.DataOutline = append(data.DataOutline, singleDataOutline{
			ObjectType:  "dataDynamic",
			DisplayName: stream.label,
			SampleSetID: sName,
			DataType: dataType{
				Type: "numberString",
				NumberStringProperties: numberStringProperties{
					Pattern: pattern{
						DigitsInteger: digitsInteger,
						DigitsDecimal: digitsDecimal,
						//To-Do signed
						IsSigned: false,
					},
					Range: mRange{
						Occuring: minmax{min, max},
						Legal:    minmax{-largestMgjsonNum, largestMgjsonNum},
					},
				},
			},
			Interpolation:         "linear",
			HasExpectedFrequencyB: false,
			SampleCount:           len(stream.values),
			MatchName:             sName,
		})

		streamSamples := []sample{}

		for _, v := range stream.values {
			integer, decimal := sides(v)
			paddedInteger := fmt.Sprintf("%0"+fmt.Sprint(digitsInteger)+"v", integer)
			paddedDecimal := fmt.Sprintf("%-0"+fmt.Sprint(digitsDecimal)+"v", decimal)
			streamSamples = append(streamSamples, sample{
				Time:  "x",
				Value: paddedInteger + "." + paddedDecimal,
			})
		}

		data.DataDynamicSamples = append(data.DataDynamicSamples, dataDynamicSample{
			SampleSetID: sName,
			Samples:     streamSamples,
		})
	}

	return data
}
