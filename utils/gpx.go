package utils

import (
	"encoding/xml"
	"log"
	"math"
	"time"
)

func millis(t time.Time) float64 {
	return float64(t.UnixNano()) / float64(time.Millisecond)
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func radiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func distanceInMBetweenEarthCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	earthRadiusM := 6378137.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusM * c
}

func angleFromCoordinate(lat1, lon1, lat2, lon2 float64) float64 {

	dLon := degreesToRadians(lon2 - lon1)

	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)

	brng := math.Atan2(y, x)

	brng = radiansToDegrees(brng)
	brng = math.Mod(brng+360, 360)
	brng = 360 - brng // count degrees counter-clockwise - remove to make clockwise

	return brng
}

var ids = map[string]int{
	// Supported GPX streams
	"lat":           0,
	"lon":           1,
	"ele":           2,
	"magvar":        3,
	"geoidheight":   4,
	"fix":           5,
	"sat":           6,
	"hdop":          7,
	"vdop":          8,
	"pdop":          9,
	"ageofdgpsdata": 10,
	"dgpsid":        11,
	// Calculated streams
	"distance":     12,
	"speed":        13,
	"acceleration": 14,
	"course":       15,
	"slope":        16,
}

func appendToStream(data SourceData, v *float64, n string) stream {
	st := data.streams[ids[n]]
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

// ReadGPX formats a compatible GPX file as a struct ready for mgJSON. If extra, will compute additional streams
func ReadGPX(src []byte, extra bool) SourceData {
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

	// One stream for each of the supported trkpt and custom fields
	data.streams = make([]stream, len(ids))
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
		data.streams[ids["lat"]] = appendToStream(data, trkpt.Lat, "lat")
		data.streams[ids["lon"]] = appendToStream(data, trkpt.Lon, "lon")
		data.streams[ids["ele"]] = appendToStream(data, trkpt.Ele, "ele")
		data.streams[ids["magvar"]] = appendToStream(data, trkpt.Magvar, "magvar")
		data.streams[ids["geoidheight"]] = appendToStream(data, trkpt.Geoidheight, "geoidheight")
		data.streams[ids["fix"]] = appendToStream(data, trkpt.Fix, "fix")
		data.streams[ids["sat"]] = appendToStream(data, trkpt.Sat, "sat")
		data.streams[ids["hdop"]] = appendToStream(data, trkpt.Hdop, "hdop")
		data.streams[ids["vdop"]] = appendToStream(data, trkpt.Vdop, "vdop")
		data.streams[ids["pdop"]] = appendToStream(data, trkpt.Pdop, "pdop")
		data.streams[ids["ageofdgpsdata"]] = appendToStream(data, trkpt.Ageofdgpsdata, "ageofdgpsdata")
		data.streams[ids["dgpsid"]] = appendToStream(data, trkpt.Dgpsid, "dgpsid")

		// Computed streams
		if extra {
			distance := 0.0
			speed := 0.0
			acceleration := 0.0
			course := 0.0
			slope := 0.0
			if i > 0 {
				prevLat := data.streams[ids["lat"]].values[i-1]
				prevLon := data.streams[ids["lon"]].values[i-1]
				distance = distanceInMBetweenEarthCoordinates(*trkpt.Lat, *trkpt.Lon, prevLat, prevLon)
				duration := (data.timing[i] - data.timing[i-1]) / 1000
				speed = distance / duration
				acceleration = speed
				course = angleFromCoordinate(*trkpt.Lat, *trkpt.Lon, prevLat, prevLon)
				prevEle := data.streams[ids["ele"]].values[i-1]
				eleDiff := *trkpt.Ele - prevEle
				slope = math.Atan2(eleDiff, distance)
				slope = radiansToDegrees(slope)
				if i > 1 {
					prevDistance := data.streams[ids["distance"]].values[i-1]
					distance += prevDistance
					prevSpeed := data.streams[ids["speed"]].values[i-1]
					speedChange := speed - prevSpeed
					acceleration = speedChange / duration
				}
			}

			data.streams[ids["distance"]] = appendToStream(data, &distance, "distance")
			data.streams[ids["speed"]] = appendToStream(data, &speed, "speed")
			data.streams[ids["acceleration"]] = appendToStream(data, &acceleration, "acceleration")
			data.streams[ids["course"]] = appendToStream(data, &course, "course")
			data.streams[ids["slope"]] = appendToStream(data, &slope, "slope")
		}
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
