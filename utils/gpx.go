package utils

import (
	"encoding/xml"
	"fmt"
	"log"
)

// ReadGPX formats a compatible GPX file as a struct ready for mgJSON
func ReadGPX(src []byte) SourceData {
	var data SourceData

	type Trkpt struct {
		XMLName       xml.Name `xml:"trkpt"`
		Lat           float64  `xml:"lat,attr"`
		Lon           float64  `xml:"lon,attr"`
		Ele           float64  `xml:"ele"`
		Time          string   `xml:"time"`
		Magvar        float64  `xml:"magvar"`
		Geoidheight   float64  `xml:"geoidheight"`
		Fix           float64  `xml:"fix"`
		Sat           float64  `xml:"sat"`
		Hdop          float64  `xml:"hdop"`
		Vdop          float64  `xml:"vdop"`
		Pdop          float64  `xml:"pdop"`
		Ageofdgpsdata float64  `xml:"ageofdgpsdata"`
		Dgpsid        float64  `xml:"dgpsid"`
	}

	type Trkseg struct {
		XMLName xml.Name `xml:"trkseg"`
		Trkpt   []Trkpt  `xml:"trkpt"`
	}

	type Trk struct {
		XMLName xml.Name `xml:"trk"`
		Trkseg  []Trkseg `xml:"trkseg"`
	}

	type Gpx struct {
		XMLName xml.Name `xml:"gpx"`
		Trk     []Trk    `xml:"trk"`
	}

	gpx := Gpx{}

	err := xml.Unmarshal(src, &gpx)
	check(err)

	if len(gpx.Trk) < 1 {
		log.Panic("Error: No GPX tracks")
	}

	// Just reading one track for now

	if len(gpx.Trk[0].Trkseg) < 1 {
		log.Panic("Error: No GPX Trkseg")
	}

	// Just reading one trkseg for now

	if len(gpx.Trk[0].Trkseg[0].Trkpt) < 1 {
		log.Panic("Error: No GPX trkpt")
	}

	for i, trkpt := range gpx.Trk[0].Trkseg[0].Trkpt {
		if i == 0 {
			fmt.Println(trkpt)
		}
	}

	return data
}

// type stream struct {
// 	label  string
// 	values []float64
// }

// // SourceData structures data relevant for mgJSON files
// type SourceData struct {
// 	timing  []float64
// 	streams []stream
// }
