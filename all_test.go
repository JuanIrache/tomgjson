package tomgjson

import (
	"fmt"
	"io/ioutil"
	"time"
)

func ExampleToMgjson() {

	stream1 := Stream{
		Label:  "Prime numbers",
		Values: []float64{2, 3, 5, 7},
	}

	stream2 := Stream{
		Label:  "Non primes",
		Values: []float64{4, 6, 8, 9},
	}

	stream3 := Stream{
		Label: "Colors",
		Strings: []string{
			"Green",
			"Yellow",
			"Red",
			"Blue",
		},
	}

	utc, _ := time.LoadLocation("UTC")
	now := time.Unix(0, 0).In(utc)
	plus10, _ := time.ParseDuration("10s")
	plus20, _ := time.ParseDuration("20s")
	plus30, _ := time.ParseDuration("30s")
	timing := []time.Time{
		now,
		now.Add(plus10),
		now.Add(plus20),
		now.Add(plus30),
	}

	data := FormattedData{
		Timing:  timing,
		Streams: []Stream{stream1, stream2, stream3},
	}

	doc, _ := ToMgjson(data, "Juan Irache")

	fmt.Println(string(doc))

	// Output:
	// {"version":"MGJSON2.0.0","creator":"Juan Irache","dynamicSamplesPresentB":true,"dynamicDataInfo":{"useTimecodeB":false,"utcInfo":{"precisionLength":3,"isGMT":true}},"dataOutline":[{"objectType":"dataDynamic","displayName":"Prime numbers","sampleSetID":"Stream0","dataType":{"type":"numberString","numberStringProperties":{"pattern":{"digitsInteger":1,"digitsDecimal":1,"isSigned":true},"range":{"occuring":{"min":2,"max":7},"legal":{"min":2,"max":7}}},"paddedStringProperties":{"maxLen":0,"maxDigitsInStrLength":0,"eventMarkerB":false}},"interpolation":"linear","hasExpectedFrequecyB":false,"sampleCount":4,"matchName":"Stream0"},{"objectType":"dataDynamic","displayName":"Non primes","sampleSetID":"Stream1","dataType":{"type":"numberString","numberStringProperties":{"pattern":{"digitsInteger":1,"digitsDecimal":1,"isSigned":true},"range":{"occuring":{"min":4,"max":9},"legal":{"min":4,"max":9}}},"paddedStringProperties":{"maxLen":0,"maxDigitsInStrLength":0,"eventMarkerB":false}},"interpolation":"linear","hasExpectedFrequecyB":false,"sampleCount":4,"matchName":"Stream1"},{"objectType":"dataDynamic","displayName":"Colors","sampleSetID":"Stream2","dataType":{"type":"paddedString","numberStringProperties":{"pattern":{"digitsInteger":0,"digitsDecimal":0,"isSigned":false},"range":{"occuring":{"min":0,"max":0},"legal":{"min":0,"max":0}}},"paddedStringProperties":{"maxLen":6,"maxDigitsInStrLength":1,"eventMarkerB":false}},"interpolation":"hold","hasExpectedFrequecyB":false,"sampleCount":4,"matchName":"Stream2"}],"dataDynamicSamples":[{"sampleSetID":"Stream0","samples":[{"time":"1970-01-01T00:00:00.000Z","value":"+2.0"},{"time":"1970-01-01T00:00:10.000Z","value":"+3.0"},{"time":"1970-01-01T00:00:20.000Z","value":"+5.0"},{"time":"1970-01-01T00:00:30.000Z","value":"+7.0"}]},{"sampleSetID":"Stream1","samples":[{"time":"1970-01-01T00:00:00.000Z","value":"+4.0"},{"time":"1970-01-01T00:00:10.000Z","value":"+6.0"},{"time":"1970-01-01T00:00:20.000Z","value":"+8.0"},{"time":"1970-01-01T00:00:30.000Z","value":"+9.0"}]},{"sampleSetID":"Stream2","samples":[{"time":"1970-01-01T00:00:00.000Z","value":{"length":"5","str":"Green "}},{"time":"1970-01-01T00:00:10.000Z","value":{"length":"6","str":"Yellow"}},{"time":"1970-01-01T00:00:20.000Z","value":{"length":"3","str":"Red   "}},{"time":"1970-01-01T00:00:30.000Z","value":{"length":"4","str":"Blue  "}}]}]}
}

func ExampleFromCSV() {
	src, _ := ioutil.ReadFile("./sample_sources/multiple-data.csv")
	converted, _ := FromCSV(src, 0)
	sample := 5
	fmt.Printf(
		`The second stream is labelled as %q and its %q at %f seconds is %v`,
		converted.Streams[1].Label,
		converted.Streams[2].Strings[sample],
		converted.Timing[sample].Sub(time.Unix(0, 0)).Seconds(),
		converted.Streams[1].Values[sample],
	)
	//Output:
	//The second stream is labelled as "Signed 1k Perlin" and its "sample 6" at 0.500000 seconds is 74.1837892462852
}

func ExampleFromGPX() {
	src, _ := ioutil.ReadFile("./sample_sources/gps-path.gpx")
	converted, _ := FromGPX(src, true)
	sample := 10
	fmt.Printf(
		`At %v the position was %v: %f, %v: %f and the %v was %f `,
		converted.Timing[sample].Format("15:04:05"),
		converted.Streams[0].Label,
		converted.Streams[0].Values[sample],
		converted.Streams[1].Label,
		converted.Streams[1].Values[sample],
		converted.Streams[8].Label,
		converted.Streams[8].Values[sample],
	)
	//Output:
	//At 11:45:55 the position was lat (°): 41.389212, lon (°): 2.147359 and the speed2d (m/s) was 3.057694
}
