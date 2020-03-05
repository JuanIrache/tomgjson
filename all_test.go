package tomgjson

import (
	"fmt"
	"io/ioutil"
	"time"
)

func ExampleToMgjson() {

	stream1 := Stream{
		label:  "Prime numbers",
		values: []float64{2, 3, 5, 7},
	}

	stream2 := Stream{
		label:  "Non primes",
		values: []float64{4, 6, 8, 9},
	}

	now := time.Unix(0, 0)
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
		timing:  timing,
		streams: []Stream{stream1, stream2},
	}

	doc := ToMgjson(data, "Juan Irache")

	fmt.Println(string(doc))

	// Output:
	// {"version":"MGJSON2.0.0","creator":"Juan Irache","dynamicSamplesPresentB":true,"dynamicDataInfo":{"useTimecodeB":false,"utcInfo":{"precisionLength":3,"isGMT":true}},"dataOutline":[{"objectType":"dataDynamic","displayName":"Prime numbers","sampleSetID":"Stream0","dataType":{"type":"numberString","numberStringProperties":{"pattern":{"digitsInteger":1,"digitsDecimal":1,"isSigned":true},"range":{"occuring":{"min":2,"max":7},"legal":{"min":-2147483648,"max":2147483648}}}},"interpolation":"linear","hasExpectedFrequecyB":false,"sampleCount":4,"matchName":"Stream0"},{"objectType":"dataDynamic","displayName":"Non primes","sampleSetID":"Stream1","dataType":{"type":"numberString","numberStringProperties":{"pattern":{"digitsInteger":1,"digitsDecimal":1,"isSigned":true},"range":{"occuring":{"min":4,"max":9},"legal":{"min":-2147483648,"max":2147483648}}}},"interpolation":"linear","hasExpectedFrequecyB":false,"sampleCount":4,"matchName":"Stream1"}],"dataDynamicSamples":[{"sampleSetID":"Stream0","samples":[{"time":"1970-01-01T01:00:00.000Z","value":"+2.0"},{"time":"1970-01-01T01:00:10.000Z","value":"+3.0"},{"time":"1970-01-01T01:00:20.000Z","value":"+5.0"},{"time":"1970-01-01T01:00:30.000Z","value":"+7.0"}]},{"sampleSetID":"Stream1","samples":[{"time":"1970-01-01T01:00:00.000Z","value":"+4.0"},{"time":"1970-01-01T01:00:10.000Z","value":"+6.0"},{"time":"1970-01-01T01:00:20.000Z","value":"+8.0"},{"time":"1970-01-01T01:00:30.000Z","value":"+9.0"}]}]}
}

func ExampleFromCSV() {
	src, _ := ioutil.ReadFile("./sample_sources/multiple-data.csv")
	converted := FromCSV(src, 25.0)
	fmt.Printf(
		`The second stream is labelled as %q and its fifth sample at %f seconds is %v`,
		converted.streams[1].label,
		converted.timing[5].Sub(time.Unix(0, 0)).Seconds(),
		converted.streams[1].values[5],
	)
	//Output:
	//The second stream is labelled as "Signed 1k Perlin" and its fifth sample at 0.500000 seconds is 74.1837892462852
}

func ExampleFromGPX() {
	src, _ := ioutil.ReadFile("./sample_sources/gps-path.gpx")
	converted := FromGPX(src, true)
	sample := 10
	fmt.Printf(
		`At %v the position was %v: %f, %v: %f and the %v was %f `,
		converted.timing[sample].Format("15:04:05"),
		converted.streams[0].label,
		converted.streams[0].values[sample],
		converted.streams[1].label,
		converted.streams[1].values[sample],
		converted.streams[8].label,
		converted.streams[8].values[sample],
	)
	//Output:
	//At 11:45:55 the position was lat: 41.389212, lon: 2.147359 and the speed2d was 3.057694
}
