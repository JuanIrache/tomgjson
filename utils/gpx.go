package utils

import (
	"encoding/xml"
	"log"
	"time"
)

func appendToStream(st stream, v *float64, n string) stream {
	if v != nil {
		st.values = append(st.values, *v)
		// Name confirmed streams
		st.label = n
	} else {
		// Appending zeros for now, could potentially add interpolated values based on previous, next and time
		st.values = append(st.values, 0)
	}
	return st
}

func millis(t time.Time) float64 {
	return float64(t.UnixNano()) / float64(time.Millisecond)
}

// ReadGPX formats a compatible GPX file as a struct ready for mgJSON
func ReadGPX(src []byte) SourceData {
	var data SourceData

	type Trkpt struct {
		XMLName       xml.Name `xml:"trkpt"`
		Time          *string  `xml:"time"`
		Lat           *float64 `xml:"lat,attr"`
		Lon           *float64 `xml:"lon,attr"`
		Ele           *float64 `xml:"ele"`
		Magvar        *float64 `xml:"magvar"`
		Geoidheight   *float64 `xml:"geoidheight"`
		Fix           *float64 `xml:"fix"`
		Sat           *float64 `xml:"sat"`
		Hdop          *float64 `xml:"hdop"`
		Vdop          *float64 `xml:"vdop"`
		Pdop          *float64 `xml:"pdop"`
		Ageofdgpsdata *float64 `xml:"ageofdgpsdata"`
		Dgpsid        *float64 `xml:"dgpsid"`
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

	// One stream for each of the supported trkpt fields
	data.streams = make([]stream, 12)
	data.timing = make([]float64, len(gpx.Trk[0].Trkseg[0].Trkpt))

	for _, st := range data.streams {
		st.values = make([]float64, len(gpx.Trk[0].Trkseg[0].Trkpt))
	}

	for i, trkpt := range gpx.Trk[0].Trkseg[0].Trkpt {
		if trkpt.Time == nil {
			log.Panic("Error: Missing timiing data in GPX")
		}
		t, err := time.Parse(time.RFC3339, *trkpt.Time)
		check(err)

		data.timing[i] = millis(t)
		data.streams[0] = appendToStream(data.streams[0], trkpt.Lat, "lat")
		data.streams[1] = appendToStream(data.streams[1], trkpt.Lon, "lon")
		data.streams[2] = appendToStream(data.streams[2], trkpt.Ele, "ele")
		data.streams[3] = appendToStream(data.streams[3], trkpt.Magvar, "magvar")
		data.streams[4] = appendToStream(data.streams[4], trkpt.Geoidheight, "geoidheight")
		data.streams[5] = appendToStream(data.streams[5], trkpt.Fix, "fix")
		data.streams[6] = appendToStream(data.streams[6], trkpt.Sat, "sat")
		data.streams[7] = appendToStream(data.streams[7], trkpt.Hdop, "hdop")
		data.streams[8] = appendToStream(data.streams[8], trkpt.Vdop, "vdop")
		data.streams[9] = appendToStream(data.streams[9], trkpt.Pdop, "pdop")
		data.streams[10] = appendToStream(data.streams[10], trkpt.Ageofdgpsdata, "ageofdgpsdata")
		data.streams[11] = appendToStream(data.streams[11], trkpt.Dgpsid, "dgpsid")
	}

	// Clean up unconfirmed streams
	for i := len(data.streams) - 1; i >= 0; i-- {
		if len(data.streams[i].label) < 1 {
			copy(data.streams[i:], data.streams[i+1:])
			data.streams[len(data.streams)-1] = stream{}
			data.streams = data.streams[:len(data.streams)-1]
		}
	}

	return data
}
